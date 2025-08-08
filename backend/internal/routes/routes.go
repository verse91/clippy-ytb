package router

import (
	"github.com/gofiber/fiber/v3"
	"github.com/supabase-community/supabase-go"
	"github.com/verse91/ytb-clipy/backend/internal/config"
	"github.com/verse91/ytb-clipy/backend/internal/controller"
	"github.com/verse91/ytb-clipy/backend/internal/middleware"
)

func SetupRoutes(router fiber.Router, supabaseClient *supabase.Client, config *config.Config) {
	userController := controller.NewUserController(supabaseClient, config)
	videoController := controller.NewVideoController(supabaseClient)

	router.Get("/", homepageHandler)

	router.Get("/user/:userID/credits", middleware.UserAuthMiddleware, func(c fiber.Ctx) error {
		return userController.GetUserCredits(&c)
	})

	router.Post("/user/:userID/credits/update", middleware.AdminAuthMiddleware, func(c fiber.Ctx) error {
		return userController.UpdateUserCredits(&c)
	})

	router.Post("/user/:userID/credits/add", middleware.AdminAuthMiddleware, func(c fiber.Ctx) error {
		return userController.AddUserCredits(&c)
	})

	router.Post("/video/download", func(c fiber.Ctx) error {
		return videoController.DownloadHandler(c)
	})

	router.Get("/video/download/:id", func(c fiber.Ctx) error {
		return videoController.GetDownloadStatus(c)
	})

	router.Post("/video/download/time-range", func(c fiber.Ctx) error {
		return videoController.DownloadTimeRangeHandler(c)
	})

	router.Get("/video/download/time-range/:id", func(c fiber.Ctx) error {
		return videoController.GetTimeRangeDownloadStatusHandler(c)
	})

	router.Get("/user/info", func(c fiber.Ctx) error {
		return userController.GetUserById(&c)
	})

	router.Get("/user/profile", func(c fiber.Ctx) error {
		return userController.UserHandler(&c)
	})
}

func homepageHandler(c fiber.Ctx) error {
	return c.SendString("Homepage")
}
