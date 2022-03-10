package handler

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

type IPushContactToSNS interface {
}

func PushCreatedContactToSNS(ctx context.Context, e events.DynamoDBEvent) {
	for _, record := range e.Records {
		fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

		log.Println("Trying to recover the dynamoDB created contact ID")
		// Print new values for attributes name and age
		id := record.Change.NewImage["id"].String()
		//age, _ := record.Change.NewImage["age"].Integer()
		log.Println("After process,pushed ID ISS: ", id)
		//fmt.Printf("Name: %s, age: %d\n", name, age)
		fmt.Println("Pushed ID is: ", id)

		//LLamo al SNS para mandarle el id y enviar el email
		//TODO enviar lo siguiente a un servicio
		// Initialize a session that the SDK will use to load
		// credentials from the shared credentials file. (~/.aws/credentials).
		sess := session.Must(session.NewSessionWithOptions(session.Options{
			SharedConfigState: session.SharedConfigEnable,
		}))

		svc := sns.New(sess)
		//TODO Llevar a una const
		topic := "arn:aws:sns:us-east-1:161142984839:meze-new-contact-topic"
		result, err := svc.Publish(&sns.PublishInput{
			Message:  &id,
			TopicArn: &topic,
		})
		if err != nil {
			fmt.Println(err.Error())
			log.Println("Error al intentar hacer publish al topic: ", err.Error())
		}

		//fmt.Println(*result.MessageId)
		log.Println("Result: ", *result.MessageId)
	}
}
