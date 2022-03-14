package handler

import (
	"contacts-crud/cmd/contacts/models"
	"context"
	"errors"
	"github.com/aws/aws-lambda-go/events"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"reflect"
	"testing"
)

func TestNewContactHandler(t *testing.T) {
	type args struct {
		contactService IContactService
	}
	tests := []struct {
		name string
		args args
		want ContactHandler
	}{
		{
			name: "Test with nil service should create a handler",
			args: args{contactService: nil},
			want: NewContactHandler(nil),
		},
		{
			name: "Test with no nil service should create a handler",
			args: args{contactService: NewMockIContactService(gomock.NewController(t))},
			want: NewContactHandler(NewMockIContactService(gomock.NewController(t))),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContactHandler(tt.args.contactService); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContactHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactHandler_AddContact(t *testing.T) {

	ctx := context.TODO()
	evt := events.APIGatewayProxyRequest{
		Resource:                        "/contacts",
		Path:                            "/",
		HTTPMethod:                      "POST",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "{ \"first_name\": \"Fulano\",  \"last_name\": \"De Tal\"}",
		IsBase64Encoded:                 false,
	}

	contact := GetValidContact()

	contactService := NewMockIContactService(gomock.NewController(t))
	contactService.EXPECT().AddContact(gomock.Any()).Return(&contact, nil)

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.AddContact(ctx, evt)

	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusCreated)

}

func TestContactHandler_AddContact_Error_On_Request_Body(t *testing.T) {

	ctx := context.TODO()
	evt := events.APIGatewayProxyRequest{
		Resource:                        "/contacts",
		Path:                            "/",
		HTTPMethod:                      "POST",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "{ \"first_name\": \"Fulano\",  \"last_name\": \"De Tal}",
		IsBase64Encoded:                 false,
	}

	contactService := NewMockIContactService(gomock.NewController(t))

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.AddContact(ctx, evt)

	assert.Error(t, err)
	assert.Equal(t, result.StatusCode, http.StatusBadRequest)

}

func TestContactHandler_AddContact_Error_Service_Result(t *testing.T) {

	ctx := context.TODO()
	evt := events.APIGatewayProxyRequest{
		Resource:                        "/contacts",
		Path:                            "/",
		HTTPMethod:                      "POST",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "{ \"first_name\": \"Fulano\",  \"last_name\": \"De Tal\"}",
		IsBase64Encoded:                 false,
	}

	contactService := NewMockIContactService(gomock.NewController(t))
	contactService.EXPECT().AddContact(gomock.Any()).Return(nil, errors.New("error from contact service"))

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.AddContact(ctx, evt)

	assert.Error(t, err)
	assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

}

func TestContactHandler_GetContact(t *testing.T) {

	ctx := context.TODO()
	evt := events.APIGatewayProxyRequest{
		Resource:                        "/contacts",
		Path:                            "/{id}",
		HTTPMethod:                      "GET",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  map[string]string{"id": "valid-id"},
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "",
		IsBase64Encoded:                 false,
	}

	contact := GetValidContact()

	contactService := NewMockIContactService(gomock.NewController(t))
	contactService.EXPECT().GetContact(gomock.Any()).Return(&contact, nil)

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.GetContact(ctx, evt)

	assert.NoError(t, err)
	assert.Equal(t, result.StatusCode, http.StatusOK)

}

func TestContactHandler_GetContact_Missing_ID_URL_Param(t *testing.T) {

	ctx := context.TODO()
	evt := events.APIGatewayProxyRequest{
		Resource:                        "/contacts",
		Path:                            "/{id}",
		HTTPMethod:                      "GET",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  nil,
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "",
		IsBase64Encoded:                 false,
	}

	contactService := NewMockIContactService(gomock.NewController(t))

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.GetContact(ctx, evt)

	assert.Error(t, err)
	assert.Equal(t, result.StatusCode, http.StatusBadRequest)

}

func TestContactHandler_GetContact_Service_Error(t *testing.T) {

	ctx := context.TODO()
	evt := events.APIGatewayProxyRequest{
		Resource:                        "/contacts",
		Path:                            "/{id}",
		HTTPMethod:                      "GET",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  map[string]string{"id": "valid-id"},
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "",
		IsBase64Encoded:                 false,
	}

	contactService := NewMockIContactService(gomock.NewController(t))
	contactService.EXPECT().GetContact(gomock.Any()).Return(nil, errors.New("error from contact service"))

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.GetContact(ctx, evt)

	assert.Error(t, err)
	assert.Equal(t, result.StatusCode, http.StatusInternalServerError)

}

func TestContactHandler_GetContact_Service_Return_Not_Found(t *testing.T) {

	ctx := context.TODO()
	evt := events.APIGatewayProxyRequest{
		Resource:                        "/contacts",
		Path:                            "/{id}",
		HTTPMethod:                      "GET",
		Headers:                         nil,
		MultiValueHeaders:               nil,
		QueryStringParameters:           nil,
		MultiValueQueryStringParameters: nil,
		PathParameters:                  map[string]string{"id": "valid-id"},
		StageVariables:                  nil,
		RequestContext:                  events.APIGatewayProxyRequestContext{},
		Body:                            "",
		IsBase64Encoded:                 false,
	}

	contactService := NewMockIContactService(gomock.NewController(t))
	contactService.EXPECT().GetContact(gomock.Any()).Return(nil, errors.New("contact not found"))

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.GetContact(ctx, evt)

	assert.Error(t, err)
	assert.Equal(t, result.StatusCode, http.StatusNotFound)

}

func TestContactHandler_UpdateContact(t *testing.T) {
	ctx := context.TODO()

	record := events.SNSEventRecord{
		EventVersion:         "",
		EventSubscriptionArn: "",
		EventSource:          "",
		SNS:                  events.SNSEntity{},
	}

	record.SNS.Message = "1"

	evt := events.SNSEvent{Records: []events.SNSEventRecord{record}}

	contact := GetValidContact()

	contactService := NewMockIContactService(gomock.NewController(t))
	contactService.EXPECT().UpdateContactStatus(gomock.Any()).Return(&contact, nil)

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.UpdateContact(ctx, evt)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestContactHandler_UpdateContact_Service_Error(t *testing.T) {
	ctx := context.TODO()

	record := events.SNSEventRecord{
		EventVersion:         "",
		EventSubscriptionArn: "",
		EventSource:          "",
		SNS:                  events.SNSEntity{},
	}

	record.SNS.Message = "1"

	evt := events.SNSEvent{Records: []events.SNSEventRecord{record}}

	contactService := NewMockIContactService(gomock.NewController(t))
	contactService.EXPECT().UpdateContactStatus(gomock.Any()).Return(nil, errors.New("error from service"))

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.UpdateContact(ctx, evt)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestContactHandler_UpdateContact_Error_Empty_Record_List(t *testing.T) {
	ctx := context.TODO()

	evt := events.SNSEvent{Records: []events.SNSEventRecord{}}

	contactService := NewMockIContactService(gomock.NewController(t))

	contactHandler := NewContactHandler(contactService)

	result, err := contactHandler.UpdateContact(ctx, evt)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func GetValidContact() models.Contact {
	return models.Contact{
		ID:        "valid-id",
		FirstName: "Meze",
		LastName:  "Law",
		Status:    "CREATED",
	}
}
