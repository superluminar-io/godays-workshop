---
title: Lab 4 - Unit testing
weight: 25
---
Write a unit test using `github.com/stretchr/testify/mock` to mock out AWS dependencies.
Test what happens when your GET /{short-url} handler can't find the short url in your storage.

## Example: How to mock AWS SDK

```go
package main

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDynamoDBClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *mockDynamoDBClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	args := m.Called(input)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func TestHandlerRedirectsToTheLongUrl(t *testing.T) {
	dynamoDBClient := new(mockDynamoDBClient)
	hc := HandlerConfig{
		DynamoDBClient: dynamoDBClient,
		DynamoDBTable:  "myTableName",
	}

	expectedDynamoDBInput := &dynamodb.GetItemInput{
		TableName: aws.String(hc.DynamoDBTable),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {
				S: aws.String("abc123"),
			},
		},
	}

	dynamoResponse := &dynamodb.GetItemOutput{
		Item: map[string]*dynamodb.AttributeValue{
			"url": {S: aws.String("https://heise.de")},
		},
	}
	dynamoDBClient.On("GetItem", expectedDynamoDBInput).Return(dynamoResponse, nil)
	req := events.APIGatewayProxyRequest{
		PathParameters: map[string]string{"short_url": "abc123"},
	}

	resp, err := hc.Handler(req)
	assert.NoError(t, err)
	dynamoDBClient.AssertExpectations(t)
	assert.Equal(t, 302, resp.StatusCode)
}
```

