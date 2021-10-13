# syntax=docker/dockerfile:1

##
## Build
##
FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY templates ./templates
COPY static ./static
COPY css ./css

RUN go build -o /carddeck
RUN CGO_ENABLED=0 GOOS=linux go build -o /carddeck .

##
## Deploy
##
FROM alpine:latest
WORKDIR /
COPY --from=build /carddeck /carddeck
EXPOSE 8888
ENTRYPOINT ["/carddeck"]
