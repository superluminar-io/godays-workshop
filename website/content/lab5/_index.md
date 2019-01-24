---
title: Lab 5 - DynamoDB streams
weight: 30
---
When pasting a URL into Twitter, Slack or similar we immediately see a preview.
Let's assume we want to save preview information for every URL we shorten.
For each newly created short URL fetch a preview of the website and store it.

- Create a new function
- Create a new table
- Hook up the function to a DynamoDB stream
- Fetch preview using e.g. https://github.com/badoux/goscraper
- Save it

## Enable DynamoDB streams

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

## Change the event type of your function

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
