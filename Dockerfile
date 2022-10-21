FROM golang:1.18-alpine3.15 AS builder
ARG SRC_PATH=.
ENV GO111MODULE=on
COPY /starters/ /starters/
COPY ${SRC_PATH}/ /build/
WORKDIR /build
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs=false -ldflags="-w -s" -o application

FROM alpine:3.15 AS dist
RUN apk add --no-cache tzdata
ENV TZ=America/Sao_Paulo
COPY --from=builder /build/application /application/
WORKDIR /application
EXPOSE 8080
ENTRYPOINT ["./application"]
