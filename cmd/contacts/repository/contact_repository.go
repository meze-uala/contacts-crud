package repository

import (
	"contacts-crud/cmd/contacts/models"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

const TABLE_NAME = "meze-contacts"

type ContactRepository struct {
	dynamoDBSession *session.Session
}

func NewContactRepository(dynamoDBSession *session.Session) ContactRepository {
	return ContactRepository{dynamoDBSession: dynamoDBSession}
}

func (cr *ContactRepository) AddContact(contact models.Contact) (*models.Contact, error) {

	//TODO ver si el cliente no deberia ser reubicado
	dynamoDbClient := dynamodb.New(cr.dynamoDBSession)

	contactItem, err := dynamodbattribute.MarshalMap(contact)
	if err != nil {
		log.Fatalf("Got error marshalling new contact item: %s", err)
		return nil, err
	}

	input := &dynamodb.PutItemInput{
		Item:      contactItem,
		TableName: aws.String(TABLE_NAME),
	}

	log.Println("Ready to insert the new contact into dynamoDB")

	_, err = dynamoDbClient.PutItem(input)

	if err != nil {
		log.Println("Got error calling PutItem: ", err)
		return nil, err
	}

	return &contact, nil

}
func (cr *ContactRepository) GetContact(id string) (*models.Contact, error) {

	var contact models.Contact
	//TODO ver si el cliente no deberia ser reubicado
	dynamoDbClient := dynamodb.New(cr.dynamoDBSession)

	log.Println("Ready to get the contact from dynamoDB")

	result, err := dynamoDbClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &contact)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	return &contact, nil
}
func (cr *ContactRepository) GetAllContacts() ([]*models.Contact, error) {
	return nil, nil
}
