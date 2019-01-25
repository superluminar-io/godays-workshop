---
title: Prerequisites
weight: 5
---
## Install software

If you have a Mac run:
```
brew install go
brew install node
brew install awscli
```

## Configure AWS access

```
$ aws configure
AWS Access Key ID [None]: <your key>
AWS Secret Access Key [None]: <your secret>
Default region name [None]: eu-central-1
Default output format [None]:
```

## Make sure AWS works

```
$ aws sts get-caller-identity
{
    "UserId": "AROAJIFDNOS32O5CUCCXO:1547722234274198000",
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:something/something/something"
}
```

## Install serverless via npm

First we install the [serverless](https://serverless.com/framework/docs/getting-started/) command line tool.
This tool allows us to build and deploy [Serverless](https://en.wikipedia.org/wiki/Serverless_computing) functions.

```
npm install -g serverless
```
