FROM golang:1.16.4-buster

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8000
CMD ["url-shortener"]
