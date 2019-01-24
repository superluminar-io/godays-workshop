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
