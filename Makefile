NAME = server

run: build
	@./bin/$(NAME)

build: clean
	@go build -o bin/$(NAME) ./cmd/app

clean:
	@rm -rf bin

unit-test:
	@go clean -testcache
	@go test `go list ./... | grep -v ./cmd/app | grep -v ./internals/database | grep -v ./mocks | grep -v ./tests`

coverage:
	@go clean -testcache
	@go test `go list ./... | grep -v ./cmd/app | grep -v ./mocks` -coverprofile=coverage.out
	@go tool cover -html=coverage.out

integration:
	@go clean -testcache
	@go test `go list ./... | grep  ./tests` -v -cover

