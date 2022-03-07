package handler

import (
	"contacts-crud/cmd/contacts/models"
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
)

type IContactService interface {
	AddContact(contact models.Contact) (*models.Contact, error)
	GetContact(id string) (*models.Contact, error)
	GetAllContacts() ([]*models.Contact, error)
}

type ContactHandler struct {
	contactService IContactService
}

func NewContactHandler(contactService IContactService) ContactHandler {
	return ContactHandler{contactService: contactService}
}

func (ch *ContactHandler) AddContact(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var requestBody models.Contact

	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,HEAD,OPTIONS,POST",
		},
	}

	err := json.Unmarshal([]byte(evt.Body), &requestBody)

	if err != nil {
		resp.StatusCode = http.StatusBadRequest
		return resp, err
	}

	createdContact, err := ch.contactService.AddContact(requestBody)

	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	log.Println("Contact created successfully with id: ", createdContact.ID)

	createdContactJson, err := json.Marshal(&createdContact)

	if err != nil {
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	resp.Body = string(createdContactJson)
	resp.StatusCode = http.StatusCreated
	return resp, nil

}

//func GetContact(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
//	return nil, nil
//}

//func GetAllContacts(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

//}
