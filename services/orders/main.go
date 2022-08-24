package main

import (
	"net/http"

	"github.com/Omar-Belghaouti/pdash/services/orders/data"
	"github.com/gofiber/fiber/v2"
)

type Respone struct {
	Message string `json:"message"`
}

func main() {
	app := fiber.New()

	// Create a new Order
	app.Post("/orders", func(c *fiber.Ctx) error {
		order := data.Order{}
		if err := c.BodyParser(&order); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Respone{Message: err.Error()})
		}
		order, status, err := data.CreateOrder(order)
		if err != nil {
			return c.Status(status).JSON(Respone{Message: err.Error()})
		}
		return c.Status(status).JSON(order)
	})

	// Get all Orders
	app.Get("/orders", func(c *fiber.Ctx) error {
		orders, status, err := data.GetOrders()
		if err != nil {
			return c.Status(status).JSON(Respone{Message: err.Error()})
		}
		return c.Status(status).JSON(orders)
	})

	// Get a Order by ID
	app.Get("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		order, status, err := data.GetOrder(id)
		if err != nil {
			return c.Status(status).JSON(Respone{Message: err.Error()})
		}
		return c.Status(status).JSON(order)
	})

	// Update a Order by ID
	app.Put("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		order := data.Order{}
		if err := c.BodyParser(&order); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Respone{Message: err.Error()})
		}
		order, status, err := data.UpdateOrder(id, order)
		if err != nil {
			return c.Status(status).JSON(Respone{Message: err.Error()})
		}
		return c.Status(status).JSON(Respone{Message: "Order updated successfully"})
	})

	// Delete a Order by ID
	app.Delete("/orders/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		status, err := data.DeleteOrder(id)
		if err != nil {
			return c.Status(status).JSON(Respone{Message: err.Error()})
		}
		return c.Status(status).JSON(Respone{
			Message: "Order deleted",
		})
	})

	app.Listen(":3002")
}
