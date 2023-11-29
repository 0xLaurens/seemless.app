FROM golang:1.21-alpine
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ADD cmd ./cmd
ADD internal ./internal

RUN go build -o api ./cmd/web/main.go
EXPOSE 3000

CMD ["./api"]