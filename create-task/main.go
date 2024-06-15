package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"serverless-todo-golang/db"
	"serverless-todo-golang/utils/logger"
	"serverless-todo-golang/utils/middleware"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logger.GetLog().Info("INFO : ", "Create Task Handler Called..")

	customCtx, err := middleware.AuthHandler(ctx, &request)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Unauthorized Access.")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized.",
		}, nil
	}

	var req db.Task
	customCtxValue := customCtx.(middleware.CustomContext)
	req.UserId = customCtxValue.UserID

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		logger.GetLog().Error("ERROR : ", "Error in unmarshalling the request body.")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	id := uuid.New().String()[:8]
	req.ID = id
	req.Completed = false
	req.CreatedAt = time.Now()
	req.IsEditing = false

	err = db.InsertTask(req)
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in creating task: %v", err.Error()))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	responseBody, err := json.Marshal(req)
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in marshalling the response: %v", err.Error()))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal server error",
		}, nil
	}

	logger.GetLog().Info("INFO : ", "Task Created Successfully.")
	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(responseBody)},
		nil
}
