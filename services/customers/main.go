package main

import (
	"net/http"

	"github.com/Omar-Belghaouti/pdash/services/customers/data"
	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	app := fiber.New()

	// Create a new Customer
	app.Post("/customers", func(c *fiber.Ctx) error {
		customer := data.Customer{}
		if err := c.BodyParser(&customer); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{Message: err.Error()})
		}
		customer, status, err := data.CreateCustomer(customer)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(customer)
	})

	// Get all Customers
	app.Get("/customers", func(c *fiber.Ctx) error {
		customers, status, err := data.GetCustomers()
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(customers)
	})

	// Get a Customer by ID
	app.Get("/customers/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		customer, status, err := data.GetCustomer(id)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(customer)
	})

	// Update a Customer by ID
	app.Put("/customers/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		customer := data.Customer{}
		if err := c.BodyParser(&customer); err != nil {
			return c.Status(http.StatusBadRequest).JSON(Response{Message: err.Error()})
		}
		customer, status, err := data.UpdateCustomer(id, customer)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(Response{Message: "Customer updated successfully"})
	})

	// Delete a Customer by ID
	app.Delete("/customers/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		status, err := data.DeleteCustomer(id)
		if err != nil {
			return c.Status(status).JSON(Response{Message: err.Error()})
		}
		return c.Status(status).JSON(Response{Message: "Customer deleted successfully"})
	})

	app.Listen(":3001")
}
