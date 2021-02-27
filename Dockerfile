FROM golang:1.12-alpine

RUN apk add --no-cache git

WORKDIR /app/bot

COPY go.mod .
COPY go.sum .


RUN go mod download

COPY . .


RUN go build -o ./build/bot .

EXPOSE 80:8081

CMD ["./build/bot"]