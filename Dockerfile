FROM golang:alpine

WORKDIR /app

COPY . .

RUN go mod download
EXPOSE 3000

RUN go build cmd/main.go

CMD ["./main"]