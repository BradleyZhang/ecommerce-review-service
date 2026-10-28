.PHONY: api gen build dev
%: build
api: 
	@buf dep update
	@buf lint
	@buf format -w
	@buf build
	@buf generate
gen:
	@go run ./cmd/gen
	@go generate ./...

build:gen api
	@mkdir -p bin/ && go build -o ./bin/review-service ./cmd/review-service

dev:build
	@./bin/review-service