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

 - Bootstrap serverless
 - `serverless init`

Deploy serverless project

## Lab 2 - Making changes

- Modify the `hello` function to take a query parameter and print a friendly greeting.

```
curl https://7gxzpfmtk5.execute-api.eu-central-1.amazonaws.com/dev/hello?name=Jan
Hello Jan
```
Hint:
- Have a look at `events.APIGatewayProxyResponse` signature 
- Use `fmt.Println` or `log.Info` for logging
- Use `serverless logs` for debugging

## Lab 3

Implement a URL-shortener (hint: use DynamoDB for storage) with two functions

## Lab 4

Write a unit test using github.com/stretchr/testify/mock to mock out AWS dependencies
Test what happens when your GET /{short-url} handler can't find the short url in your storage

## Lab 5
 
Fetch details about the URL using DynamoDB streams
``Use another Lambda function store a short description and a preview image url in another DynamoDB table