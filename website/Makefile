WEBSITE := www.go-serverless-workshop.com

public:
	hugo -b /

.PHONY: bucket
bucket:
	aws s3 mb s3://$(WEBSITE)
	aws s3 website s3://$(WEBSITE) --index-document index.html --error-document error.html

.PHONY: deploy
deploy: public
	aws s3 sync --acl public-read --sse --delete public s3://$(WEBSITE)

.PHONY: clean
clean:
	rm -fr public

.PHONY: guard-%
guard-%:
	@ if [ "${${*}}" = "" ]; then \
	    echo "Environment variable $* not set"; \
	    exit 1; \
	fi

.PHONY: cloudformation
cloudformation: guard-GITHUB_OAUTH_TOKEN guard-SLACK_HOOK_URL
	aws cloudformation deploy \
		--stack-name go-serverless-workshop-website \
		--region us-east-1 \
		--template-file cloudformation.yaml \
		--capabilities CAPABILITY_IAM \
		--parameter-overrides GithubOauthToken=$(GITHUB_OAUTH_TOKEN) SlackHookUrl=$(SLACK_HOOK_URL)
