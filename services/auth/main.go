package main

import (
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/Omar-Belghaouti/pdash/services/auth/data"
	_ "github.com/Omar-Belghaouti/pdash/services/auth/docs"
	"github.com/Omar-Belghaouti/pdash/services/auth/pb"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/grpc"
)

type Respone struct {
	Message string `json:"message"`
}

// @title pdash auth service
// @version 1.0
// @description pdash auth service
// @contact.name Omar Belghaouti
// @contact.email omarbelghaouti@gmail.com
// @host localhost:8004
// @BasePath /
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

		// Swagger
		app.Get("/swagger/*", swagger.HandlerDefault)

		// Create a new user
		app.Post("/users", CreateUser)

		// Login a user
		app.Post("/users/login", LoginUser)

		app.Listen("0.0.0.0:3004")
	}()

	wg.Wait()
}

// CreateUser creates a new user
// @Summary Create a new user
// @Description Create a new user
// @ID create-user
// @Accept  json
// @Produce  json
// @Param user body data.User true "User"
// @Success 201 {object} data.User
// @Failure 400 {object} Respone
// @Failure 500 {object} Respone
// @Router /users [post]
func CreateUser(c *fiber.Ctx) error {
	user := data.User{}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Respone{Message: err.Error()})
	}
	user, status, err := data.CreateUser(user)
	if err != nil {
		return c.Status(status).JSON(Respone{Message: err.Error()})
	}
	return c.Status(status).JSON(user)
}

// LoginUser logs in a user
// @Summary Login a user
// @Description Login a user
// @ID login-user
// @Accept  json
// @Produce  json
// @Param user body data.LoginUserRequest true "User"
// @Success 200 {object} data.LoginUserResponse
// @Failure 400 {object} Respone
// @Failure 500 {object} Respone
// @Router /users/login [post]
func LoginUser(c *fiber.Ctx) error {
	req := data.LoginUserRequest{}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Respone{Message: err.Error()})
	}
	user, status, err := data.LoginUser(req)
	if err != nil {
		return c.Status(status).JSON(Respone{Message: err.Error()})
	}
	return c.Status(status).JSON(user)
}
