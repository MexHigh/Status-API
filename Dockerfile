# Build Go backend
FROM golang:1.16 AS go-builder
WORKDIR /go/src/app
COPY server/ .
RUN go get -d -v ./...
RUN CGO_ENABLED=1 GOOS=linux go install -a -ldflags '-linkmode external -extldflags "-static"' .


# Transpile React frontend to static files
FROM node:15 AS react-builder
WORKDIR /tmp
COPY frontend/ .
RUN npm install
RUN npm run build


# Collect builds in scratch image
FROM scratch
LABEL maintainer="Leon Schmidt"

# Copy CA-Certs and Timezone info
COPY --from=go-builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=go-builder /usr/share/zoneinfo /usr/share/zoneinfo
# Copy compiled go binary
COPY --from=go-builder /go/bin/status-api /status-api
# Copy static frontend files
COPY --from=react-builder /tmp/build /frontend/build
# Copy example config
COPY server/config.example.json /config.json

#VOLUME /config.json
#VOLUME /db.sqlite
EXPOSE 3002

ENTRYPOINT ["/status-api"]
CMD ["--config", "/config.json"]