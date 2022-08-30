package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/Omar-Belghaouti/pdash/services/auth/data"
	"github.com/Omar-Belghaouti/pdash/services/auth/pb"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
)

type Respone struct {
	Message string `json:"message"`
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	// Start the gRPC server
	go func() {
		defer wg.Done()
		lis, err := net.Listen("tcp", "0.0.0.0:4004")
		if err != nil {
			log.Fatalf("failed to listen: %s", err.Error())
		}
		defer lis.Close()

		s := grpc.NewServer()
		pb.RegisterAuthServiceServer(s, &server{})
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
			AllowMethods: "POST",
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

		// Login a user
		app.Post("/users/login", func(c *fiber.Ctx) error {
			req := data.LoginUserRequest{}
			if err := c.BodyParser(&req); err != nil {
				return c.Status(http.StatusBadRequest).JSON(Respone{Message: err.Error()})
			}
			user, status, err := data.LoginUser(req)
			if err != nil {
				return c.Status(status).JSON(Respone{Message: err.Error()})
			}
			return c.Status(status).JSON(user)
		})

		app.Listen("0.0.0.0:3004")
	}()

	wg.Wait()
}
