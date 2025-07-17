package response

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v3"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"messsage"`
	Data    interface{} `json:"data"`
}

// I just give 2 examples to use fiber.Map or struct depends on you
func SuccessResponse(c fiber.Ctx, code int, data interface{}) error {

	response := ResponseData{
		Code:    code,
		Message: msg[code],
		Data:    data,
	}

	prettyJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}

func ErrorResponse(c fiber.Ctx, code int, data interface{}) error {

	response := fiber.Map{
		"code":    code,
		"message": msg[code],
		"data":    nil,
	}

	prettyJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).Send(prettyJSON)
}
