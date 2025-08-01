package service

import (
	"encoding/json"
	"fmt"

	"github.com/supabase-community/supabase-go"
)

// UserProfile represents a user profile from the database
type UserProfile struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Credits int    `json:"credits"`
}

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

	resp, _, err := us.supabaseClient.From("profiles").Select("credits", "", false).Eq("id", userID).Single().Execute()
	if err != nil {
		return 0, fmt.Errorf("user profile not found: %w", err)
	}

	var profile UserProfile
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

	if credits < 0 {
		return fmt.Errorf("invalid credit value: credits cannot be negative")
	}

	// First check if user exists
	_, _, err := us.supabaseClient.From("profiles").Select("id", "", false).Eq("id", userID).Single().Execute()
	if err != nil {
		// Create user profile if it doesn't exist
		err = us.createUserProfile(userID)
		if err != nil {
			return fmt.Errorf("failed to create user profile: %w", err)
		}
	}

	// Get current credits and add new credits atomically
	currentCredits, err := us.GetUserCredits(userID)
	if err != nil {
		return fmt.Errorf("failed to get current credits: %w", err)
	}

	newTotal := currentCredits + credits
	data := map[string]interface{}{
		"credits": newTotal,
	}

	_, _, err = us.supabaseClient.From("profiles").Update(data, "", "").Eq("id", userID).Execute()
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

	_, _, err := us.supabaseClient.From("profiles").Insert(data, false, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to create user profile: %w", err)
	}

	return nil
}
