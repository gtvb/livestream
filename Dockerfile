FROM golang:latest as builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ls-server

FROM alpine:latest

COPY --from=builder /app/.env .env
COPY --from=builder /ls-server /ls-server

EXPOSE 3333
CMD ["/ls-server"]