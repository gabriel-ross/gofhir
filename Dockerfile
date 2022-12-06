# syntax=docker/dockerfile:1

FROM golang:1.19

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN mkdir -p bin
RUN go build -o ./bin/gofhir ./cmd/gofhir.go

EXPOSE ${PORT}
CMD ./bin/gofhir