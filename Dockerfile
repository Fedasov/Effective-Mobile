FROM golang:1.25-alpine AS build

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go mod tidy

RUN go build -o app ./cmd/server

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/app .

EXPOSE 8080

CMD ["./app"]