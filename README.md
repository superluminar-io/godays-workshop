## Prerequisites

### Install software

If you have a Mac run:
```
brew install go
brew install node
brew install awscli
```

### Configure AWS access

```
$ aws configure
AWS Access Key ID [None]: <your key>
AWS Secret Access Key [None]: <your secret>
Default region name [None]: eu-central-1
Default output format [None]:
```

### Make sure AWS works

```
$ aws sts get-caller-identity
{
    "UserId": "AROAJIFDNOS32O5CUCCXO:1547722234274198000",
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:something/something/something"
}
```

## Lab 1 - Basics

### Bootstrap serverless app

First we install the [serverless](https://serverless.com/framework/docs/getting-started/) command line tool.
This tool allows us to build and deploy [Serverless](https://en.wikipedia.org/wiki/Serverless_computing) functions.

```
npm install -g serverless
```

Now we bootstrap the Serverless "Hello World" example.
Change directories into your `$GOPATH` and create the usual Go boilerplate.

```
mkdir -p $GOPATH/src/github.com/<YOUR_NAME>
cd $GOPATH/src/github.com/<YOUR_NAME>
```

Now generate your template.

```
serverless create -t aws-go-dep -p godays-workshop
```

This will create an example app with the relevant configuration files and directory structure.

### Build and deploy
Since Lambda requires us to provide the compiled binary, we have to build it beforehand.
```
make build
```
This compiles the functions `hello` and `world` and places the binaries in the `bin/` subdirectory.

To deploy the functions and create the AWS components (Lambda function, API Gateway, DNS Entries) we run:

```
serverless deploy
```

### Run the functions

To run the function we can use either the HTTP endpoint (via API Gateway) or use the Serverless tool to invoke the function directly.
```
ENDPOINT=$(sls info -v | grep ServiceEndpoint: | cut -d ' ' -f2)
curl ${ENDPOINT}/hello

serverless invoke -f hello
```

## Lab 2 - Making changes

Modify the `hello` function to take a query parameter and print a friendly greeting.
```
ENDPOINT=$(sls info -v | grep ServiceEndpoint: | cut -d ' ' -f2)
curl ${ENDPOINT}/hello?name=World
Hello World
```

Hint:
- Have a look at `events.APIGatewayProxyRequest` [signature](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go#L6)
- Use `fmt.Println` or `log.Info` for logging
- Use `serverless logs` for debugging
- [Serverless AWS Docs](https://serverless.com/framework/docs/providers/aws/)

## Lab 3

Implement a URL-shortener (hint: use DynamoDB for storage) with two functions

## Lab 4

Write a unit test using `github.com/stretchr/testify/mock` to mock out AWS dependencies.
Test what happens when your GET /{short-url} handler can't find the short url in your storage.

### Example: How to mock AWS SDK

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

## Lab 5
 
Fetch details about the URL using DynamoDB streams
``Use another Lambda function store a short description and a preview image url in another DynamoDB table
