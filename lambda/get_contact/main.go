package main

import (
	"contacts-crud/cmd/contacts/handler"
	"contacts-crud/cmd/contacts/repository"
	"contacts-crud/cmd/contacts/service"
	"contacts-crud/internal"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {

	internal.ConfigureDynamoDB()

	//Create repo
	contactRepo := repository.NewContactRepository(internal.Dyna)
	//Create service
	contactService := service.NewContactService(&contactRepo)
	//Create Handler
	contactHandler := handler.NewContactHandler(&contactService)

	lambda.Start(contactHandler.GetContact)

}
