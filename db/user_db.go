package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"serverless-todo-golang/utils/constants"
	"serverless-todo-golang/utils/logger"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	ID        string    `json:"id,omitempty" dynamodbav:"id"`
	Name      string    `json:"name,omitempty" dynamodbav:"name"`
	Email     string    `json:"email,omitempty" dynamodbav:"email"`
	Password  string    `json:"password,omitempty" dynamodbav:"password"`
	CreatedAt time.Time `json:"created_at,omitempty" dynamodbav:"created_at"`
}

var userClient *dynamodb.Client
var userTableName string

func init() {
	userTableName = os.Getenv("USER_TABLE_NAME")
	if userTableName == "" {
		log.Fatal("missing environment variable USER_TABLE_NAME")
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	userClient = dynamodb.NewFromConfig(cfg)
}

func CheckIfUserWithThisEmailExists(email string) (bool, error) {
	logger.GetLog().Info("INFO : ", "CheckIfUserWithThisEmailExists Called..")

	input := &dynamodb.ScanInput{
		TableName:        aws.String(userTableName),
		FilterExpression: aws.String("#email = :email"),
		ExpressionAttributeNames: map[string]string{
			"#email": "email",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := userClient.Scan(context.Background(), input)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in scanning the user table.")
		return false, err
	}

	return len(result.Items) > 0, nil
}

func InsertUser(user User) error {
	logger.GetLog().Info("INFO : ", "Inserting User Called..")

	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in marshalling the user.")
		return err
	}

	_, err = userClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(userTableName),
		Item:      item,
	})

	if err != nil {
		logger.GetLog().Error("ERROR : ", fmt.Sprintf("Error in inserting the user: %v", err.Error()))
		return err
	}

	return nil
}

func FindTheUserWithEmail(email string) (User, error) {
	logger.GetLog().Info("INFO : ", "FindTheUserWithEmail Called..")

	var user User
	input := &dynamodb.ScanInput{
		TableName:        aws.String(userTableName),
		FilterExpression: aws.String("email = :email"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":email": &types.AttributeValueMemberS{Value: email},
		},
	}

	result, err := userClient.Scan(context.Background(), input)
	if err != nil {
		logger.GetLog().Error("ERROR : ", "Error in scanning the user table.")
		return user, err
	}

	if len(result.Items) == 0 {
		logger.GetLog().Error("ERROR : ", "User not found.")
		return user, errors.New("user not found")
	}

	item := result.Items[0]

	passwordAV, ok := item[constants.PASSWORD]
	if !ok {
		logger.GetLog().Error("ERROR : ", "password attribute not found")
		return user, errors.New("password attribute not found")
	}
	user.Password = passwordAV.(*types.AttributeValueMemberS).Value

	nameAV, ok := item[constants.USER_NAME]
	if !ok {
		logger.GetLog().Error("ERROR : ", "name attribute not found")
		return user, errors.New("name attribute not found")
	}
	user.Name = nameAV.(*types.AttributeValueMemberS).Value

	emailAV, ok := item[constants.EMAIL]
	if !ok {
		logger.GetLog().Error("ERROR : ", "email attribute not found")
		return user, errors.New("email attribute not found")
	}
	user.Email = emailAV.(*types.AttributeValueMemberS).Value

	idAV, ok := item[constants.USER_ID]
	if !ok {
		logger.GetLog().Error("ERROR : ", "ID attribute not found")
		return user, errors.New("ID attribute not found")
	}
	user.ID = idAV.(*types.AttributeValueMemberS).Value

	return user, nil
}
