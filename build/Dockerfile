FROM golang:1.23.4-bullseye as build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

WORKDIR /app/cmd/auth-service-task

RUN ["go", "build", "-o", "main", "."]

EXPOSE 8080

CMD ["/app/cmd/auth-service-task/main"]