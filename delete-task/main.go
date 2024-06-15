package main

import (
	"context"
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

const pathParameterName = "id"

func handler(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	logger.GetLog().Info("INFO : ", "Delete Task Handler Called..")

	_, err := middleware.AuthHandler(ctx, &request)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Unauthorized Access.")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       "Unauthorized.",
		}, nil
	}

	taskId := request.PathParameters[pathParameterName]

	err = db.DeleteTaskById(taskId)
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in deleting task: %v", err.Error()))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	logger.GetLog().Info("INFO : ", "Task Deleted Successfully.")
	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       fmt.Sprintf("Deleted Successfully %v", taskId)},
		nil
}
