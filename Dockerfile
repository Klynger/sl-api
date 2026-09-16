FROM golang:1.26-alpine

WORKDIR /sl-api
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./bin/api ./cmd/api \
  && go build -o ./bin/migrate ./cmd/migrate

CMD ["/sl-api/bin/api"]
EXPOSE 8080
