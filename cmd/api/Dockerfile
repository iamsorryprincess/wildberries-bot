FROM golang:1.24.3-alpine AS build

ENV GOPATH="" CGO_ENABLED=0
ARG DIR

COPY go.mod go.mod
COPY go.sum go.sum
COPY cmd/$DIR cmd/$DIR
COPY internal internal

RUN go mod download
RUN env GOOS=linux GOARCH=arm64 go build -o /bin/app cmd/$DIR/main.go

FROM alpine:3.20.1

COPY --from=build /bin/app /app/
WORKDIR /app/

ARG PORT
EXPOSE $PORT

CMD ["./app"]