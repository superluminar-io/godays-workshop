## Prerequisites

### Install software

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

First we bootstrap the Serverless "Hello World" example.
```
serverless create -t aws-go-dep -p godays-workshop
```
This will create an example app with the relevant configuration files and directory structure.

### Build and deploy
Since lambda requires us to provide the compiled binary, we have to build it beforehand.
```
make build
```
This compiles the functions `hello` and `world` and places the binaries in the `bin/` subdirectory.

To deploy the functions and create the aws components (Lambda function, API Gateway, DNS Entries) we run:

```
serverless deploy
```

### Run the functions

To run the function we can use either the HTTP endpoint (via API Gateway) or the Serverless.
```
ENDPOINT=$(sls info -v | grep ServiceEndpoint: | cut -d ' ' -f2)
curl ${ENDPOINT}/hello

serverless invoke -f hello
```

## Lab 2 - Making changes

Modify the `hello` function to take a query parameter and print a friendly greeting.

```
curl https://7gxzpfmtk5.execute-api.eu-central-1.amazonaws.com/dev/hello?name=Jan
Hello Jan
```

Hint:
- Have a look at `events.APIGatewayProxyResponse` [signature](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go#L6)
- Use `fmt.Println` or `log.Info` for logging
- Use `serverless logs` for debugging
- [Serverless AWS Docs](https://serverless.com/framework/docs/providers/aws/)

## Lab 3

Implement a URL-shortener (hint: use DynamoDB for storage) with two functions

## Lab 4

Write a unit test using github.com/stretchr/testify/mock to mock out AWS dependencies
Test what happens when your GET /{short-url} handler can't find the short url in your storage

## Lab 5
 
Fetch details about the URL using DynamoDB streams
``Use another Lambda function store a short description and a preview image url in another DynamoDB table