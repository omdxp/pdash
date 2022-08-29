package data

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/Omar-Belghaouti/pdash/services/auth/util"
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
	collection = client.Database("db").Collection("users")
}

// User struct is a representation of a User document
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username string             `bson:"username" json:"username"`
	Password string             `bson:"password" json:"password"`
	Fullname string             `bson:"fullname" json:"fullname"`
	Email    string             `bson:"email" json:"email"`
}

// CreateUser creates a new user
func CreateUser(user User) (User, int, error) {
	// check if user already exists
	existingUser, _ := GetUserByUsername(user.Username)
	if existingUser.ID != primitive.NilObjectID {
		return user, http.StatusConflict, errors.New("user already exists")
	}
	user.ID = primitive.NewObjectID()
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	user.Password = hashedPassword
	_, err = collection.InsertOne(ctx, user)
	if err != nil {
		return user, http.StatusInternalServerError, err
	}
	return user, http.StatusCreated, nil
}

// GetUserByUsername returns a user by username
func GetUserByUsername(username string) (User, error) {
	var user User
	filter := bson.M{"username": username}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
