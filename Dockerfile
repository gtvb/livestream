FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/cespare/reflex@latest

COPY . .

EXPOSE 3333

CMD ["reflex", "-r", "\\.go$", "-s", "--", "go", "run", "main.go"]
