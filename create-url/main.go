package main

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"net/url"
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	log "github.com/sirupsen/logrus"
)

// Shorten shortens a URL and will return an error if the URL does not validate.
// The implementation is a bit naive but good enough for a show case.
func Shorten(u string) (string, error) {
	if _, err := url.ParseRequestURI(u); err != nil {
		return "", err
	}
	hash := fnv.New64a()
	hash.Write([]byte(u))
	return strconv.FormatUint(hash.Sum64(), 36), nil
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.WithField("body", request.Body).Info("Received request")
	var data map[string]string
	err := json.Unmarshal([]byte(request.Body), &data)
	if err != nil {
		log.WithField("error", err).Error("Error while reading request")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	u, ok := data["url"]
	if !ok {
		return events.APIGatewayProxyResponse{Body: "no url", StatusCode: 400}, nil
	}

	s, err := Shorten(u)
	if err != nil {
		log.WithField("error", err).Error("Malformed URL")
		return events.APIGatewayProxyResponse{Body: "Malformed URL", StatusCode: 400}, nil
	}
	sess := session.Must(session.NewSession())
	dynamodbclient := dynamodb.New(sess)
	_, err = dynamodbclient.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(os.Getenv("DYNAMO_DB_TABLE")),
		Item: map[string]*dynamodb.AttributeValue{
			"short_url": &dynamodb.AttributeValue{S: aws.String(s)},
			"url":       &dynamodb.AttributeValue{S: aws.String(u)},
		}})
	if err != nil {
		log.WithField("error", err).Error("Couldn't save URL")
		return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 500}, nil
	}
	body := fmt.Sprintf("%s/%s", request.Headers["Host"], s)
	return events.APIGatewayProxyResponse{Body: body, StatusCode: 201}, nil
}

func main() {
	lambda.Start(Handler)
}
