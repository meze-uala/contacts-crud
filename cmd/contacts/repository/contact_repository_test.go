package repository

import (
	"contacts-crud/cmd/contacts/models"
	"contacts-crud/internal"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	dynamock "github.com/gusaul/go-dynamock"
	"github.com/stretchr/testify/assert"
	"log"
	"reflect"
	"testing"
)

var mock *dynamock.DynaMock

func TestNewContactRepository(t *testing.T) {
	type args struct {
		dynamoDB *internal.DynamoDB
	}
	tests := []struct {
		name string
		args args
		want ContactRepository
	}{
		{
			name: "Test with nil dynamoDB should pass",
			args: args{dynamoDB: nil},
			want: NewContactRepository(nil),
		},
		{
			name: "Test with no nil dynamoDB should pass",
			args: args{dynamoDB: &internal.DynamoDB{}},
			want: NewContactRepository(&internal.DynamoDB{}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewContactRepository(tt.args.dynamoDB); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewContactRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestContactRepository_AddContact_Error_Putting_Item(t *testing.T) {

	contact := map[string]interface{}{
		"id":         make(chan string),
		"first_name": "Meze",
	}

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	contactItem, _ := dynamodbattribute.MarshalMap(contact)

	result := dynamodb.PutItemOutput{
		Attributes:            nil,
		ConsumedCapacity:      nil,
		ItemCollectionMetrics: nil,
	}

	mock.ExpectPutItem().ToTable("meze-contacts").WithItems(contactItem).WillReturns(result)

	addResult, err := contactRepository.AddContact(models.Contact{})

	assert.Error(t, err)
	assert.Nil(t, addResult)

}

func TestContactRepository_AddContact(t *testing.T) {

	contact := GetValidContact()

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	contactItem, err := dynamodbattribute.MarshalMap(contact)

	if err != nil {
		log.Fatalf("Got error marshalling new contact item: %s", err)
	}

	result := dynamodb.PutItemOutput{
		Attributes:            nil,
		ConsumedCapacity:      nil,
		ItemCollectionMetrics: nil,
	}

	mock.ExpectPutItem().ToTable("meze-contacts").WithItems(contactItem).WillReturns(result)

	addResult, err := contactRepository.AddContact(contact)

	assert.NoError(t, err)
	assert.NotNil(t, addResult)

}

func TestContactRepository_GetContact(t *testing.T) {
	contact := map[string]interface{}{
		"id": "valid-id",
	}

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	contactItem, err := dynamodbattribute.MarshalMap(contact)

	if err != nil {
		log.Fatalf("Got error marshalling contact item: %s", err)
	}

	result := dynamodb.GetItemOutput{
		ConsumedCapacity: nil,
		Item:             contactItem,
	}

	mock.ExpectGetItem().ToTable("meze-contacts").WithKeys(contactItem).WillReturns(result)
	addResult, err := contactRepository.GetContact("valid-id")

	assert.NoError(t, err)
	assert.NotNil(t, addResult)
}

func TestContactRepository_GetContact_Error_Getting_Result(t *testing.T) {
	contact := map[string]interface{}{
		"id":         "valid-id",
		"first_name": "Meze",
	}

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	contactItem, err := dynamodbattribute.MarshalMap(contact)

	if err != nil {
		log.Fatalf("Got error marshalling contact item: %s", err)
	}

	result := dynamodb.GetItemOutput{
		ConsumedCapacity: nil,
		Item:             contactItem,
	}

	mock.ExpectGetItem().ToTable("meze-contacts").WithKeys(contactItem).WillReturns(result)
	addResult, err := contactRepository.GetContact("valid-id")

	assert.Error(t, err)
	assert.Nil(t, addResult)
}

func TestContactRepository_GetContact_No_Results(t *testing.T) {
	contact := map[string]interface{}{
		"id": "valid-id",
	}

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	contactItem, err := dynamodbattribute.MarshalMap(contact)

	if err != nil {
		log.Fatalf("Got error marshalling contact item: %s", err)
	}

	result := dynamodb.GetItemOutput{
		ConsumedCapacity: nil,
		Item:             nil,
	}

	mock.ExpectGetItem().ToTable("meze-contacts").WithKeys(contactItem).WillReturns(result)
	addResult, err := contactRepository.GetContact("valid-id")

	assert.Error(t, err)
	assert.Nil(t, addResult)
}

func TestContactRepository_GetContact_Error_On_Unmarshall(t *testing.T) {
	contact := map[string]interface{}{
		"id": "valid-id",
	}

	mapChannel := map[string]interface{}{
		"mapField": make(chan string),
	}
	contactResponse := map[string]interface{}{
		"id":         make(chan int),
		"first_name": mapChannel,
	}

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	contactItem, err := dynamodbattribute.MarshalMap(contact)
	contactItemResponse, err := dynamodbattribute.MarshalMap(contactResponse)

	if err != nil {
		log.Fatalf("Got error marshalling contact item: %s", err)
	}

	result := dynamodb.GetItemOutput{
		ConsumedCapacity: nil,
		Item:             contactItemResponse,
	}

	mock.ExpectGetItem().ToTable("meze-contacts").WithKeys(contactItem).WillReturns(result)
	addResult, err := contactRepository.GetContact("valid-id")

	assert.Error(t, err)
	assert.Nil(t, addResult)
}

func TestContactRepository_GetAllContacts(t *testing.T) {
	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	_, err := contactRepository.GetAllContacts()

	assert.NoError(t, err)
}

func TestContactRepository_UpdateContactStatus(t *testing.T) {
	contact := GetValidContact()

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	result := dynamodb.UpdateItemOutput{
		Attributes:            nil,
		ConsumedCapacity:      nil,
		ItemCollectionMetrics: nil,
	}

	mock.ExpectUpdateItem().ToTable(TableName).WithKeys(map[string]*dynamodb.AttributeValue{"id": {
		S: aws.String("valid-id"),
	}}).WillReturns(result)
	addResult, err := contactRepository.UpdateContactStatus(contact.ID)

	assert.NoError(t, err)
	assert.NotNil(t, addResult)
}

func TestContactRepository_UpdateContactStatus_Error(t *testing.T) {
	contact := GetValidContact()

	Dyna := new(internal.DynamoDB)
	Dyna.Db, mock = dynamock.New()

	contactRepository := NewContactRepository(Dyna)

	result := dynamodb.UpdateItemOutput{
		Attributes:            nil,
		ConsumedCapacity:      nil,
		ItemCollectionMetrics: nil,
	}

	mock.ExpectUpdateItem().ToTable(TableName).WithKeys(map[string]*dynamodb.AttributeValue{"id": {
		S: aws.String("invalid-id-just-for-error-way"),
	}}).WillReturns(result)
	addResult, err := contactRepository.UpdateContactStatus(contact.ID)

	assert.Error(t, err)
	assert.Nil(t, addResult)
}

func GetValidContact() models.Contact {
	return models.Contact{
		ID:        "valid-id",
		FirstName: "Meze",
		LastName:  "Law",
		Status:    "CREATED",
	}
}
