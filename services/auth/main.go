package main

import (
	"net/http"

	"github.com/Omar-Belghaouti/pdash/services/auth/data"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Respone struct {
	Message string `json:"message"`
}

func main() {
	app := fiber.New()

	// CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))

	// Create a new user
	app.Post("/users", func(c *fiber.Ctx) error {
		user := data.User{}
		if err := c.BodyParser(&user); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Respone{Message: err.Error()})
		}
		user, status, err := data.CreateUser(user)
		if err != nil {
			return c.Status(status).JSON(Respone{Message: err.Error()})
		}
		return c.Status(status).JSON(user)
	})

	app.Listen("0.0.0.0:3004")
}
