FROM golang:1.22.4-alpine

WORKDIR /app

COPY . .

RUN go build -o posts cmd/blog/main.go

CMD ["./posts", "--config=config/prod.yaml"]