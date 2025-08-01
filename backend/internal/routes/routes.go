package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/supabase-community/supabase-go"
	"github.com/verse91/ytb-clipy/backend/internal/controller"
	"github.com/verse91/ytb-clipy/backend/internal/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(router fiber.Router, supabaseClient *supabase.Client) {
	router.Get("/", homepageHandler)

	// User credit management routes with authentication
	router.Get("/user/:userID/credits", middleware.UserAuthMiddleware, func(c fiber.Ctx) error {
		userController := controller.NewUserController(supabaseClient)
		return userController.GetUserCredits(c)
	})

	// Admin only routes for credit management
	router.Post("/user/:userID/credits/update", middleware.AdminAuthMiddleware, func(c fiber.Ctx) error {
		userController := controller.NewUserController(supabaseClient)
		return userController.UpdateUserCredits(c)
	})

	router.Post("/user/:userID/credits/add", middleware.AdminAuthMiddleware, func(c fiber.Ctx) error {
		userController := controller.NewUserController(supabaseClient)
		return userController.AddUserCredits(c)
	})

	// Video routes
	router.Post("/video/download", func(c fiber.Ctx) error {
		videoController := controller.NewVideoController(supabaseClient)
		return videoController.DownloadHandler(c)
	})

	router.Get("/userinfo/", func(c fiber.Ctx) error {
		userController := controller.NewUserController(supabaseClient)
		return userController.GetUserById(c)
	})

	router.Get("/user/", func(c fiber.Ctx) error {
		userController := controller.NewUserController(supabaseClient)
		return userController.UserHandler(c)
	})
}

// Homepage handler
func homepageHandler(c fiber.Ctx) error {
	return c.SendString("Homepage")
}
