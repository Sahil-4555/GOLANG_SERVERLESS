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

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logger.GetLog().Info("INFO : ", "Sign In Handler Called..")

	var req db.User
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		logger.GetLog().Error("ERROR : ", "Error in unmarshalling the request body.")
		return events.APIGatewayV2HTTPResponse{}, err
	}

	var IsValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,4}$`).MatchString
	if !IsValidEmail(req.Email) {
		logger.GetLog().Error("ERROR : ", "Invalid Email Format.")
		return events.APIGatewayV2HTTPResponse{}, errors.New("Invalid Email.")
	}

	user, err := db.FindTheUserWithEmail(req.Email)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in finding the user with email.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       err.Error(),
		}, errors.New("Error in finding the user with email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Email Or Password Not Matched")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       "Email Or Password Not Matched",
		}, errors.New("Email Or Password Not Matched")
	}

	tokenData := crypto.UserTokenData{
		ID: user.ID,
	}

	token, err := crypto.GenerateAuthToken(tokenData)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in generating the token.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       err.Error(),
		}, err
	}

	loginData := db.User{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	data := map[string]interface{}{
		"token": token,
		"data":  loginData,
	}

	responseBody, err := json.Marshal(data)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in marshalling the response body.")
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusInternalServerError,
			// Body:       "Internal server error",
		}, errors.New("Internal server error")
	}

	logger.GetLog().Info("INFO : ", "User Logged In Successfully.")
	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(responseBody)},
		nil
}
