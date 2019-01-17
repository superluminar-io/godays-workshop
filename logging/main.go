package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler() (events.APIGatewayProxyResponse, error) {
	fmt.Print("Hi, this ends up in the logs.")
	log.Print("This will also end up in the logs.")
	resp := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "OK",
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}
