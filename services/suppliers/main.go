package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/Omar-Belghaouti/pdash/services/suppliers/data"
	"github.com/Omar-Belghaouti/pdash/services/suppliers/pb"
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
		lis, err := net.Listen("tcp", "0.0.0.0:4003")
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

		app.Listen("0.0.0.0:3003")
	}()

	wg.Wait()
}
