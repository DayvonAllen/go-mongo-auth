package repo

import (
	"context"
	"example.com/app/domain"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type AuthRepoImpl struct {
	*domain.User
}

func Login(username, password string) (*domain.User, string, error) {
	var login domain.Authentication
	var user domain.User
	opts := options.FindOne()
	_ = dbConnection.Collection.FindOne(context.TODO(), bson.D{{"USERNAME", username}},opts).Decode(&user)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		panic(err)
	}

	if string(hashedPassword) == user.Password {
		token, _ := login.GenerateJWT(user)
		return &user, token, nil
	}

	return nil, "", fmt.Errorf("error logging in")
}