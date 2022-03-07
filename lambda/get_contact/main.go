package main

import (
	"contacts-crud/cmd/contacts/handler"
	"contacts-crud/cmd/contacts/repository"
	"contacts-crud/cmd/contacts/service"
	"contacts-crud/internal"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	internal.InitDynamoDB()

	//Create repo
	contactRepo := repository.NewContactRepository(internal.DynamoDBSession)
	//Create service
	contactService := service.NewContactService(&contactRepo)
	//Create Handler
	contactHandler := handler.NewContactHandler(&contactService)

	lambda.Start(contactHandler.GetContact)

}
