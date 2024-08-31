FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY .env ./
COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /server_bin

EXPOSE 8080

CMD ["/server_bin"]

