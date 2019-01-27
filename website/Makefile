HUGO_VERSION := 0.53
ifeq ($(shell uname),Darwin)
	OS := macOS
else
	OS := Linux
endif

.PHONY: clean
clean:
	rm -fr public tmp

.PHONY: clobber
clobber: clean
	rm -fr hugo

.PHONY: deploy
deploy: public guard-WEBSITE_BUCKET guard-CLOUDFRONT_DISTRIBUTION_ID ## Deploys the website to the S3 bucket
	aws s3 sync public/ s3://$(WEBSITE_BUCKET)/ --delete
	aws configure set preview.cloudfront true
	aws cloudfront create-invalidation --distribution-id=$(CLOUDFRONT_DISTRIBUTION_ID) --paths /

.PHONY: guard-%
guard-%:
	@ if [ "${${*}}" = "" ]; then \
	    echo "Environment variable $* not set"; \
	    exit 1; \
	fi

hugo:
	mkdir -p tmp
	curl -LsS https://github.com/gohugoio/hugo/releases/download/v$(HUGO_VERSION)/hugo_$(HUGO_VERSION)_$(OS)-64bit.tar.gz | tar xzf - hugo

develop: hugo ## Start a development server
	./hugo serve -D

.PHONY: build
build: public

public: hugo ## Build the website
	./hugo -b https://go-serverless-workshop.com/

install: hugo ## Install dependencies (hugo)

.PHONY: cloudformation
cloudformation: guard-GITHUB_OAUTH_TOKEN guard-SLACK_HOOK_URL
	aws cloudformation deploy \
		--stack-name go-serverless-workshop-website \
		--region us-east-1 \
		--template-file cloudformation.yaml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides GithubOauthToken=$(GITHUB_OAUTH_TOKEN) SlackHookUrl=$(SLACK_HOOK_URL)
