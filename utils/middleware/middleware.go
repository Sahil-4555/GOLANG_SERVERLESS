package middleware

import (
	"context"
	"errors"
	"fmt"
	"log"

	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt"
)

var JWT_API_KEY = "IAMSAHIL"

type CustomContext struct {
	context.Context
	UserID string
}

func ValidateToken(t string, k string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			log.Fatalf("unexpected signin method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(k), nil
	})
	if err != nil {
		return nil, err
	}

	return token, err
}

func AuthHandler(ctx context.Context, request *events.APIGatewayV2HTTPRequest) (context.Context, error) {
	bearerToken := request.Headers["authorization"]

	fmt.Println("Bearer Token--->: ", bearerToken)
	if !strings.HasPrefix(bearerToken, "Bearer ") {
		fmt.Println("Bearer Token token missing or invalid")
		return ctx, errors.New("authorization token missing or invalid")
	}

	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")

	fmt.Println("Token String--->: ", tokenString)
	token, err := ValidateToken(tokenString, JWT_API_KEY)
	if err != nil {
		fmt.Println("Bearer Token invalid")
		return ctx, errors.New("authorization token invalid")
	}

	claims := token.Claims.(jwt.MapClaims)

	userId := claims["userData"].(map[string]interface{})["ID"].(string)
	customCtx := CustomContext{
		Context: ctx,
		UserID:  userId,
	}
	fmt.Println("User ID--->: ", userId)
	return customCtx, nil

}

func GenerateToken(userData interface{}) (string, error) {
	token := jwt.New(jwt.SigningMethodHS512)
	claims := make(jwt.MapClaims)
	claims["userData"] = userData
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(JWT_API_KEY))
	return tokenString, err
}
