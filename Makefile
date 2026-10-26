.PHONY: api
%: gen apiHTTP 
api:
	@protoc --proto_path=./api \
	--go_out=paths=source_relative:./api \
	--go-grpc_out=paths=source_relative:./api \
    ./api/review/v1/review.proto
apiHTTP:
	@protoc --proto_path=./api \
	-I third_party/googleapis \
	--grpc-gateway_out=paths=source_relative:./api \
	--grpc-gateway_opt paths=source_relative \
	--go_out=paths=source_relative:./api \
	--go-grpc_out=paths=source_relative:./api \
    ./api/review/v1/review.proto 

gen:
	@go run ./cmd/gen

build:gen apiHTTP
	@mkdir -p bin/ && go build -o ./bin/ ./...