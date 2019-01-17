## Prerequisites



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