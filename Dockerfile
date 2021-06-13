# Base image
FROM golang:1.15.4-alpine3.12 AS builder
WORKDIR /go/src/go-api-base
ADD go.* ./
RUN go mod download
ADD . .
RUN mv .config.json.example .config.json
RUN go build

# Starting API
FROM alpine:3.12
WORKDIR /go/src/go-api-base
COPY --from=builder /go/src/go-api-base .
CMD ["./go-api-base", "serve"]