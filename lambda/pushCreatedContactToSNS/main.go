package main

import (
	"contacts-crud/cmd/pushContactToSNS/handler"
	"github.com/aws/aws-lambda-go/lambda"
	"log"
)

func main() {
	//TODO llevar esto a un handler, para crear el handler tmb no pegarle directamente la func
	log.Println("Entre al main del lambda de publish to sns")
	lambda.Start(handler.PushCreatedContactToSNS)
	log.Println("Despues de invocar al lambda del publish del sns")
}
