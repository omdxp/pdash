package main

import (
	"context"
	"net/http"

	"github.com/Omar-Belghaouti/pdash/services/customers/data"
	"github.com/Omar-Belghaouti/pdash/services/customers/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.CustomerServiceServer
}

// GetCustomer implementation for Customer gRPC server
func (s *server) GetCustomer(ctx context.Context, in *pb.Customer) (*pb.Customer, error) {
	customer, sc, err := data.GetCustomer(in.Id)
	if err != nil {
		if sc == http.StatusNotFound {
			return nil, status.Errorf(codes.NotFound, "Customer not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal error: %s", err.Error())
	}
	res := &pb.Customer{
		Id:        customer.ID.Hex(),
		Name:      customer.Name,
		CreatedAt: customer.CreatedAt,
		UpdatedAt: customer.UpdatedAt,
	}
	return res, nil
}
