FROM golang:1.17

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["gin-mock-s3"]

EXPOSE 8888
