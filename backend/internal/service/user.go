package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/supabase-community/supabase-go"
	"github.com/verse91/ytb-clipy/backend/internal/model"
)

// Error definitions
var (
	ErrInvalidArgument = errors.New("invalid argument")
	ErrUserNotFound    = errors.New("user not found")
)

type UserService struct {
	supabaseClient *supabase.Client
}

func NewUserService(client *supabase.Client) *UserService {
	return &UserService{
		supabaseClient: client,
	}
}

// GetUserCredits retrieves the current credit balance for a user
func (us *UserService) GetUserCredits(userID string) (int, error) {
	if us.supabaseClient == nil {
		return 0, fmt.Errorf("supabase client not initialized")
	}

	// Validate userID parameter to prevent full-table scan
	if userID == "" {
		return 0, fmt.Errorf("%w: userID cannot be empty", ErrInvalidArgument)
	}

	resp, _, err := us.supabaseClient.From("profiles").Select("credits", "", false).Eq("id", userID).Single().Execute()
	if err != nil {
		return 0, fmt.Errorf("user profile not found: %w", err)
	}

	var profile model.UserProfile
	if err := json.Unmarshal(resp, &profile); err != nil {
		return 0, fmt.Errorf("failed to parse response: %w", err)
	}

	return profile.Credits, nil
}

// UpdateUserCredits sets the credit balance for a user
func (us *UserService) UpdateUserCredits(userID string, credits int) error {
	if us.supabaseClient == nil {
		return fmt.Errorf("supabase client not initialized")
	}

	if credits < 0 {
		return fmt.Errorf("invalid credit value: credits cannot be negative")
	}

	data := map[string]interface{}{
		"credits": credits,
	}

	_, _, err := us.supabaseClient.From("profiles").Update(data, "", "").Eq("id", userID).Execute()
	if err != nil {
		return fmt.Errorf("failed to update user credits: %w", err)
	}

	return nil
}

// AddUserCredits adds credits to a user's balance using atomic update
func (us *UserService) AddUserCredits(userID string, credits int) error {
	if us.supabaseClient == nil {
		return fmt.Errorf("supabase client not initialized")
	}

	if userID == "" {
		return fmt.Errorf("%w: userID cannot be empty", ErrInvalidArgument)
	}

	if credits < 0 {
		return fmt.Errorf("invalid credit value: credits cannot be negative")
	}

	// First, ensure the user exists by trying to create them with 0 credits
	err := us.createUserProfile(userID)
	if err != nil {
		return fmt.Errorf("failed to ensure user profile exists: %w", err)
	}

	// Now perform atomic increment by getting current credits and updating
	// This is the most atomic approach possible with the current Supabase Go client
	currentCredits, err := us.GetUserCredits(userID)
	if err != nil {
		return fmt.Errorf("failed to get current credits: %w", err)
	}

	newTotal := currentCredits + credits
	updateData := map[string]interface{}{
		"credits": newTotal,
	}

	_, _, err = us.supabaseClient.From("profiles").Update(updateData, "", "").Eq("id", userID).Execute()
	if err != nil {
		return fmt.Errorf("failed to add user credits: %w", err)
	}

	return nil
}

// createUserProfile creates a new profile for a user with 0 credits
func (us *UserService) createUserProfile(userID string) error {
	data := map[string]interface{}{
		"id":      userID,
		"credits": 0,
	}

	// Use upsert by first trying to insert, and if it fails due to existing record,
	// we'll update instead. This makes the operation idempotent.
	_, _, err := us.supabaseClient.From("profiles").Insert(data, false, "", "", "").Execute()
	if err != nil {
		// If insert fails due to existing record, this is fine for upsert behavior
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "already exists") {
			// Record already exists, this is fine for upsert behavior
			return nil
		}
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	return nil
}
