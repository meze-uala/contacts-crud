package handler

import (
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	_ "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestNewPushContactToSNSHandler(t *testing.T) {
	type args struct {
		service IPushContactToSNS
	}
	tests := []struct {
		name string
		args args
		want PushContactToSNSHandler
	}{
		{
			name: "Test with nil service should pass",
			args: args{service: nil},
			want: NewPushContactToSNSHandler(nil),
		},
		{
			name: "Test with no nil service should pass",
			args: args{service: NewMockIPushContactToSNS(gomock.NewController(t))},
			want: NewPushContactToSNSHandler(NewMockIPushContactToSNS(gomock.NewController(t))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPushContactToSNSHandler(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPushContactToSNSHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPushContactToSNSHandler_PushCreatedContactToSNS(t *testing.T) {

	record := events.DynamoDBEventRecord{
		AWSRegion:      "",
		Change:         events.DynamoDBStreamRecord{},
		EventID:        "",
		EventName:      "",
		EventSource:    "",
		EventVersion:   "",
		EventSourceArn: "",
		UserIdentity:   nil,
	}

	dynamoDbIDAttr := events.NewStringAttribute("1")

	newImage := map[string]events.DynamoDBAttributeValue{
		"id": dynamoDbIDAttr,
	}

	record.Change.NewImage = newImage
	record.EventName = "INSERT"
	records := []events.DynamoDBEventRecord{record}
	event := events.DynamoDBEvent{Records: records}

	mockedMessage := "Mocked message"
	mockedResult := sns.PublishOutput{
		MessageId:      &mockedMessage,
		SequenceNumber: nil,
	}

	snsService := NewMockIPushContactToSNS(gomock.NewController(t))
	snsService.EXPECT().PublishContactIDToSNS(gomock.Any()).Return(&mockedResult, nil)

	snsHandler := NewPushContactToSNSHandler(snsService)

	result, err := snsHandler.PushCreatedContactToSNS(context.TODO(), event)

	assert.NoError(t, err)
	assert.NotNil(t, result)

}

func TestPushContactToSNSHandler_PushCreatedContactToSNS_Error_From_Service(t *testing.T) {

	record := events.DynamoDBEventRecord{
		AWSRegion:      "",
		Change:         events.DynamoDBStreamRecord{},
		EventID:        "",
		EventName:      "",
		EventSource:    "",
		EventVersion:   "",
		EventSourceArn: "",
		UserIdentity:   nil,
	}

	dynamoDbIDAttr := events.NewStringAttribute("1")

	newImage := map[string]events.DynamoDBAttributeValue{
		"id": dynamoDbIDAttr,
	}

	record.Change.NewImage = newImage
	record.EventName = "INSERT"

	records := []events.DynamoDBEventRecord{record}
	event := events.DynamoDBEvent{Records: records}

	snsService := NewMockIPushContactToSNS(gomock.NewController(t))
	snsService.EXPECT().PublishContactIDToSNS(gomock.Any()).Return(nil, errors.New("error from service"))

	snsHandler := NewPushContactToSNSHandler(snsService)

	result, err := snsHandler.PushCreatedContactToSNS(context.TODO(), event)

	assert.Error(t, err)
	assert.Nil(t, result)

}

//TODO esto puede o no funcar dependiendo de que version este en el proyecto (si separo por type, el SNS se bugea)
/*
func TestPushContactToSNSHandler_PushCreatedContactToSNS_Error_Record_Is_Not_An_Insert(t *testing.T) {

	record := events.DynamoDBEventRecord{
		AWSRegion:      "",
		Change:         events.DynamoDBStreamRecord{},
		EventID:        "",
		EventName:      "",
		EventSource:    "",
		EventVersion:   "",
		EventSourceArn: "",
		UserIdentity:   nil,
	}

	dynamoDbIDAttr := events.NewStringAttribute("1")

	newImage := map[string]events.DynamoDBAttributeValue{
		"id": dynamoDbIDAttr,
	}

	record.Change.NewImage = newImage
	record.EventName = "UPDATE"

	records := []events.DynamoDBEventRecord{record}
	event := events.DynamoDBEvent{Records: records}

	snsService := NewMockIPushContactToSNS(gomock.NewController(t))

	snsHandler := NewPushContactToSNSHandler(snsService)

	result, err := snsHandler.PushCreatedContactToSNS(context.TODO(), event)

	assert.Error(t, err)
	assert.Nil(t, result)

}
*/
func TestPushContactToSNSHandler_PushCreatedContactToSNS_No_Records_Provided(t *testing.T) {

	records := []events.DynamoDBEventRecord{}
	event := events.DynamoDBEvent{Records: records}

	snsService := NewMockIPushContactToSNS(gomock.NewController(t))
	snsHandler := NewPushContactToSNSHandler(snsService)

	result, err := snsHandler.PushCreatedContactToSNS(context.TODO(), event)

	assert.Error(t, err)
	assert.Nil(t, result)

}
