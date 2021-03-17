package domain

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"strings"
	"time"
)

type Authentication struct {
	Id primitive.ObjectID
	Email string
}

type LoginDetails struct {
	Email string
	Password string
}

type Claims struct {
	jwt.StandardClaims
	Id primitive.ObjectID
	Email string
}

var newKey = make([]byte, 64)
var _, _ = io.ReadFull(rand.Reader, newKey)


func (l Authentication) GenerateJWT(msg User) (string, error){

	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(5 * time.Minute).Unix(),
		},
		Id: msg.Id,
		Email: msg.Email,
	}
	// always better to use a pointer with JSON
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	signedString, err := token.SignedString(newKey)

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}
	return signedString, nil
}

func (l Authentication) SignToken(token []byte) ([]byte, error) {
	// second arg is a private key, key needs to be the same size as hasher
	// sha512 is 64 bits
	h := hmac.New(sha256.New, newKey)

	// hash is a writer
	_, err := h.Write(token)
	if err != nil {
		return nil, fmt.Errorf("error in signMessage while hashing message: %w", err)
	}

	// returns signature value
	signature := h.Sum(nil)

	return signature, nil
}

func (l Authentication) VerifySignature(token, sig []byte) (bool, error) {
	// sign message
	newSig, err := l.SignToken(token)

	if err != nil {
		return false, fmt.Errorf("error verifying signature: %w", err)
	}

	// compare it
	return hmac.Equal(newSig, sig), nil
}

func(l Authentication) IsLoggedIn(cookie string) (*Authentication, error)  {

	xs := strings.SplitN(cookie, " ", 2)

	tokenValue := strings.SplitN(xs[1], "|", 2)

	validSig, err := l.VerifySignature([]byte(tokenValue[0]), []byte(tokenValue[1]))
	if err != nil {
		return nil, fmt.Errorf("error...: %w", err)
	}

	if !validSig {
		return nil, fmt.Errorf("invalid signature: %w", err)
	}

	var jwtValue string

	if len(xs) == 2 {
		jwtValue = xs[1]
	}

	token, err := jwt.ParseWithClaims(jwtValue, &Claims{},func(t *jwt.Token)(interface{}, error) {
		if t.Method.Alg() == jwt.SigningMethodHS256.Alg() {
			//verify token(we pass in our key to be verified)
			return newKey, nil
		}
		return nil, fmt.Errorf("invalid signing method")
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	isEqual := token.Valid

	if isEqual {
		// user is logged in at this point
		// because we receive an interface type we need to assert which type we want to use that inherits it
		claims := token.Claims.(*Claims)

		l.Id = claims.Id
		l.Email = claims.Email
		return &l, err
	}

	return nil, fmt.Errorf("token is not valid")
}