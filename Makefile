.PHONY: vet
vet:
	@go vet $(go list ./...)

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test: 
	@go test -cover -coverpkg=./... -coverprofile coverage.out -covermode=atomic ./... --count=1 
	@cat coverage.out | grep -v "testdata/" | grep -v "pb.go" | grep -v "examples" > coverage.tmp
	@mv coverage.tmp coverage.out

.PHONY: testview
testview:
	@go tool cover -html=coverage.out

.PHONY: testcover
testcover:
	@go tool cover -func=coverage.out | tail -n 1