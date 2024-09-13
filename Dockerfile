FROM golang:1.22-alpine AS builder
RUN apk add --update make git curl
ARG MODULE_NAME=avito
COPY . /home/${MODULE_NAME}/
WORKDIR /home/${MODULE_NAME}/
RUN go build -o super_tender cmd/main.go
CMD ["./super_tender"]