# FROM golang:1.24-alpine3.19 AS builder
# ARG SRC_PATH=.
# ENV GO111MODULE=on
# COPY ../colibri-sdk-go/ /colibri-sdk-go/
# COPY ${SRC_PATH}/ /build/
# WORKDIR /build
# RUN go mod download
# RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags="-w -s" -o application

FROM alpine:3.19 AS dist
ARG APP_SRC=.
RUN apk add --no-cache tzdata
ENV TZ=America/Sao_Paulo
COPY ${APP_SRC}/migrations/ /application/migrations/
COPY ${APP_SRC}/application /application/
WORKDIR /application
EXPOSE 8080
ENTRYPOINT ["./application"]
