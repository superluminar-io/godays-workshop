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
- Can you find your functions logs?
- Can you find your functions metrics?

{{% expand "Need help?" %}}
You can either explore your function through the console or via the CLI.

For the console:

- Go to the [AWS console](https://console.aws.amazon.com/lambda/home)
- Select your function
- Click the "Monitoring" tab. Various metrics can be found here.
- Click "View logs in CloudWatch" to see your function's log output

For the CLI:

- Run `serverless logs -f hello`
{{% /expand %}}
