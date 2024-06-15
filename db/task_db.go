package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"serverless-todo-golang/utils/constants"
	"serverless-todo-golang/utils/logger"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Task struct {
	ID        string    `json:"id,omitempty" dynamodbav:"id"`
	Name      string    `json:"name,omitempty" dynamodbav:"name"`
	UserId    string    `json:"user_id,omitempty" dynamodbav:"user_id"`
	Completed bool      `json:"completed" dynamodbav:"completed"`
	CreatedAt time.Time `json:"created_at,omitempty" dynamodbav:"created_at"`
	IsEditing bool      `json:"is_editing" dynamodbav:"is_editing"`
}

var todoClient *dynamodb.Client
var todoTableName string

func init() {
	todoTableName = os.Getenv("TASK_TABLE_NAME")

	if todoTableName == "" {
		log.Fatal("missing environment variable TASK_TABLE_NAME")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	todoClient = dynamodb.NewFromConfig(cfg)
}

func InsertTask(task Task) error {
	logger.GetLog().Info("INFO : ", "Inserting Task Called..")

	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in marshalling the task.")
		return err
	}

	_, err = todoClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(todoTableName),
		Item:      item,
	})

	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in inserting the task: %v", err.Error()))
		return err
	}

	return nil
}

func UpdateTaskById(task Task, taskId string) error {
	logger.GetLog().Info("INFO : ", "Update Task Called..")

	update := expression.Set(expression.Name(constants.TASK_NAME), expression.Value(task.Name))
	updateExpression, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in building the update expression.")
		return err
	}

	_, err = todoClient.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(todoTableName),
		Key: map[string]types.AttributeValue{
			constants.TASK_ID: &types.AttributeValueMemberS{Value: taskId},
		},
		UpdateExpression:          updateExpression.Update(),
		ExpressionAttributeNames:  updateExpression.Names(),
		ExpressionAttributeValues: updateExpression.Values(),
	})
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in updating the task: %v", err.Error()))
		return fmt.Errorf("failed to update task name: %v", err)
	}

	return nil
}

func UpdateTaskToCompletedById(task Task, taskId string) error {
	logger.GetLog().Info("INFO : ", "Update Task to Completed Called..")

	update := expression.Set(expression.Name(constants.COMPLETED), expression.Value(task.Completed))
	updateExpression, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in building the update expression.")
		return err
	}

	_, err = todoClient.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: aws.String(todoTableName),
		Key: map[string]types.AttributeValue{
			constants.TASK_ID: &types.AttributeValueMemberS{Value: taskId},
		},
		UpdateExpression:          updateExpression.Update(),
		ExpressionAttributeNames:  updateExpression.Names(),
		ExpressionAttributeValues: updateExpression.Values(),
	})
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in updating the task: %v", err.Error()))
		return fmt.Errorf("failed to update task to completed: %v", err)
	}

	return nil
}

func GetAllTasksWIthUserId(userId string) ([]Task, error) {
	logger.GetLog().Info("INFO : ", "Get All Tasks Called..")

	var tasks []Task
	filter := expression.Name("user_id").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in building the expression.")
		return nil, fmt.Errorf("failed to build expression: %v", err)
	}

	result, err := todoClient.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:                 aws.String(todoTableName),
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in querying the tasks: %v", err.Error()))
		return nil, fmt.Errorf("failed to query tasks: %v", err)
	}

	for _, item := range result.Items {
		var task Task
		if err := attributevalue.UnmarshalMap(item, &task); err != nil {
			logger.GetLog().Error("ERROR : ", "Error in unmarshalling the task.")
			return nil, fmt.Errorf("failed to unmarshal task: %v", err)
		}
		tasks = append(tasks, task)
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks, nil

}

func DeleteTaskById(taskId string) error {
	logger.GetLog().Info("INFO : ", "Delete Task Called..")

	_, err := todoClient.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: aws.String(todoTableName),
		Key: map[string]types.AttributeValue{
			constants.TASK_ID: &types.AttributeValueMemberS{Value: taskId},
		},
	})

	if err != nil {
		if strings.Contains(err.Error(), "ConditionalCheckFailedException") {
			logger.GetLog().Error("ERROR : ", fmt.Sprintf("task with ID %s not found", taskId))
			return fmt.Errorf("task with ID %s not found", taskId)
		}
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in deleting the task: %v", err.Error()))
		return fmt.Errorf("failed to delete task: %v", err)
	}

	return nil
}
