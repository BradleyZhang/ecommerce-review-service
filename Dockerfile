FROM golang:1.26.9-alpine AS builder
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=${GOPROXY}

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o review-service ./cmd/review-service

FROM alpine:3.22

COPY --from=builder /build/review-service /

ENTRYPOINT ["/review-service"]
