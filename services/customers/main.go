package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Omar-Belghaouti/pdash/services/customers/data"
	"github.com/Omar-Belghaouti/pdash/services/customers/pb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Response struct {
	Message string `json:"message"`
}

func main() {
	var wg sync.WaitGroup
	grpcAuthClient := make(chan pb.AuthServiceClient)

	wg.Add(3)
	// Start the grpc client to Auth gRPC server
	go func() {
		defer wg.Done()
		for {
			log.Print("Dialing Auth gRPC server on port 4004")
			cc, err := grpc.Dial("auth:4004", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("failed to dial: %s", err.Error())
			}
			grpcAuthClient <- pb.NewAuthServiceClient(cc)
		}
	}()

	// Start the grpc server
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", "0.0.0.0:4001")
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

		// CORS
		app.Use(cors.New(cors.Config{
			AllowOrigins: "*",
			AllowMethods: "GET, POST, PUT, DELETE",
		}))

		// Auth middleware
		app.Use(func(c *fiber.Ctx) error {
			authClient := <-grpcAuthClient
			authHeader := c.Get("Authorization")
			if authHeader == "" {
				return c.Status(http.StatusUnauthorized).JSON(Response{Message: "Unauthorized"})
			}
			fields := strings.Fields(authHeader)
			if len(fields) != 2 || fields[0] != "Bearer" {
				return c.Status(http.StatusUnauthorized).JSON(Response{Message: "Unauthorized"})
			}
			token := fields[1]
			_, err := authClient.VerifyToken(context.Background(), &pb.Auth{AccessToken: token})
			if err != nil {
				s, ok := status.FromError(err)
				if ok {
					if s.Code() == codes.Unauthenticated {
						return c.Status(http.StatusUnauthorized).JSON(Response{Message: "Unauthorized"})
					}
					return c.Status(http.StatusInternalServerError).JSON(Response{Message: "Internal server error: " + err.Error()})
				}
				return c.Status(http.StatusInternalServerError).JSON(Response{Message: "Internal server error: " + err.Error()})
			}
			return c.Next()
		})

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

		app.Listen("0.0.0.0:3001")
	}()

	wg.Wait()
}
