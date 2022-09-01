package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Omar-Belghaouti/pdash/services/orders/data"
	_ "github.com/Omar-Belghaouti/pdash/services/orders/docs"
	"github.com/Omar-Belghaouti/pdash/services/orders/pb"
	"github.com/antoniodipinto/ikisocket"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/websocket/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Response struct {
	Message string `json:"message"`
}

type OrdersData struct {
	Length int `json:"length"`
}

type EventMessage struct {
	Event string     `json:"event"`
	Data  OrdersData `json:"data"`
}

// @title pdash orders service
// @version 1.0
// @description pdash orders service
// @contact.name Omar Belghaouti
// @contact.email omarbelghaouti@gmail.com
// @host localhost:8002
// @BasePath /
func main() {
	var wg sync.WaitGroup
	grpcCustomerClient := make(chan pb.CustomerServiceClient)
	grpcSupplierClient := make(chan pb.SupplierServiceClient)
	grpcAuthClient := make(chan pb.AuthServiceClient)

	wg.Add(4)
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

	// Start the grpc client to Customers gRPC server
	go func() {
		defer wg.Done()
		for {
			log.Print("Dialing Customers gRPC server on port 4001")
			cc, err := grpc.Dial("customers:4001", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("failed to dial: %s", err.Error())
			}
			grpcCustomerClient <- pb.NewCustomerServiceClient(cc)
		}
	}()

	// Start the grpc client to Suppliers gRPC server
	go func() {
		defer wg.Done()
		for {
			log.Print("Dialing Suppliers gRPC server on port 4003")
			cc, err := grpc.Dial("suppliers:4003", grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatalf("failed to dial: %s", err.Error())
			}
			grpcSupplierClient <- pb.NewSupplierServiceClient(cc)
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

		// Swagger
		app.Get("/swagger/*", swagger.HandlerDefault)

		// Setup websocket
		app.Use("/ws", func(c *fiber.Ctx) error {
			if websocket.IsWebSocketUpgrade(c) {
				c.Locals("allowed", true)
				return c.Next()
			}
			return fiber.ErrUpgradeRequired
		})

		// Websocket
		app.Get("/ws", ikisocket.New(func(kws *ikisocket.Websocket) {}))

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

		// Create a new Order
		app.Post("/orders", func(c *fiber.Ctx) error {
			order := data.Order{}
			if err := c.BodyParser(&order); err != nil {
				return c.Status(http.StatusBadRequest).JSON(Response{Message: err.Error()})
			}
			order, status, err := data.CreateOrder(order, <-grpcCustomerClient, <-grpcSupplierClient)
			if err != nil {
				return c.Status(status).JSON(Response{Message: err.Error()})
			}
			b, _ := json.Marshal(EventMessage{
				Event: "orders",
				Data: OrdersData{
					Length: data.GetOrdersLength(),
				},
			})
			ikisocket.Broadcast(b)
			return c.Status(status).JSON(order)
		})

		// Get all Orders
		app.Get("/orders", func(c *fiber.Ctx) error {
			supplierID := c.Query("supplier_id")
			customerID := c.Query("customer_id")
			if strings.TrimSpace(supplierID) != "" && strings.TrimSpace(customerID) != "" {
				return c.Status(http.StatusBadRequest).JSON(Response{Message: "supplier_id and customer_id are mutually exclusive"})
			}
			if strings.TrimSpace(supplierID) != "" {
				orders, status, err := data.GetOrdersBySupplierID(supplierID, <-grpcSupplierClient)
				if err != nil {
					return c.Status(status).JSON(Response{Message: err.Error()})
				}
				return c.Status(status).JSON(orders)
			}
			if strings.TrimSpace(customerID) != "" {
				orders, status, err := data.GetOrdersByCustomerID(customerID, <-grpcCustomerClient)
				if err != nil {
					return c.Status(status).JSON(Response{Message: err.Error()})
				}
				return c.Status(status).JSON(orders)
			}
			orders, status, err := data.GetOrders()
			if err != nil {
				return c.Status(status).JSON(Response{Message: err.Error()})
			}
			return c.Status(status).JSON(orders)
		})

		// Get a Order by ID
		app.Get("/orders/:id", GetOrderByID)

		// Update a Order by ID
		app.Put("/orders/:id", UpdateOrderByID)

		// Delete a Order by ID
		app.Delete("/orders/:id", DeleteOrderByID)

		app.Listen("0.0.0.0:3002")
	}()

	wg.Wait()
}

// CreateOrder creates a new Order
// @Summary Create a new Order
// @Description Create a new Order
// @ID create-order
// @Accept  json
// @Produce  json
// @Param order body data.Order true "Order"
// @Success 201 {object} data.Order
// @Failure 401 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /orders [post]
func CreateOrder() {}

// GetOrders returns all Orders
// @Summary Get all Orders
// @Description Get all Orders
// @ID get-orders
// @Accept  json
// @Produce  json
// @Param supplier_id query string false "Supplier ID"
// @Param customer_id query string false "Customer ID"
// @Success 200 {array} data.Order
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /orders [get]
func GetOrders() {}

// GetOrderByID returns a Order by ID
// @Summary Get a Order by ID
// @Description Get a Order by ID
// @ID get-order-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} data.Order
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /orders/{id} [get]
func GetOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order, status, err := data.GetOrder(id)
	if err != nil {
		return c.Status(status).JSON(Response{Message: err.Error()})
	}
	return c.Status(status).JSON(order)
}

// UpdateOrderByID updates a Order by ID
// @Summary Update a Order by ID
// @Description Update a Order by ID
// @ID update-order-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Param order body data.Order true "Order"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /orders/{id} [put]
func UpdateOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	order := data.Order{}
	if err := c.BodyParser(&order); err != nil {
		return c.Status(http.StatusBadRequest).JSON(Response{Message: err.Error()})
	}
	order, status, err := data.UpdateOrder(id, order)
	if err != nil {
		return c.Status(status).JSON(Response{Message: err.Error()})
	}
	return c.Status(status).JSON(Response{Message: "Order updated successfully"})
}

// DeleteOrderByID deletes a Order by ID
// @Summary Delete a Order by ID
// @Description Delete a Order by ID
// @ID delete-order-by-id
// @Accept  json
// @Produce  json
// @Param id path string true "Order ID"
// @Success 200 {object} Response
// @Failure 401 {object} Response
// @Failure 404 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /orders/{id} [delete]
func DeleteOrderByID(c *fiber.Ctx) error {
	id := c.Params("id")
	status, err := data.DeleteOrder(id)
	if err != nil {
		return c.Status(status).JSON(Response{Message: err.Error()})
	}
	b, _ := json.Marshal(EventMessage{
		Event: "orders",
		Data: OrdersData{
			Length: data.GetOrdersLength(),
		},
	})
	ikisocket.Broadcast(b)
	return c.Status(status).JSON(Response{
		Message: "Order deleted",
	})
}
