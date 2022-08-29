package data

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/Omar-Belghaouti/pdash/services/auth/token"
	"github.com/Omar-Belghaouti/pdash/services/auth/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	ctx        context.Context
	collection *mongo.Collection
	config     util.Config
	tokenMaker *token.PasetoMaker
)

func init() {
	var err error
	config, err = util.LoadConfig(".")
	if err != nil {
		log.Fatalf("cannot load config: %s", err.Error())
	}
	tokenMaker, err = token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatalf("cannot create token maker: %s", err.Error())
	}
	ctx = context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %s", err.Error())
	}
	collection = client.Database("db").Collection("users")
}

// User struct is a representation of a User document
type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	Fullname  string             `bson:"fullname" json:"fullname"`
	Email     string             `bson:"email" json:"email"`
	CreatedAt string             `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt string             `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// LoginUserRequest is the request body for the LoginUser endpoint
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUserResponse is the response of the LoginUser function
type LoginUserResponse struct {
	AccessToken string `json:"access_token"`
	User        User   `json:"user"`
}

// CreateUser creates a new user
func CreateUser(user User) (User, int, error) {
	// check if user already exists
	existingUser, _ := getUserByUsername(user.Username)
	if existingUser.ID != primitive.NilObjectID {
		return user, http.StatusConflict, errors.New("user already exists")
	}
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	user.UpdatedAt = user.CreatedAt
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

// LoginUser logs in a user
func LoginUser(req LoginUserRequest) (LoginUserResponse, int, error) {
	user, err := getUserByUsername(req.Username)
	if err != nil {
		return LoginUserResponse{}, http.StatusNotFound, err
	}
	if err := util.CheckPassword(req.Password, user.Password); err != nil {
		return LoginUserResponse{}, http.StatusUnauthorized, err
	}
	accessToken, err := tokenMaker.CreateToken(user.Username, config.AccessTokenDuration)
	if err != nil {
		return LoginUserResponse{}, http.StatusInternalServerError, err
	}
	res := LoginUserResponse{
		AccessToken: accessToken,
		User:        user,
	}
	return res, http.StatusOK, nil
}

// getUserByUsername returns a user by username
func getUserByUsername(username string) (User, error) {
	var user User
	filter := bson.M{"username": username}
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return user, err
	}
	return user, nil
}
