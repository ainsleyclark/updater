format:
	go fmt ./...

lint:
	golangci-lint run ./...

mock:
	rm -rf mocks && mockery --all --keeptree

test:
	go clean -testcache && go test -race $$(go list ./... | grep -v tests)

test-v:
	go clean -testcache && go test -race $$(go list ./... | grep -v tests) -v

test-cover:
	go clean -testcache && go test -v -cover -race $$(go list ./... | grep -v tests)

cover:
	gopherbadger -md="README.md"

all:
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test

travis:
	$(MAKE) format
	$(MAKE) lint
	$(MAKE) test-cover