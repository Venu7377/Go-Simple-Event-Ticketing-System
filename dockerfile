FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN chmod +x wait-for-it.sh

RUN go build -o myapp ./

EXPOSE 8080

CMD ["./wait-for-it.sh", "mysql:3306", "--", "./myapp"]
