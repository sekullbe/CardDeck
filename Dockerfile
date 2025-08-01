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
COPY pkg ./pkg
COPY templates ./templates
COPY decks ./decks
COPY css ./css

RUN go build -o /carddeck
RUN CGO_ENABLED=0 GOOS=linux go build -o /carddeck .

##
## Deploy
##
FROM scratch
WORKDIR /
COPY --from=build /carddeck /carddeck
EXPOSE 8888
ENTRYPOINT ["/carddeck"]
