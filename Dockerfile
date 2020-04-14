FROM golang:alpine

WORKDIR /app
COPY . /app
RUN go build -o ./bin/main .

CMD ["/app/bin/main", "server", "--dokku"]
