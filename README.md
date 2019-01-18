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

Implement a URL shortener using [`DynamoDB`](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Introduction.html) for storage with two functions.

- Implement one function to create a shortened URL via HTTP POST
- Implement one function to retrieve the full URL via HTTP GET issuing a `302 Found` redirect

Below you find an example interaction with your service:

```
$ curl -v -XPOST -d url=https://godays.io https://$ENDPOINT/create-url

> POST /create-url HTTP/1.1
< HTTP/1.1 Created 201
Created short url: https://$ENDPOINT/${short-url}

$ curl -v http://$ENDPOINT/${short-url}

> GET /${short-url} HTTP/1.1
< HTTP/1.1 302 Found
< Location: https://godays.io
```

Here is an integration test you can run against your service to see if it works.
- Download [integration_test.go](https://raw.githubusercontent.com/superluminar-io/godays-workshop/master/integration_test.go)
- Copy to root of your project
- Run it with: `go test -integrationTest -endpoint=$(sls info -v | awk '/ServiceEndpoint/ { print $2 }')`

### Hints

- Use the aws-sdk-go to talk to DynamoDB
- Create a DynamoDB table using Cloudformation
- Give your Lambda functions permissions to access the DynamoDB table with IAM
- Inject the DynamoDB table via environment variables
- Using path parameters with API Gateway and Lambda

#### Use the aws-sdk-go with DynamoDB

Use the [aws-sdk-go](https://aws.amazon.com/sdk-for-go/) to talk to DynamoDB. Find examples [here](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/using-dynamodb-with-go-sdk.html).
The code below sets up the SDK and issues a `PutItem` request.

```go
import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)
// Create a new AWS session and fail immediately on error
sess := session.Must(session.NewSession())
// Create the DynamoDB client
dynamodbclient := dynamodb.New(sess)
_, err = dynamodbclient.PutItem(&dynamodb.PutItemInput{
	TableName: aws.String(os.Getenv("DYNAMO_DB_TABLE")),
	Item: map[string]*dynamodb.AttributeValue{
		"short_url": &dynamodb.AttributeValue{S: aws.String(s)},
		"url":       &dynamodb.AttributeValue{S: aws.String(u)},
	}})
if err != nil {
	log.WithField("error", err).Error("Couldn't save URL")
}
```

#### Creating a DynamoDB table with Cloudformation

`serverless.yml` supports arbitrary Cloudformation [resources](https://serverless.com/framework/docs/providers/aws/guide/resources/) under the `resources` key.

```
resources:
  Resources:
    DynamoDBTable:
      Type: AWS::DynamoDB::Table
      Properties:
        TableName: url-shortener
        KeySchema:
          - AttributeName: "id"
            KeyType: "HASH"
        ProvisionedThroughput:
          ReadCapacityUnits: "1"
          WriteCapacityUnits: "1"
        AttributeDefinitions:
          - AttributeName: "id"
            AttributeType: "S"
```

#### Give your Lambda functions permissions to access DynamoDB
Every AWS Lambda function needs permission to interact with other AWS infrastructure resources.
Permissions are set via an AWS IAM Role which is automatically created and is shared by all of your functions.
You can set additional [permissions](https://serverless.com/framework/docs/providers/aws/guide/iam/) via the `serverless.yml` file:

```
provider:
  iamRoleStatements:
    -  Effect: "Allow"
       Action:
         - "dynamodb:PutItem"
         - "dynamodb:GetItem"
       Resource:
         Fn::GetAtt:
           - DynamoDBTable
           - Arn
```

#### Inject the DynamoDB table via environment variables

Use the twelve factor app pattern to [store config in the environment](https://12factor.net/config) to make the Lambda functions aware of your table.
Make these settings in `serverless.yml`:

```
provider:
  environment: # Service wide environment variables
    DYNAMO_DB_TABLE:
      Ref: DynamoDBTable # References resource by name

resources:
  Resources:
    DynamoDBTable:
      Type: AWS::DynamoDB::Table
```


### Configure URL routing and request parameters

For `GET /${short-url}` to work you need to configure path parameters for your function.
In your `serverless.yml` add the following:

```
functions:
  get-url:
    handler: bin/get-url
    events:
      - http:
          path: /{short_url}
          method: get
          request:
            parameters:
              paths:
                short_url: true
```

Details can be found [here](https://serverless.com/framework/docs/providers/aws/events/apigateway#request-parameters).

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
 
When pasting a URL into Twitter, Slack or similar we immediately see a preview.
Let's assume we want to save preview information for every URL we shorten.
For each newly created short URL fetch a preview of the website and store it.

- Create a new function
- Create a new table
- Hook up the function to a DynamoDB stream
- Fetch preview using e.g. https://github.com/badoux/goscraper
- Save it

### Enable DynamoDB streams

DynamoDB supports listening to events via [streams](https://docs.amazonaws.cn/en_us/amazondynamodb/latest/developerguide/Streams.Lambda.html).
To make use of this we must enable streaming on our table. Various `StreamViewTypes` exist. We only care about new data being written here:

```
resources:
  Resources:
    DynamoDBTable:
      Type: AWS::DynamoDB::Table
      Properties:
        StreamSpecification: # Add this to your existing table
          StreamViewType: NEW_IMAGE
```

### Change the event type of your function

To listen to events from DynamoDB you must change the [event type](https://serverless.com/framework/docs/providers/aws/events/streams/).
This will subscribe the Lambda function to the DynamoDB stream.

```
functions:
  unfurl:
    handler: handler.unfurl
    events:
      - stream:
          type: dynamodb
          arn:
            Fn::GetAtt:
            - DynamoDBTable
            - StreamArn
```
