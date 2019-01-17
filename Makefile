.PHONY: build clean deploy present integration-test

GOPATH := $(shell go env GOPATH)

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

present: $(GOPATH)/bin/present
	$(GOPATH)/bin/present

integration-test:
	go test -integrationTest -endpoint=$(shell sls info -v | awk '/ServiceEndpoint/ { print $$2 }')

$(GOPATH)/bin/present:
	go get golang.org/x/tools/cmd/present
