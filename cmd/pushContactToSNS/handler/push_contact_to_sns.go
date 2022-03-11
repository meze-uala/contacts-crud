package handler

//go:generate mockgen -source push_contact_to_sns.go -destination mock_push_contact_to_sns_handler.go -package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/sns"
	"log"
)

type IPushContactToSNS interface {
	PublishContactIDToSNS(id string) (*sns.PublishOutput, error)
}

type PushContactToSNSHandler struct {
	pushToSNSService IPushContactToSNS
}

func NewPushContactToSNSHandler(service IPushContactToSNS) PushContactToSNSHandler {
	return PushContactToSNSHandler{pushToSNSService: service}
}

func (pch *PushContactToSNSHandler) PushCreatedContactToSNS(ctx context.Context, e events.DynamoDBEvent) (*sns.PublishOutput, error) {
	for _, record := range e.Records {
		if record.EventName == "INSERT" {
			fmt.Printf("Processing request data for event ID %s, type %s.\n", record.EventID, record.EventName)

			id := record.Change.NewImage["id"].String()

			result, err := pch.pushToSNSService.PublishContactIDToSNS(id)

			if err != nil {
				log.Println("An error occurred when trying to publish to SNS. Error: ", err.Error())
				return nil, err
			}

			log.Println("Result from sns publish: ", result.String())
			return result, nil
		}
		log.Println("record is not an insert")
	}

	return nil, errors.New("no valid records for proccesing found")
}
