package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2alpha/v2"
	"github.com/aws/aws-cdk-go/awscdkapigatewayv2integrationsalpha/v2"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"

	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const UserTableDynamoDBAttributeName = "id"
const TaskTableDynanoDBAttributeName = "id"
const envUserTableName = "USER_TABLE_NAME"
const envTaskTableName = "TASK_TABLE_NAME"

const SignInFunctionDirectory = "../sign-in"
const SignUpFunctionDirectory = "../sign-up"
const CreateTaskFunctionDirectory = "../create-task"
const UpdateTaskFunctionDirectory = "../update-task"
const DeleteTaskFunctionDirectory = "../delete-task"
const GetTasksFunctionDirectory = "../get-tasks"
const UpdateTaskCompletedFunctionDirectory = "../complete-task"

type ServerlessTODOStackProps struct {
	awscdk.StackProps
}

func NewServerlessTODOStack(scope constructs.Construct, id string, props *ServerlessTODOStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// USER DYNAMODB TABLE
	userDynamoDBTable := awsdynamodb.NewTable(stack, jsii.String("user-serverless-dynamodb-table"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String(UserTableDynamoDBAttributeName),
				Type: awsdynamodb.AttributeType_STRING,
			},
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		})

	// TASK DYNAMODB TABLE
	taskDynamoDBTable := awsdynamodb.NewTable(stack, jsii.String("task-serverless-dynamodb-table"),
		&awsdynamodb.TableProps{
			PartitionKey: &awsdynamodb.Attribute{
				Name: jsii.String(TaskTableDynanoDBAttributeName),
				Type: awsdynamodb.AttributeType_STRING,
			},
			RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
		})

	// API GATEWAY
	todoAPI := awscdkapigatewayv2alpha.NewHttpApi(stack, jsii.String("todo-serverless-http-api"),
		&awscdkapigatewayv2alpha.HttpApiProps{
			CorsPreflight: &awscdkapigatewayv2alpha.CorsPreflightOptions{
				AllowOrigins: &[]*string{jsii.String("*")},
				AllowHeaders: &[]*string{jsii.String("*")},
				AllowMethods: &[]awscdkapigatewayv2alpha.CorsHttpMethod{
					awscdkapigatewayv2alpha.CorsHttpMethod_GET,
					awscdkapigatewayv2alpha.CorsHttpMethod_POST,
					awscdkapigatewayv2alpha.CorsHttpMethod_PUT,
					awscdkapigatewayv2alpha.CorsHttpMethod_DELETE,
					awscdkapigatewayv2alpha.CorsHttpMethod_OPTIONS,
				},
			},
		})

	funcEnvVar := &map[string]*string{envUserTableName: userDynamoDBTable.TableName(), envTaskTableName: taskDynamoDBTable.TableName()}

	// SIGNIN FUNCTION
	signInFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("signin-serverless-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Environment:  funcEnvVar,
			Entry:        jsii.String(SignInFunctionDirectory),
			LogRetention: awslogs.RetentionDays_ONE_WEEK, // Set log retention period
		})

	userDynamoDBTable.GrantReadData(signInFunction)

	signInFunctionIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("signin-function-serverless-integration"), signInFunction, nil)

	todoAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/signIn"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: signInFunctionIntg,
	})

	// SIGNUP FUNCTION
	signUpFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("signup-serverless-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Environment:  funcEnvVar,
			Entry:        jsii.String(SignUpFunctionDirectory),
			LogRetention: awslogs.RetentionDays_ONE_WEEK, // Set log retention period
		})

	userDynamoDBTable.GrantReadWriteData(signUpFunction)

	signUpIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("signup-function-serverless-integration"), signUpFunction, nil)

	todoAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/signUp"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: signUpIntg,
	})

	// CREATE TASK
	createTaskFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("create-task-serverless-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Environment:  funcEnvVar,
			Entry:        jsii.String(CreateTaskFunctionDirectory),
			LogRetention: awslogs.RetentionDays_ONE_WEEK, // Set log retention period
		})

	taskDynamoDBTable.GrantWriteData(createTaskFunction)

	createTaskIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("create-task-function-serverless-integration"), createTaskFunction, nil)

	todoAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/create-task"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_POST},
		Integration: createTaskIntg,
	})

	// UPDATE TASK
	updateTaskFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("update-task-serverless-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Environment:  funcEnvVar,
			Entry:        jsii.String(UpdateTaskFunctionDirectory),
			LogRetention: awslogs.RetentionDays_ONE_WEEK, // Set log retention period
		})

	taskDynamoDBTable.GrantReadWriteData(updateTaskFunction)

	updateTaskIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("update-task-function-serverless-integration"), updateTaskFunction, nil)

	todoAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/update-task/{id}"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_PUT},
		Integration: updateTaskIntg,
	})

	// DELETE TASK
	deleteTaskFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("delete-task-serverless-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Environment:  funcEnvVar,
			Entry:        jsii.String(DeleteTaskFunctionDirectory),
			LogRetention: awslogs.RetentionDays_ONE_WEEK, // Set log retention period
		})

	taskDynamoDBTable.GrantReadWriteData(deleteTaskFunction)

	deleteTaskIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("delete-task-function-serverless-integration"), deleteTaskFunction, nil)

	todoAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/delete-task/{id}"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_DELETE},
		Integration: deleteTaskIntg,
	})

	// GET ALL TASKS
	getTasksFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("get-tasks-serverless-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Environment:  funcEnvVar,
			Entry:        jsii.String(GetTasksFunctionDirectory),
			LogRetention: awslogs.RetentionDays_ONE_WEEK, // Set log retention period
		})

	taskDynamoDBTable.GrantReadData(getTasksFunction)

	getTasksIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("get-tasks-function-serveless-integration"), getTasksFunction, nil)

	todoAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/get-tasks"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_GET},
		Integration: getTasksIntg,
	})

	// UPADTE TASK TO COMPLETED
	updateTaskCompletedFunction := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("update-task-completed-serverless-function"),
		&awscdklambdagoalpha.GoFunctionProps{
			Environment:  funcEnvVar,
			Entry:        jsii.String(UpdateTaskCompletedFunctionDirectory),
			LogRetention: awslogs.RetentionDays_ONE_WEEK, // Set log retention period
		})

	taskDynamoDBTable.GrantReadWriteData(updateTaskCompletedFunction)

	updateTaskCompletedIntg := awscdkapigatewayv2integrationsalpha.NewHttpLambdaIntegration(jsii.String("update-task-completed-function-serverless-integration"), updateTaskCompletedFunction, nil)

	todoAPI.AddRoutes(&awscdkapigatewayv2alpha.AddRoutesOptions{
		Path:        jsii.String("/update-task-completed/{id}"),
		Methods:     &[]awscdkapigatewayv2alpha.HttpMethod{awscdkapigatewayv2alpha.HttpMethod_PUT},
		Integration: updateTaskCompletedIntg,
	})

	// OUTPUT API GATEWAY URL
	awscdk.NewCfnOutput(stack, jsii.String("output"), &awscdk.CfnOutputProps{
		Value:       todoAPI.Url(),
		Description: jsii.String("API Gateway endpoint"),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewServerlessTODOStack(app, "ServerlessTODOAppStack", &ServerlessTODOStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	return nil
}
