package main

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/Omar-Belghaouti/pdash/pb"
	"github.com/Omar-Belghaouti/pdash/services/orders/data"
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Respone struct {
	Message string `json:"message"`
}

func main() {
	var wg sync.WaitGroup
	grpcCustomerClient := make(chan pb.CustomerServiceClient)
	grpcSupplierClient := make(chan pb.SupplierServiceClient)

	wg.Add(3)
	// Start the grpc client to Customers gRPC server
	go func() {
		defer wg.Done()
		for {
			log.Print("Dialing Customers gRPC server on port 4001")
			cc, err := grpc.Dial("localhost:4001", grpc.WithTransportCredentials(insecure.NewCredentials()))
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
			log.Print("Dialing Suppliers gRPC server on port 4001")
			cc, err := grpc.Dial("localhost:4003", grpc.WithTransportCredentials(insecure.NewCredentials()))
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

		// Create a new Order
		app.Post("/orders", func(c *fiber.Ctx) error {
			order := data.Order{}
			if err := c.BodyParser(&order); err != nil {
				return c.Status(http.StatusBadRequest).JSON(Respone{Message: err.Error()})
			}
			order, status, err := data.CreateOrder(order, <-grpcCustomerClient, <-grpcSupplierClient)
			if err != nil {
				return c.Status(status).JSON(Respone{Message: err.Error()})
			}
			return c.Status(status).JSON(order)
		})

		// Get all Orders
		app.Get("/orders", func(c *fiber.Ctx) error {
			supplierID := c.Query("supplier_id")
			customerID := c.Query("customer_id")
			if strings.TrimSpace(supplierID) != "" && strings.TrimSpace(customerID) != "" {
				return c.Status(http.StatusBadRequest).JSON(Respone{Message: "supplier_id and customer_id are mutually exclusive"})
			}
			if strings.TrimSpace(supplierID) != "" {
				orders, status, err := data.GetOrdersBySupplierID(supplierID, <-grpcSupplierClient)
				if err != nil {
					return c.Status(status).JSON(Respone{Message: err.Error()})
				}
				return c.Status(status).JSON(orders)
			}
			if strings.TrimSpace(customerID) != "" {
				orders, status, err := data.GetOrdersByCustomerID(customerID, <-grpcCustomerClient)
				if err != nil {
					return c.Status(status).JSON(Respone{Message: err.Error()})
				}
				return c.Status(status).JSON(orders)
			}
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
	}()

	wg.Wait()
}
