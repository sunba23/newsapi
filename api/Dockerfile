FROM golang:1.24.0-bookworm AS build

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN go build -o /usr/local/bin ./...

FROM ubuntu:22.04

WORKDIR /app

COPY --from=build /usr/local/bin .

CMD ["/app/news"]
