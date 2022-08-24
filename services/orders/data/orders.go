package data

import (
	"context"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        context.Context
	collection *mongo.Collection
)

func init() {
	ctx = context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err.Error())
	}
	collection = client.Database("db").Collection("orders")
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
func CreateOrder(order Order) (Order, int, error) {
	order.ID = primitive.NewObjectID()
	order.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	order.UpdatedAt = order.CreatedAt
	_, err := collection.InsertOne(ctx, order)
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

// GetOrder returns a Order by ID
func GetOrder(id string) (Order, int, error) {
	var order Order
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return order, http.StatusBadRequest, err
	}
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&order)
	if err != nil {
		return order, http.StatusNotFound, err
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
	return http.StatusOK, nil
}