package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"serverless-todo-golang/db"
	"serverless-todo-golang/utils/crypto"
	"serverless-todo-golang/utils/logger"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logger.GetLog().Info("INFO : ", "Sign Up Handler Called..")

	var req db.User
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		logger.GetLog().Error("ERROR : ", "Error in unmarshalling the request body.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       err.Error(),
		}, err
	}

	var IsValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,4}$`).MatchString
	if !IsValidEmail(req.Email) {
		logger.GetLog().Error("ERROR : ", "Invalid Email Format.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       "Invalid Email.",
		}, errors.New("Invalid Email.")
	}

	ok, err := db.CheckIfUserWithThisEmailExists(req.Email)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in finding the user with email.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       err.Error(),
		}, errors.New("Error in finding the user with email")
	}

	if ok {
		logger.GetLog().Error("ERROR : ", "User with this email already exists.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       "User with this email already exists.",
		}, errors.New("User with this email already exists.")
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Failed to hash password.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       "Failed to hash password.",
		}, errors.New("Failed to hash password.")
	}

	id := uuid.New().String()[:8]
	tokenData := crypto.UserTokenData{
		ID: id,
	}

	token, err := crypto.GenerateAuthToken(tokenData)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in generating the token.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       err.Error(),
		}, errors.New("Error in generating the token")
	}

	user := db.User{
		ID:        id,
		Email:     req.Email,
		Password:  string(hashPassword),
		Name:      req.Name,
		CreatedAt: time.Now(),
	}

	err = db.InsertUser(user)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in inserting the user.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       err.Error(),
		}, errors.New("Error in inserting the user")
	}

	data := map[string]interface{}{
		"token": token,
		"data":  user,
	}

	responseBody, err := json.Marshal(data)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in marshalling the response body.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusInternalServerError,
			// Body:       "Internal server error ----->",
		}, errors.New("Error in marshalling the response body.")
	}

	logger.GetLog().Info("INFO : ", "User Signed Up Successfully.")
	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(responseBody)},
		nil
}
