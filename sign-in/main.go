package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"serverless-todo-golang/db"
	"serverless-todo-golang/utils/crypto"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	var req db.User
	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		fmt.Println("Error in unmarshalling the request body: ", err)
		return events.APIGatewayV2HTTPResponse{}, err
	}

	fmt.Println("Request Body: --->", req)

	var IsValidEmail = regexp.MustCompile(`^[a-zA-Z0-9._-]+@[a-zA-Z.-]+\.[a-zA-Z]{2,4}$`).MatchString
	if !IsValidEmail(req.Email) {
		fmt.Println("Invalid Email. --->", req.Email)
		return events.APIGatewayV2HTTPResponse{}, errors.New("Invalid Email.")
	}

	user, err := db.FindTheUserWithEmail(req.Email)
	if err != nil {
		fmt.Println("Error in finding the user with email: --->", err)
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusBadRequest,
			// Body:       err.Error(),
		}, errors.New("Error in finding the user with email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		fmt.Println("Email Or Password Not Matched --->", err)
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
		fmt.Println("Error in generating the token: --->", err)
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
		fmt.Println("Error in marshalling the response body: --->", err)
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusInternalServerError,
			// Body:       "Internal server error",
		}, errors.New("Internal server error")
	}

	fmt.Println("User logged in successfully. --->", string(responseBody))
	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(responseBody)},
		nil
}
