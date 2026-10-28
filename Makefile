.PHONY: api gen build dev
%: build
gen:
	@go run ./cmd/gen
	@go generate ./...

build:gen
	@mkdir -p bin/ && go build -o ./bin/review-service ./cmd/review-service

dev:build
	@echo "===== build finish ====="
	@./bin/review-service