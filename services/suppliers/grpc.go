package main

import (
	"context"
	"net/http"

	"github.com/Omar-Belghaouti/pdash/services/suppliers/data"
	"github.com/Omar-Belghaouti/pdash/services/suppliers/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	pb.SupplierServiceServer
}

// GetSupplier implementation for Supplier gRPC server
func (s *server) GetSupplier(ctx context.Context, in *pb.Supplier) (*pb.Supplier, error) {
	supplier, sc, err := data.GetSupplier(in.Id)
	if err != nil {
		if sc == http.StatusNotFound {
			return nil, status.Errorf(codes.NotFound, "Supplier not found: %s", err.Error())
		}
		return nil, status.Errorf(codes.Internal, "Internal error: %s", err.Error())
	}
	res := &pb.Supplier{
		Id:        supplier.ID.Hex(),
		Name:      supplier.Name,
		CreatedAt: supplier.CreatedAt,
		UpdatedAt: supplier.UpdatedAt,
	}
	return res, nil
}
