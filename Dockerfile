FROM golang:1.22
WORKDIR /app
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o tender_service ./cmd/main.go
CMD ["./tender_service"]