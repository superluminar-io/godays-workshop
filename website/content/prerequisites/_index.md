---
title: Prerequisites
weight: 5
---

You will need some software to work through this workshop.

## Install software

If you have a Mac run:
```shell
brew install go
brew install node
brew install awscli
```

## Configure AWS access

```shell
$ aws configure
AWS Access Key ID [None]: <your key>
AWS Secret Access Key [None]: <your secret>
Default region name [None]: eu-central-1
Default output format [None]:
```

## Make sure AWS works

```shell
$ aws sts get-caller-identity
{
    "UserId": "AROAJIFDNOS32O5CUCCXO:1547722234274198000",
    "Account": "123456789012",
    "Arn": "arn:aws:sts::123456789012:something/something/something"
}
```

## Install serverless via npm

Install the [serverless](https://serverless.com/framework/docs/getting-started/) command line tool.
This tool allows us to build and deploy [Serverless](https://en.wikipedia.org/wiki/Serverless_computing) functions.

```shell
npm install -g serverless
```

## Install dep

We will use a serverless template that uses `dep` for dependency  management.

```shell
go get -u github.com/golang/dep/cmd/dep
```
