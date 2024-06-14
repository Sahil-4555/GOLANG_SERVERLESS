package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"serverless-todo-golang/db"
	"serverless-todo-golang/utils/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	customCtx, err := middleware.AuthHandler(ctx, &request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized.",
		}, nil
	}

	customCtxValue := customCtx.(middleware.CustomContext)
	userId := customCtxValue.UserID

	tasks, err := db.GetAllTasksWIthUserId(userId)
	if err != nil {
		fmt.Println("Error in updating the task: --->", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	responseBody, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Error in marshalling the response body: --->", err)
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusInternalServerError,
			// Body:       "Internal server error",
		}, errors.New("internal server error")
	}

	fmt.Println("getted all tasks succesfully --->", string(responseBody))
	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(responseBody)},
		nil
}
