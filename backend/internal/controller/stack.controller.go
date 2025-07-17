package controller

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type StackController struct{}

func NewStackController() *StackController {
	return &StackController{}
}

func (sc *StackController) StackHandler(c fiber.Ctx) error {
	data := fiber.Map{
		"message": "hello",
		"id":      123,
	}

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

