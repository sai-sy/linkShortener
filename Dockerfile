FROM golang:1.25.4 as dev
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod .
COPY . .
RUN go mod download
CMD ["air"]
