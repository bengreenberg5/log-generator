FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod main.go ./
RUN go build -o /log-generator .

FROM alpine:latest
COPY --from=build /log-generator /log-generator

# Default: 5 KiB/s. Override with: docker run log-generator 10
ENTRYPOINT ["/log-generator"]
CMD ["5"]
