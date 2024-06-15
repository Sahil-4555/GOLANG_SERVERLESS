package main

import (
	"context"
	"encoding/json"
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
	logger.GetLog().Info("INFO : ", "Complete Task Handler Called..")

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

	taskId := request.PathParameters[pathParameterName]

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		logger.GetLog().Error("ERROR : ", "Error in unmarshalling the request body.")
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = db.UpdateTaskToCompletedById(req, taskId)
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in updating task: %v", err.Error()))
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	logger.GetLog().Info("INFO : ", "Task Updated Successfully.")
	return events.APIGatewayV2HTTPResponse{
		StatusCode: http.StatusCreated,
		Body:       fmt.Sprintf("Updated Successfully"),
	}, nil
}
