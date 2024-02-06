FROM golang:1.21.6-alpine3.19

WORKDIR /app

COPY . .

RUN go mod download

ENV ARIONURL_CONFIG="./configs/local.yaml"
ENV ARIONURL_DB="postgresql://admin:root@localhost:5432/arion"

RUN go build -o ./ cmd/main.go


CMD ["./main"]