# Two-stage build.

#
# Build app binary.
FROM golang:1.9 AS app_builder
ENV APP_PATH "/go/src/github.com/AAbouZaid/consul-ssh-conf-generator"
COPY . ${APP_PATH}
WORKDIR ${APP_PATH}
RUN go get -v -t .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/consul2ssh

#
# Run app binray.
FROM scratch
COPY --from=app_builder /app/consul2ssh .
EXPOSE 8001
ENTRYPOINT ["./consul2ssh", "listen"]
