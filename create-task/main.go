package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"serverless-todo-golang/db"
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

	customCtx, err := middleware.AuthHandler(ctx, &request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized.",
		}, nil
	}

	var req db.Task
	customCtxValue := customCtx.(middleware.CustomContext)
	req.UserId = customCtxValue.UserID

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		fmt.Println("Error in unmarshalling the request body: ", err)
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
		fmt.Println("Error in inserting the task: --->", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	responseBody, err := json.Marshal(req)
	if err != nil {
		fmt.Println("Error in marshalling the response body: --->", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Internal server error ----->",
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       string(responseBody)},
		nil
}
