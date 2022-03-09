package repository

import (
	"contacts-crud/cmd/contacts/models"
	"contacts-crud/internal"
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

const TableName = "meze-contacts"

type ContactRepository struct {
	dynamoDB *internal.DynamoDB
}

func NewContactRepository(dynamoDB *internal.DynamoDB) ContactRepository {
	return ContactRepository{dynamoDB: dynamoDB}
}

func (cr *ContactRepository) AddContact(contact models.Contact) (*models.Contact, error) {

	//Esto no puede dar un error, ya que lo valido antes en el handler.
	contactItem, _ := dynamodbattribute.MarshalMap(contact)

	input := &dynamodb.PutItemInput{
		Item:      contactItem,
		TableName: aws.String(TableName),
	}

	log.Println("Ready to insert the new contact into dynamoDB")

	_, err := cr.dynamoDB.Db.PutItem(input)

	if err != nil {
		log.Println("Got error calling PutItem: ", err)
		return nil, err
	}

	return &contact, nil

}
func (cr *ContactRepository) GetContact(id string) (*models.Contact, error) {

	var contact models.Contact

	log.Println("Ready to get the contact from dynamoDB")

	result, err := cr.dynamoDB.Db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(id),
			},
		},
	})
	if err != nil {
		log.Println("Got error calling GetItem: ", err)
		return nil, err
	}

	if result.Item == nil {
		return nil, errors.New("contact not found")
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &contact)
	if err != nil {
		log.Println("Failed to unmarshal Record: ", err)
		return nil, err
	}

	return &contact, nil
}

func (cr *ContactRepository) GetAllContacts() ([]*models.Contact, error) {
	return nil, nil
}
