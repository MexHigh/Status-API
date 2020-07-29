FROM golang:1.14

WORKDIR /go/src/app
COPY src/ .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3002
VOLUME /go/src/app/config.json

CMD ["app"]