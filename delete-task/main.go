package main

import (
	"context"
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

const pathParameterName = "id"

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {

	_, err := middleware.AuthHandler(ctx, &request)
	if err != nil {
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized.",
		}, nil
	}

	taskId := request.PathParameters[pathParameterName]

	err = db.DeleteTaskById(taskId)
	if err != nil {
		fmt.Println("Error in deleting the task: --->", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       fmt.Sprintf("Deleted Successfully %v", taskId)},
		nil
}
