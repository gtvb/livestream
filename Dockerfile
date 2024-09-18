FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@latest

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ls-server

EXPOSE 3333

CMD ["/ls-server"]