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
	collection = client.Database("db").Collection("suppliers")
}

// Supplier struct is a representation of a Supplier document
type Supplier struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	CreatedAt string             `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Suppliers is a slice of Supplier structs
type Suppliers []Supplier

// CreateSupplier creates a new Supplier document
func CreateSupplier(supplier Supplier) (Supplier, int, error) {
	supplier.ID = primitive.NewObjectID()
	supplier.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	supplier.UpdatedAt = supplier.CreatedAt
	_, err := collection.InsertOne(ctx, supplier)
	if err != nil {
		return supplier, http.StatusInternalServerError, err
	}
	return supplier, http.StatusCreated, err
}

// GetSuppliers returns all Suppliers
func GetSuppliers() (Suppliers, int, error) {
	var suppliers Suppliers
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		return suppliers, http.StatusInternalServerError, err
	}
	if err := cursor.All(ctx, &suppliers); err != nil {
		return suppliers, http.StatusInternalServerError, err
	}
	if suppliers == nil {
		return Suppliers{}, http.StatusNotFound, nil
	}
	return suppliers, http.StatusOK, nil
}

// GetSupplier returns a Supplier by ID
func GetSupplier(id string) (Supplier, int, error) {
	var supplier Supplier
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return supplier, http.StatusBadRequest, err
	}
	filter := bson.M{"_id": objectID}
	err = collection.FindOne(ctx, filter).Decode(&supplier)
	if err != nil {
		return supplier, http.StatusNotFound, err
	}
	return supplier, http.StatusOK, nil
}

// UpdateSupplier updates a Supplier by ID
func UpdateSupplier(id string, supplier Supplier) (Supplier, int, error) {
	// check if supplier exists
	_, status, err := GetSupplier(id)
	if err != nil {
		return supplier, status, err
	}
	supplierID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return supplier, http.StatusBadRequest, err
	}
	supplier.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	_, err = collection.UpdateOne(ctx, bson.M{"_id": supplierID}, bson.M{"$set": supplier})
	if err != nil {
		return supplier, http.StatusInternalServerError, err
	}
	return supplier, http.StatusOK, nil
}

// DeleteSupplier deletes a Supplier by ID
func DeleteSupplier(id string) (int, error) {
	// check if supplier exists
	_, status, err := GetSupplier(id)
	if err != nil {
		return status, err
	}
	supplierID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return http.StatusBadRequest, err
	}
	_, err = collection.DeleteOne(ctx, bson.M{"_id": supplierID})
	if err != nil {
		return http.StatusInternalServerError, err
	}
	return http.StatusOK, nil
}
