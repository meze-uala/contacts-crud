package internal

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
	"log"
)

type MockSNS struct {
	SnsClient snsiface.SNSAPI
}

func (msns MockSNS) Publish(input sns.PublishInput) (*sns.PublishOutput, error) {
	log.Println("Log from mocked publish")
	if &input == nil {
		return nil, errors.New("no input provided")
	}

	mockedMessage := "Good test"

	result := sns.PublishOutput{
		MessageId:      &mockedMessage,
		SequenceNumber: nil,
	}

	return &result, nil
}

var SNSClient snsiface.SNSAPI

func ConfigureSNS() {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := sns.New(sess)
	SNSClient = snsiface.SNSAPI(svc)
}
