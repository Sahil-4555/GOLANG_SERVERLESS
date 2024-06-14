package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"serverless-todo-golang/utils/constants"
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
		fmt.Println("CheckIfUserWithThisEmailExists Error: ", err)
		return false, err
	}
	fmt.Println("CheckIfUserWithThisEmailExists: ------>", result.Items)
	return len(result.Items) > 0, nil
}

func InsertUser(user User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		fmt.Println("InsertUser: ", err)
		return err
	}

	fmt.Println("InsertUser: ------>", item)

	_, err = userClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(userTableName),
		Item:      item,
	})

	if err != nil {
		fmt.Println("InsertUser Error: ", err)
		return err
	}

	return nil
}

func FindTheUserWithEmail(email string) (User, error) {

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
		fmt.Println("FindTheUserWithEmail Error: ", err)
		return user, err
	}

	if len(result.Items) == 0 {
		return user, errors.New("user not found")
	}

	item := result.Items[0]

	passwordAV, ok := item[constants.PASSWORD]
	if !ok {
		return user, errors.New("password attribute not found")
	}
	user.Password = passwordAV.(*types.AttributeValueMemberS).Value

	nameAV, ok := item[constants.USER_NAME]
	if !ok {
		return user, errors.New("name attribute not found")
	}
	user.Name = nameAV.(*types.AttributeValueMemberS).Value

	emailAV, ok := item[constants.EMAIL]
	if !ok {
		return user, errors.New("email attribute not found")
	}
	user.Email = emailAV.(*types.AttributeValueMemberS).Value

	idAV, ok := item[constants.USER_ID]
	if !ok {
		return user, errors.New("ID attribute not found")
	}
	user.ID = idAV.(*types.AttributeValueMemberS).Value

	fmt.Println("FindTheUserWithEmail: ------>", user)
	return user, nil
}
