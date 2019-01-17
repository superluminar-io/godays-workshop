package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	log "github.com/sirupsen/logrus"
)

type HandlerConfig struct {
	DynamoDBTable  string
	DynamoDBClient dynamodbiface.DynamoDBAPI
}

func (hc *HandlerConfig) Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	s := request.PathParameters["short_url"]
	log.WithField("short_url", s).Info("Got short URL")

	result, err := hc.DynamoDBClient.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(hc.DynamoDBTable),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {S: aws.String(s)},
		},
	})
	if err != nil {
		log.WithField("error", err).Info("Couldn't get data from DynamoDB")
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}
	if result.Item == nil {
		return events.APIGatewayProxyResponse{StatusCode: 404, Body: "nix"}, nil
	}
	loc := result.Item["url"]
	return events.APIGatewayProxyResponse{StatusCode: 302, Headers: map[string]string{"Location": *loc.S}}, nil
}

func main() {
	sess := session.Must(session.NewSession())
	hc := &HandlerConfig{
		DynamoDBTable:  os.Getenv("DYNAMO_DB_TABLE"),
		DynamoDBClient: dynamodb.New(sess),
	}
	lambda.Start(hc.Handler)
}
