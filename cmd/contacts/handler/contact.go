package handler

//go:generate mockgen -source contact.go -destination mock_contact_handler.go -package handler

import (
	"contacts-crud/cmd/contacts/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"log"
	"net/http"
	"net/url"
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

func (ch *ContactHandler) GetContact(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var requestBody models.Contact

	resp := events.APIGatewayProxyResponse{
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "GET,HEAD,OPTIONS,POST",
		},
	}

	rawIDParam, found := evt.PathParameters["id"]

	if !found {
		log.Println("Missing contact id on URL")
		return resp, errors.New("missing the ID in the url")
	}

	value, err := url.QueryUnescape(rawIDParam)
	if nil != err {
		log.Println("Error al intentar parsear el id de la url: ", err.Error())
		return resp, err
	}
	requestBody.ID = value

	retrievedContact, err := ch.contactService.GetContact(requestBody.ID)

	if err != nil {
		fmt.Println("An error ocurred trying to get the contact. Error: ", err.Error())
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	if retrievedContact == nil {
		fmt.Println("Contact with id " + requestBody.ID + "not found!")
		resp.Body = "Contact not found"
		resp.StatusCode = http.StatusNotFound
		return resp, nil
	}

	retrievedContactJson, err := json.Marshal(&retrievedContact)

	fmt.Println("Contact with id " + requestBody.ID + " retrieved successfully!")

	resp.Body = string(retrievedContactJson)
	resp.StatusCode = http.StatusOK
	return resp, nil
}

//func GetAllContacts(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

//}
