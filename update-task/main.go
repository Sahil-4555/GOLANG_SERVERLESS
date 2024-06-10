package main

import (
	"context"
	"encoding/json"
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

	taskId := request.PathParameters[pathParameterName]

	if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
		fmt.Println("Error in unmarshalling the request body: ", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	err = db.UpdateTaskById(req, taskId)
	if err != nil {
		fmt.Println("Error in updating the task: --->", err)
		return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayV2HTTPResponse{
			StatusCode: http.StatusCreated,
			Body:       fmt.Sprintf("Updated Successfully %v", taskId)},
		nil
}
