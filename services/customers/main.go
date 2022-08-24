package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/Omar-Belghaouti/pdash/pb"
	"github.com/Omar-Belghaouti/pdash/services/customers/data"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	// Start the grpc server
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", ":4001")
		if err != nil {
			log.Fatalf("failed to listen: %s", err.Error())
		}
		defer lis.Close()

		s := grpc.NewServer()
		pb.RegisterCustomerServiceServer(s, &server{})
		reflection.Register(s)

		log.Print("Starting Customer gRPC server on port 4001")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err.Error())
		}
	}()

	// Start the http server
	go func() {
		defer wg.Done()
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
	}()

	wg.Wait()
}
