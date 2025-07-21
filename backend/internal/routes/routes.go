package router

import (
	c "github.com/verse91/ytb-clipy/backend/internal/controller"

	"github.com/gofiber/fiber/v3"
	"github.com/supabase-community/supabase-go"
)

// SetupRoutes configures all API routes
func SetupRoutes(router fiber.Router, supabaseClient *supabase.Client) {
	router.Get("/", homepageHandler)
	// router.Post("/video", c.NewVideoController(videoRepo).VideoProcessHandler)
	router.Get("/video/download", c.NewVideoController(supabaseClient).DownloadHandler)
	router.Get("/userinfo/", c.NewUserController().GetUserById).Name("user")
	router.Get("/user/", c.NewUserController().UserHandler).Name("user") // /?name=...&id=...&age=...
}

// Homepage handler
func homepageHandler(c fiber.Ctx) error {
	return c.SendString("Homepage")
}
