---
title: Lab 2 - Making changes
weight: 15
---

Modify the `hello` function to take a query parameter and print a friendly greeting.

```
ENDPOINT=$(sls info -v | grep ServiceEndpoint: | cut -d ' ' -f2)
curl ${ENDPOINT}/hello?name=World
Hello World
```

Hint:

- Have a look at [events.APIGatewayProxyRequest](https://github.com/aws/aws-lambda-go/blob/master/events/apigw.go#L6) signature
- Use `fmt.Println` or `log.Info` for logging
- Use `serverless logs` for debugging
- [Serverless AWS Docs](https://serverless.com/framework/docs/providers/aws/)

