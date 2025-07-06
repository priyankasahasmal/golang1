package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func main() {
	fmt.Println("hello world")
	app := fiber.New()
	app.Get("/abc", func(c *fiber.Ctx) error {
		return c.SendString("hello world!")
	})
	app.Listen(":3000")
}