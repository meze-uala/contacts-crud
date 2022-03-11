package service

//go:generate mockgen -source push_contact_to_sns_service.go -destination mock_push_contact_to_sns_service.go -package service

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"log"
)

var pushContactToSNSTopic = "arn:aws:sns:us-east-1:161142984839:meze-new-contact-topic"

type PushContactToSNSService struct {
	SNSClient snsiface.SNSAPI
}

func NewPushContactToSNSService(snsClient snsiface.SNSAPI) PushContactToSNSService {
	return PushContactToSNSService{SNSClient: snsClient}
}

func (pcs *PushContactToSNSService) PublishContactIDToSNS(id string) (*sns.PublishOutput, error) {

	messageToPush := id

	result, err := pcs.SNSClient.Publish(&sns.PublishInput{
		Message:  &messageToPush,
		TopicArn: &pushContactToSNSTopic,
	})

	if err != nil {
		log.Println("Error trying to publish to the SNS Topic (on contact create): ", err.Error())
		return nil, err
	}

	return result, nil
}
