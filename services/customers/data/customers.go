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
	collection = client.Database("db").Collection("customers")
}

// Customer struct is a representation of a Customer document
type Customer struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt string             `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Customers is a slice of Customer structs
type Customers []Customer

// CreateCustomer creates a new Customer document
func CreateCustomer(customer Customer) (Customer, int, error) {
	customer.ID = primitive.NewObjectID()
	customer.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	customer.UpdatedAt = customer.CreatedAt
	_, err := collection.InsertOne(ctx, customer)
	if err != nil {
		return customer, http.StatusInternalServerError, err
	}
	return customer, http.StatusCreated, err
}

// GetCustomers returns all Customers
func GetCustomers() (Customers, int, error) {
	var customers Customers
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return customers, http.StatusInternalServerError, err
	}
	if err := cursor.All(ctx, &customers); err != nil {
		return customers, http.StatusInternalServerError, err
	}
	if customers == nil {
		return Customers{}, http.StatusNotFound, nil
	}
	return customers, http.StatusOK, err
}

// GetCustomer returns a single Customer
func GetCustomer(id string) (Customer, int, error) {
	var customer Customer
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customer, http.StatusBadRequest, err
	}
	filter := bson.M{"_id": objectID}
	err = collection.FindOne(ctx, filter).Decode(&customer)
	if err != nil {
		return customer, http.StatusNotFound, err
	}
	return customer, http.StatusOK, err
}

// UpdateCustomer updates a single Customer
func UpdateCustomer(id string, customer Customer) (Customer, int, error) {
	// check if supplier exists
	_, status, err := GetCustomer(id)
	if err != nil {
		return customer, status, err
	}
	customerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return customer, http.StatusBadRequest, err
	}
	customer.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": customerID}, bson.M{"$set": customer})
	if err != nil {
		return customer, http.StatusInternalServerError, err
	}
	return customer, http.StatusOK, err
}

// DeleteCustomer deletes a single Customer
func DeleteCustomer(id string) (int, error) {
	// check if customer exists
	_, status, err := GetCustomer(id)
	if err != nil {
		return status, err
	}
	customerID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": customerID})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, err
}
