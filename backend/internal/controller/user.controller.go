package controller

import (
	"github.com/verse91/ytb-clipy/backend/internal/service"
	"github.com/verse91/ytb-clipy/backend/pkg/response"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type Person struct {
	Name string `query:"name"`
	ID   int    `query:"id"`
	Age  int    `query:"age"`
}

type UserController struct {
	UserService *service.UserService
}

func NewUserController() *UserController {
	return &UserController{
		UserService: service.NewUserService(),
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

	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

// controller -> service -> repo -> models -> database
func (uc *UserController) GetUserById(c fiber.Ctx) error {
	// if err := someFunctionThatMightFail(); err != nil {
	// 	return response.ErrorResponse(c, 200003, "not neccessary")
	// }
	return response.SuccessResponse(c, response.SuccessCode, []string{"abc"})
}
