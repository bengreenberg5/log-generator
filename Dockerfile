FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN go build -o /log-generator .

FROM alpine:latest
COPY --from=build /log-generator /log-generator

ENTRYPOINT ["/log-generator"]
