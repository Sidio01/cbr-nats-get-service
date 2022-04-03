FROM golang:1.17.8
WORKDIR /cbr_api
COPY . .
RUN go build -o cbr .
CMD [ "/cbr_api/cbr" ]