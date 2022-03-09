package internal

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoDB struct {
	Db dynamodbiface.DynamoDBAPI
}

// Dyna - object from DynamoDB
var Dyna *DynamoDB

// ConfigureDynamoDB - init func for open connection to aws dynamodb

func ConfigureDynamoDB() {
	Dyna = new(DynamoDB)
	awsSession := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))
	svc := dynamodb.New(awsSession)
	Dyna.Db = dynamodbiface.DynamoDBAPI(svc)
}
