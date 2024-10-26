# syntax=docker/dockerfile:1

# Build the application from source
FROM docker.io/library/golang:1.23 AS build-stage

WORKDIR /app

COPY apps ./apps
COPY services ./services
COPY infra ./infra
COPY domain ./domain

COPY go.mod go.sum ./
RUN go mod download

COPY deployments/pizzeria/*.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /service

FROM build-stage AS run-test-stage
RUN go test -v ./...

FROM alpine AS build-release-stage

WORKDIR /

COPY --from=build-stage /service /service

EXPOSE 1234

# USER nonroot:nonroot

ENTRYPOINT ["/service"]