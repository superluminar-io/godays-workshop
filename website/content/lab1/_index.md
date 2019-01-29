---
title: Lab 1 - Basics
weight: 10
---

In this lab we will create our first Serverless application and explore it.

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

{{% notice tip %}}
Pay attention to the log output. Especially lines prefixed with _CloudFormation - ..._.
To take a look under the hood look into the `.serverless` folder. You'll see two JSON
files. This is the Cloudformation stack for your application.
{{% /notice %}}

## Run the functions

To run the function we can use either the HTTP endpoint (via API Gateway) or use the Serverless tool to invoke the function directly.
Try it out. Figure out the URL to your service by taking a look at the shell output from the command you ran above.

{{% expand "Need help?" %}}
```shell
# Find the endpoint in the output
ENDPOINT=$(sls info -v | grep ServiceEndpoint: | cut -d ' ' -f2)
# Run curl
curl ${ENDPOINT}/hello

# Or invoke the function directly.
serverless invoke -f hello
```
{{% /expand %}}

Do you notice any differences when running the function via `curl` or when you invoke it directly?
