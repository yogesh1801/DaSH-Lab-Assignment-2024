FROM golang:1.22

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
COPY input.txt ./

EXPOSE 8080

RUN CGO_ENABLED=0 GOOS=linux go build -o /client_bin

CMD ["/client_bin"]