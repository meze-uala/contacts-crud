package internal

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
)

var DynamoDBSession *session.Session

func InitDynamoDB() {

	log.Println("DynamoDB init >>>")

	sess := session.Must(session.NewSessionWithOptions(
		session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

	DynamoDBSession = sess
}
