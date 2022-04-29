package main

import (
	"calculator/service"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	app := fiber.New()
	app.Use(logger.New())

	app.Static("/", "./wwwroot", fiber.Static{
		Index:         "index.html",
		CacheDuration: time.Second * 10,
	})

	app.Post("/calculate", Calculate)

	app.Listen(":3000")

}

func Calculate(c *fiber.Ctx) error {
	request := c.Body()

	tokens, err := service.SplitExpression(string(request))
	if err != nil {
		return fiber.ErrBadRequest
	}

	result, err := service.Evaluate(tokens)
	if err != nil {
		return fiber.ErrBadRequest
	}

	return c.SendString(result)
}
