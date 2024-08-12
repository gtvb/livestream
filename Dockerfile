FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /ls-server

EXPOSE 3333

CMD ["/ls-server"]
