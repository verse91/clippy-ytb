package router

import (
	c "github.com/verse91/ytb-clipy/backend/internal/controller"

	"github.com/gofiber/fiber/v3"
)

// SetupRoutes configures all API routes
func SetupRoutes(router fiber.Router) {
	router.Get("/", homepageHandler)
	router.Get("/stack", c.NewStackController().StackHandler)
	// router.Post("/video", c.NewVideoController().VideoProcessHandler)
	router.Get("/video/download", c.NewVideoController().DownloadHandler)
	router.Get("/userinfo/", c.NewUserController().GetUserById).Name("user")
	router.Get("/user/", c.NewUserController().UserHandler).Name("user") // /?name=...&id=...&age=...
}

// Homepage handler
func homepageHandler(c fiber.Ctx) error {
	return c.SendString("Homepage")
}
