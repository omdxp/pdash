package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/Omar-Belghaouti/pdash/services/suppliers/data"
	"github.com/Omar-Belghaouti/pdash/services/suppliers/pb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
		lis, err := net.Listen("tcp", ":4003")
		if err != nil {
			log.Fatalf("failed to listen: %s", err.Error())
		}
		defer lis.Close()

		s := grpc.NewServer()
		pb.RegisterSupplierServiceServer(s, &server{})
		reflection.Register(s)

		log.Print("Starting Supplier gRPC server on port 4003")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %s", err.Error())
		}
	}()

	// Start the http server
	go func() {
		defer wg.Done()
		app := fiber.New()

		// CORS
		app.Use(cors.New())

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

		app.Listen("localhost:3003")
	}()

	wg.Wait()
}
