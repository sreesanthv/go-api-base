FROM golang:1.15.4-alpine3.12
WORKDIR /go/src/go-api-base
ADD go.* ./
RUN go mod download

ADD . .
RUN mv .config.json.example .config.json
RUN go build
CMD ["./go-api-base", "serve"]