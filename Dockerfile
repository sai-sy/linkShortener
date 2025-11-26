FROM golang:1.25.4 as dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/main

CMD ["go", "run", "."]
