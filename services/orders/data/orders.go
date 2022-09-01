package data

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Omar-Belghaouti/pdash/services/orders/pb"
	"github.com/go-redis/redis/v9"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ctx        context.Context
	collection *mongo.Collection
	rdb        *redis.Client
)

func init() {
	ctx = context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err.Error())
	}
	collection = client.Database("db").Collection("orders")
	rdb = redis.NewClient(&redis.Options{
		Addr: "redis:6379",
	})
}

// MarshalBinary is a marshalling function for Order
func (order Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(order)
}

// Order struct is a representation of a Order document
type Order struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	SupplierID primitive.ObjectID `bson:"supplier_id" json:"supplier_id"`
	CustomerID primitive.ObjectID `bson:"customer_id" json:"customer_id"`
	TotalPrice float64            `bson:"total_price" json:"total_price"`
	CreatedAt  string             `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt  string             `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Orders is a slice of Order structs
type Orders []Order

// CreateOrder creates a new Order document
func CreateOrder(order Order, grpcCustomerClient pb.CustomerServiceClient, grpcSupplierClient pb.SupplierServiceClient) (Order, int, error) {
	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	order.UpdatedAt = order.CreatedAt
	// check if customer exists
	_, err := grpcCustomerClient.GetCustomer(ctx, &pb.Customer{
		Id: order.CustomerID.Hex(),
	})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.NotFound {
				return order, http.StatusNotFound, err
			} else {
				return order, http.StatusInternalServerError, err
			}
		} else {
			return order, http.StatusInternalServerError, err
		}
	}
	// check if supplier exists
	_, err = grpcSupplierClient.GetSupplier(ctx, &pb.Supplier{
		Id: order.SupplierID.Hex(),
	})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.NotFound {
				return order, http.StatusNotFound, err
			} else {
				return order, http.StatusInternalServerError, err
			}
		} else {
			return order, http.StatusInternalServerError, err
		}
	}
	_, err = collection.InsertOne(ctx, order)
	if err != nil {
		return order, http.StatusInternalServerError, err
	}
	return order, http.StatusCreated, err
}

// GetOrders returns all Orders
func GetOrders() (Orders, int, error) {
	var orders Orders
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return orders, http.StatusInternalServerError, err
	}
	if err := cursor.All(ctx, &orders); err != nil {
		return orders, http.StatusInternalServerError, err
	}
	if orders == nil {
		return Orders{}, http.StatusNotFound, nil
	}
	return orders, http.StatusOK, nil
}

// GetOrdersByCustomerID returns all Orders by Customer ID
func GetOrdersByCustomerID(id string, grpcCustomerClient pb.CustomerServiceClient) (Orders, int, error) {
	// check if customer exists
	_, err := grpcCustomerClient.GetCustomer(ctx, &pb.Customer{
		Id: id,
	})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.NotFound {
				return Orders{}, http.StatusNotFound, err
			} else {
				return Orders{}, http.StatusInternalServerError, err
			}
		} else {
			return Orders{}, http.StatusInternalServerError, err
		}
	}
	var orders Orders
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return orders, http.StatusInternalServerError, err
	}
	cursor, err := collection.Find(ctx, bson.M{"customer_id": oid})
	if err != nil {
		return orders, http.StatusInternalServerError, err
	}
	if err := cursor.All(ctx, &orders); err != nil {
		return orders, http.StatusInternalServerError, err
	}
	if orders == nil {
		return Orders{}, http.StatusNotFound, nil
	}
	return orders, http.StatusOK, nil
}

// GetOrdersBySupplierID returns all Orders by Supplier ID
func GetOrdersBySupplierID(id string, grpcSupplierClient pb.SupplierServiceClient) (Orders, int, error) {
	// check if supplier exists
	_, err := grpcSupplierClient.GetSupplier(ctx, &pb.Supplier{
		Id: id,
	})
	if err != nil {
		s, ok := status.FromError(err)
		if ok {
			if s.Code() == codes.NotFound {
				return Orders{}, http.StatusNotFound, err
			} else {
				return Orders{}, http.StatusInternalServerError, err
			}
		} else {
			return Orders{}, http.StatusInternalServerError, err
		}
	}
	var orders Orders
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return orders, http.StatusInternalServerError, err
	}
	cursor, err := collection.Find(ctx, bson.M{"supplier_id": oid})
	if err != nil {
		return orders, http.StatusInternalServerError, err
	}
	if err := cursor.All(ctx, &orders); err != nil {
		return orders, http.StatusInternalServerError, err
	}
	if orders == nil {
		return Orders{}, http.StatusNotFound, nil
	}
	return orders, http.StatusOK, nil
}

// GetOrder returns a Order by ID
func GetOrder(id string) (Order, int, error) {
	var order Order
	// get order from cache
	val, err := rdb.Get(ctx, id).Result()
	if err == redis.Nil {
		// order not in cache
		objectID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return order, http.StatusBadRequest, err
		}
		err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
		if err != nil {
			return order, http.StatusNotFound, err
		}
		// set order in cache for 5 minutes
		err = rdb.Set(ctx, id, order, time.Minute*5).Err()
		if err != nil {
			return order, http.StatusInternalServerError, err
		}
		return order, http.StatusOK, nil
	} else if err != nil {
		return order, http.StatusInternalServerError, err
	}
	// order in cache
	err = json.Unmarshal([]byte(val), &order)
	if err != nil {
		return order, http.StatusInternalServerError, err
	}
	return order, http.StatusOK, nil
}

// UpdateOrder updates a Order by ID
func UpdateOrder(id string, order Order) (Order, int, error) {
	// check if order exists
	_, status, err := GetOrder(id)
	if err != nil {
		return order, status, err
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return order, http.StatusBadRequest, err
	}
	order.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": order})
	if err != nil {
		return order, http.StatusInternalServerError, err
	}
	// set order in cache for 5 minutes
	err = rdb.Set(ctx, id, order, time.Minute*5).Err()
	if err != nil {
		return order, http.StatusInternalServerError, err
	}
	return order, http.StatusOK, nil
}

// DeleteOrder deletes a Order by ID
func DeleteOrder(id string) (int, error) {
	// check if order exists
	_, status, err := GetOrder(id)
	if err != nil {
		return status, err
	}
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// remove order from cache
	err = rdb.Del(ctx, id).Err()
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}

// GetOrdersLength returns the number of Orders
func GetOrdersLength() int {
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return 0
	}
	orders := Orders{}
	if err := cursor.All(ctx, &orders); err != nil {
		return 0
	}
	return len(orders)
}
