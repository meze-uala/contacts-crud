package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
	"log"
	"net/http"
)

type Contact struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Status    string `json:"status"`
}

func main() {
	lambda.Start(HandleAddContactRequest)
}

func HandleAddContactRequest(ctx context.Context, evt events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var requestBody Contact

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

	// Initialize a session that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials
	// and region from the shared configuration file ~/.aws/config.
	//sess := session.Must(session.NewSessionWithOptions(session.Options{
	//	SharedConfigState: session.SharedConfigEnable,
	//}))
	//TODO ademas de refactorizar en capas como servicios y repo, una vez que esto funcione, quitar los hardcore :P
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
		Credentials: credentials.NewStaticCredentials("ASIASLBG62CDURE3S4NG",
			"tNHnx4P6Cbln+aCZhDtjW0m6yJRgdQROPpezZXjH",
			"IQoJb3JpZ2luX2VjEN///////////wEaCXVzLWVhc3QtMSJHMEUCIQDxCoDqFJA+fo9+77ph8u+jJP9w5aWuUbPlnWNOqFzA3AIgcnJJ4mkQRJlzlk4ILSzx0HeDRIrAHk8n2k215UtWVCIqnAMISBADGgwxNjExNDI5ODQ4MzkiDKkN7wOk70UHV3unnir5App7KCyvBWGngctXqNW7GBWN0q8DQCLTOv4hYy5cP6DKYD85ALWofj/cdn7CSK/P512bAZ9DTrsoa+RHmGxO7799XV8mn+JvnUJX3m8hiUIYScEb1QhDKmUhObv2K1NzpTlNEzBquyz3zdMrFY1fqNkRF85M+I8+mG8iV+8mb9zk7DwAE5URBkTwlEt4LExQc7qDoP2B8iI8L+ElnPrNOqQtLqLIaTl6aVPI/8CErICDtOfAYgd+GzX8Y3posbmMY+eSkBtX2+HUrYTYMrEBR5rUngepUJqTZVbhVr1dDKqbAL3RRD9ZsCAFV0cyRSqyZ8APPRZRetC4sW2tq1Bi6PrHnwctiY7uMbQYPtg86Da27fnL30c1oQvtUMclWwNtJvHh1igFqKGFTnfSvl6FTBaL9psf367oyyXT5hqdedSNY/vBPwwAaDdhMbAzKzTYAO30ovuZnvDLLcLHbUKg38TzdYb+EOcg75GcxnYFCoR+YjeKJKWsjDnhMIWrg5EGOqYBi2kFOAxwhsvyqUEV1aRQfYHxKibAeHkVAJ/WJ+xKarsBLOElMXmrGp7/GlPTYUpcAeByeSQcNRSOGJjGG5qZRlT2hAXrnj3sKxyJSh9Q0Aeamvm1RyMKXSy+ZK/3n9t03p6XDe/fHZChb/qIM2ev40hYTLMs6/JYrs/abbxT5F6O+TeHOnOrEoWe/8I1BMSRmDmpMfQzTuvCMQaMbEZdiapS+y8XIg=="),
	})

	// Create DynamoDB client
	dynamoDbClient := dynamodb.New(sess)

	//Fin dynamoDB
	uniqueID, _ := uuid.NewUUID()

	requestBody.ID = uniqueID.String()
	requestBody.Status = "CREATED"

	contactItem, err := dynamodbattribute.MarshalMap(requestBody)
	if err != nil {
		log.Fatalf("Got error marshalling new contact item: %s", err)
	}

	//Persistencia

	tableName := "meze-contacts"

	input := &dynamodb.PutItemInput{
		Item:      contactItem,
		TableName: aws.String(tableName),
	}

	_, err = dynamoDbClient.PutItem(input)
	if err != nil {
		log.Fatalf("Got error calling PutItem: %s", err)
		resp.Body = err.Error()
		resp.StatusCode = http.StatusInternalServerError
		return resp, err
	}

	fmt.Println("¡¡¡Successfully added contact!!!")

	resp.Body = evt.Body
	resp.StatusCode = http.StatusCreated
	return resp, nil

}
