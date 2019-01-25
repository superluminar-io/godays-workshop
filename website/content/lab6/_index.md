---
title: Lab 6 - Tracing with X-Ray
weight: 30
---
To analyse and debug your serverless application [AWS X-Ray](https://aws.amazon.com/xray/) 
comes in very handy. It is a distributed tracing tool, that measures 
how long it takes your application to complete certain tasks. 
E.g. fetching data from a database system.

X-Ray can also visualize the connection between the components of your application
and can therefor help you reason about it while it grows and get more complex.

To enable this functionality within your serverless application, we want to use
the [tracing plugin](https://www.npmjs.com/package/serverless-plugin-tracing).

To instrument your code to send data to X-Ray we use the [aws-sdk-go](https://aws.amazon.com/sdk-for-go/).

To complete this lab you have to:

 - Install and activate the tracing plugin
 - Give IAM permissions to your function to use X-Ray
 - Instrument your code
 - Generate some traces by calling your function
 - Analyse the data within the AWS Console

## IAM Permission to use X-Ray

The tracing data is stored in `segments` and `telemetryRecords` in X-Ray. In order for
the function to be able to store this data, we need to provide 
it with the corresponding IAM permission.  

```yaml
provider:
  name: aws
  tracing: true # enable tracing 
  iamRoleStatements:
    - Effect: "Allow" 
      Action:
        - "xray:PutTraceSegments"
        - "xray:PutTelemetryRecords"
      Resource:
        - "*"
```

## Install and activate the plugin

The serverless.com framework comes with a wide variety of [plugins](https://github.com/serverless/plugins).
Plugins are node modules that need to be installed and activated in order to be available.

```
npm install --save-dev serverless-plugin-tracing
```

```yaml
plugins:
  - serverless-plugin-tracing
```
