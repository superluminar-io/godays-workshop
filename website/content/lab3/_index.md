---
title: Lab 3 - A URL shortener
weight: 20
---
Implement a URL shortener using [DynamoDB](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Introduction.html) for storage with two functions.

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

{{<mermaid>}}
sequenceDiagram
    participant Browser
    participant APIGateway
    participant Lambda
    participant DynamoDB
    Browser->>APIGateway: POST /create-url
    APIGateway->>Lambda: Invoke
    Lambda->>DynamoDB: PutItem
    DynamoDB-->>Lambda: OK
    Lambda-->>APIGateway: {"url": "foo"}
    APIGateway-->>Browser: HTTP 201 Created {"url": "foo"}
{{< /mermaid >}}

## Hints

- Use the aws-sdk-go to talk to DynamoDB
- Create a DynamoDB table using Cloudformation
- Give your Lambda functions permissions to access the DynamoDB table with IAM
- Inject the DynamoDB table via environment variables
- Using path parameters with API Gateway and Lambda
- Generate a short unique ID for the URL
- Run the integration test to see if your service works as expected

### Use the aws-sdk-go with DynamoDB

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

### Creating a DynamoDB table with Cloudformation

`serverless.yml` supports arbitrary Cloudformation [resources](https://serverless.com/framework/docs/providers/aws/guide/resources/) under the `resources` key.

```yaml
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

### Give your Lambda functions permissions to access DynamoDB
Every AWS Lambda function needs permission to interact with other AWS infrastructure resources.
Permissions are set via an AWS IAM Role which is automatically created and is shared by all of your functions.
You can set additional [permissions](https://serverless.com/framework/docs/providers/aws/guide/iam/) via the `serverless.yml` file:

```yaml
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

### Inject the DynamoDB table via environment variables

Use the twelve factor app pattern to [store config in the environment](https://12factor.net/config) to make the Lambda functions aware of your table.
Make these settings in `serverless.yml`:

```yaml
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

```yaml
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

### Generate a short unique ID for the URL

To generate a short/unique ID you can use this function:

```go
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
```

### Run integration test

Here is an integration test you can run against your service to see if it works.

- Download [integration_test.go](https://raw.githubusercontent.com/superluminar-io/godays-workshop/master/integration_test.go)
- Copy to root of your project
- Run it with: `go test -integrationTest -endpoint=$(sls info -v | awk '/ServiceEndpoint/ { print $2 }')`


