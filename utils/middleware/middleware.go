package middleware

import (
	"context"
	"errors"
	"fmt"
	"serverless-todo-golang/utils/logger"

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
			logger.GetLog().Error("ERROR : ", fmt.Sprintf("unexpected signin method: %v", token.Header["alg"]))
			return nil, fmt.Errorf("unexpected signin method: %v", token.Header["alg"])
		}

		return []byte(k), nil
	})
	if err != nil {
		logger.GetLog().Error("ERROR : ", err.Error())
		return nil, err
	}

	return token, err
}

func AuthHandler(ctx context.Context, request *events.APIGatewayV2HTTPRequest) (context.Context, error) {
	logger.GetLog().Info("INFO : ", "AuthHandler Called..")
	bearerToken := request.Headers["authorization"]

	if !strings.HasPrefix(bearerToken, "Bearer ") {
		logger.GetLog().Error("ERROR : ", "Authorization token missing or invalid")
		return ctx, errors.New("authorization token missing or invalid")
	}

	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")

	token, err := ValidateToken(tokenString, JWT_API_KEY)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Authorization token invalid")
		return ctx, errors.New("authorization token invalid")
	}

	claims := token.Claims.(jwt.MapClaims)

	userId := claims["userData"].(map[string]interface{})["ID"].(string)
	customCtx := CustomContext{
		Context: ctx,
		UserID:  userId,
	}

	return customCtx, nil

}

func GenerateToken(userData interface{}) (string, error) {
	logger.GetLog().Info("INFO : ", "GenerateToken Called..")
	token := jwt.New(jwt.SigningMethodHS512)
	claims := make(jwt.MapClaims)
	claims["userData"] = userData
	claims["exp"] = time.Now().Add(time.Hour * 24 * 30).Unix()
	token.Claims = claims
	tokenString, err := token.SignedString([]byte(JWT_API_KEY))
	return tokenString, err
}
