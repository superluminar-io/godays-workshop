---
title: Lab 1 - Basics
weight: 10
---
## Bootstrap serverless app

First we bootstrap the Serverless "Hello World" example.
Change directories into your `$GOPATH` and create the usual Go boilerplate.

```shell
mkdir -p $GOPATH/src/github.com/<YOUR_NAME>
cd $GOPATH/src/github.com/<YOUR_NAME>
```

Now generate your template.

```shell
serverless create -t aws-go-dep -p godays-workshop
```

This will create an example app with the relevant configuration files and directory structure.

## Build and deploy
Since Lambda requires us to provide the compiled binary, we have to build it beforehand.
```shell
make build
```
This compiles the functions `hello` and `world` and places the binaries in the `bin/` subdirectory.

To deploy the functions and create the AWS components (Lambda function, API Gateway, DNS Entries) we run:

```shell
make deploy
```

This will run `serverless deploy` to deploy our functions.

*Hint*: Pay attention to the log output. Especially lines prefixed with _CloudFormation - ..._

## Run the functions

To run the function we can use either the HTTP endpoint (via API Gateway) or use the Serverless tool to invoke the function directly.
```shell
ENDPOINT=$(sls info -v | grep ServiceEndpoint: | cut -d ' ' -f2)
curl ${ENDPOINT}/hello

serverless invoke -f hello
```

