package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"serverless-todo-golang/db"
	"serverless-todo-golang/utils/logger"
	"serverless-todo-golang/utils/middleware"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logger.GetLog().Info("INFO : ", "Get Tasks Handler Called..")

	customCtx, err := middleware.AuthHandler(ctx, &request)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Unauthorized Access.")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized.",
		}, nil
	}

	customCtxValue := customCtx.(middleware.CustomContext)
	userId := customCtxValue.UserID

	tasks, err := db.GetAllTasksWIthUserId(userId)
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in getting all tasks: %v", err.Error()))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	responseBody, err := json.Marshal(tasks)
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in marshalling the response: %v", err.Error()))
		return events.APIGatewayV2HTTPResponse{
			// StatusCode: http.StatusInternalServerError,
			// Body:       "Internal server error",
		}, errors.New("internal server error")
	}

	logger.GetLog().Info("INFO : ", "Tasks Fetched Successfully.")
	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(responseBody)},
		nil
}
