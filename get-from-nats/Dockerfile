FROM golang:1.17.8
WORKDIR /nats_getter
COPY . .
RUN go build -o getter .
CMD [ "/nats_getter/getter" ]