include .env
## Build the executable for the current environment
build:
	@go build -ldflags '-w -s' -o ${BINARY}-${VERSION}
## Build executable files in linux environment
static:
	@set CGO_ENABLED=0
	@set GOOS=linux
	@set GOARCH=amd64
	@go build -ldflags '-w -s' -o ${BINARY}-${VERSION}
## Run all tests
test:
	go test $$(go list ./... | grep -v /vendor/) -cover

.PHONY: help
## Show help
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make target'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
	helpMessage = match(lastLine, /^## (.*)/); \
	if (helpMessage) { \
	helpCommand = substr($$1, 0, index($$1, ":")-1); \
	helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
	printf " %-$(TARGET_MAX_CHAR_NUM)s %s\n", helpCommand, helpMessage; \
	} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)
