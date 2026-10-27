.PHONY: api
%: gen api
# api:
# 	@protoc --proto_path=./api \
# 	-I third_party/googleapis \
# 	--grpc-gateway_out=paths=source_relative:./api \
# 	--grpc-gateway_opt paths=source_relative \
# 	--go_out=paths=source_relative:./api \
# 	--go-grpc_out=paths=source_relative:./api \
#     ./api/review/v1/review.proto 

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
	@mkdir -p bin/ && go build -o ./bin/ ./...