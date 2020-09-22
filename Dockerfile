FROM golang:1.14 AS builder
WORKDIR /go/src/app
COPY src/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go install -v ./...

FROM scratch
LABEL maintainer="Leon Schmidt"
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /go/bin/status-api /status-api
#VOLUME /go/src/app/config.json
EXPOSE 3002
ENTRYPOINT ["/status-api"]
