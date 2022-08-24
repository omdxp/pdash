package main

import (
	"net/http"

	"github.com/Omar-Belghaouti/pdash/services/suppliers/data"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	app := fiber.New()

	// Create a new Supplier
	app.Post("/suppliers", func(c *fiber.Ctx) error {
		supplier := data.Supplier{}
		if err := c.BodyParser(&supplier); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{Message: err.Error()})
		}
		supplier, status, err := data.CreateSupplier(supplier)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(supplier)
	})

	// Get all Suppliers
	app.Get("/suppliers", func(c *fiber.Ctx) error {
		suppliers, status, err := data.GetSuppliers()
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(suppliers)
	})

	// Get a Supplier by ID
	app.Get("/suppliers/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		supplier, status, err := data.GetSupplier(id)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(supplier)
	})

	// Update a Supplier by ID
	app.Put("/suppliers/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		supplier := data.Supplier{}
		if err := c.BodyParser(&supplier); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{Message: err.Error()})
		}
		supplier, status, err := data.UpdateSupplier(id, supplier)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(Response{Message: "Supplier updated successfully"})
	})

	// Delete a Supplier by ID
	app.Delete("/suppliers/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		status, err := data.DeleteSupplier(id)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(Response{Message: "Supplier deleted successfully"})
	})

	app.Listen(":3003")
}
