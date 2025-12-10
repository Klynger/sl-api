FROM golang:alpine

WORKDIR /shopping-list-backend
COPY . .

RUN go build main.go

CMD ["/shopping-list-backend/main"]
EXPOSE 8080
