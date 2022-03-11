package main

import (
	"contacts-crud/cmd/pushContactToSNS/handler"
	"contacts-crud/cmd/pushContactToSNS/service"
	"contacts-crud/internal"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	internal.ConfigureSNS()

	pushContactToSNSService := service.NewPushContactToSNSService(internal.SNSClient)
	pushContactToSNSHandler := handler.NewPushContactToSNSHandler(&pushContactToSNSService)

	lambda.Start(pushContactToSNSHandler.PushCreatedContactToSNS)
}
