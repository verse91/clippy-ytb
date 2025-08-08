package controller

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"github.com/supabase-community/supabase-go"
	"github.com/verse91/ytb-clipy/backend/internal/service"
	"github.com/verse91/ytb-clipy/backend/pkg/response"
)

type Person struct {
	Name string `query:"name"`
	ID   int    `query:"id"`
	Age  int    `query:"age"`
}

type UserController struct {
	UserService *service.UserService
}

func NewUserController(supabaseClient *supabase.Client) *UserController {
	return &UserController{
		UserService: service.NewUserService(supabaseClient),
	}
}

func (uc *UserController) UserHandler(c fiber.Ctx) error {
	p := new(Person)

	if err := c.Bind().Query(p); err != nil {
		return err
	}

	data := fiber.Map{
		"message": "hello " + p.Name,
		"id":      p.ID,
		"age":     p.Age,
	}

	return c.Status(fiber.StatusOK).JSON(data)
}

// GetUserCredits returns the current credit balance for a user
func (uc *UserController) GetUserCredits(c fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return response.ErrorResponse(c, 400, "User ID is required")
	}

	credits, err := uc.UserService.GetUserCredits(userID)
	if err != nil {
		return response.ErrorResponse(c, 500, "Failed to get user credits")
	}

	return response.SuccessResponse(c, response.SuccessCode, fiber.Map{
		"user_id": userID,
		"credits": credits,
	})
}

// UpdateUserCredits sets the credit balance for a user (admin only)
func (uc *UserController) UpdateUserCredits(c fiber.Ctx) error {
	// Additional authorization check - verify admin key is present
	adminKey := c.Get("X-Admin-Key")
	if adminKey == "" {
		return response.ErrorResponse(c, 401, "Admin authentication required")
	}

	userID := c.Params("userID")
	if userID == "" {
		return response.ErrorResponse(c, 400, "User ID is required")
	}

	creditsStr := c.FormValue("credits")
	credits, err := strconv.Atoi(creditsStr)
	if err != nil {
		return response.ErrorResponse(c, 400, "Invalid credits value")
	}

	if credits < 0 {
		return response.ErrorResponse(c, 400, "Credits cannot be negative")
	}

	err = uc.UserService.UpdateUserCredits(userID, credits)
	if err != nil {
		return response.ErrorResponse(c, 500, "Failed to update user credits")
	}

	return response.SuccessResponse(c, response.SuccessCode, fiber.Map{
		"user_id": userID,
		"credits": credits,
		"message": "Credits updated successfully",
	})
}

// AddUserCredits adds credits to a user's balance
func (uc *UserController) AddUserCredits(c fiber.Ctx) error {
	userID := c.Params("userID")
	if userID == "" {
		return response.ErrorResponse(c, 400, "User ID is required")
	}

	creditsStr := c.FormValue("credits")
	credits, err := strconv.Atoi(creditsStr)
	if err != nil {
		return response.ErrorResponse(c, 400, "Invalid credits value")
	}

	if credits <= 0 {
		return response.ErrorResponse(c, 400, "Credits must be positive")
	}

	err = uc.UserService.AddUserCredits(userID, credits)
	if err != nil {
		return response.ErrorResponse(c, 500, "Failed to add user credits")
	}

	return response.SuccessResponse(c, response.SuccessCode, fiber.Map{
		"user_id":       userID,
		"credits_added": credits,
		"message":       "Credits added successfully",
	})
}

// controller -> service -> repo -> models -> database
func (uc *UserController) GetUserById(c fiber.Ctx) error {
	// if err := someFunctionThatMightFail(); err != nil {
	// 	return response.ErrorResponse(c, 200003, "not neccessary")
	// }
	return response.SuccessResponse(c, response.SuccessCode, []string{"abc"})
}
