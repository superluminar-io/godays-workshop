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

- Have a look at [ApiGatewayEvent](https://github.com/aws/aws-lambda-go/blob/master/events/README_ApiGatewayEvent.md)
- Use `fmt.Println` or `log.Info` for logging
- Use `serverless logs` for debugging
- [Serverless docs](https://serverless.com/framework/docs/providers/aws/events/apigateway/#request-parameters)

