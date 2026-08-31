FROM golang:1.25.0-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go test ./

RUN go build ./

EXPOSE 8080

CMD ["./deploy-demo"]