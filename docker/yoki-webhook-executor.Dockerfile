FROM golang:1.21 as builder

COPY ./webhooks/ /src/webhooks/
COPY ./common /src/common/
COPY ./yoki-event-worker /src/yoki-event-worker
COPY ./yoki-webhook-executor /src/yoki-webhook-executor
COPY ./go.work /src
COPY ./go.work.sum /src

WORKDIR /src/yoki-webhook-executor

# Build with static binaries attached (to be able to use small images)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o yoki-webhook-executor .

FROM alpine

COPY --from=builder /src/yoki-webhook-executor/yoki-webhook-executor /app/
# https://stackoverflow.com/questions/52969195/docker-container-running-golang-http-client-getting-error-certificate-signed-by
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

## Add the wait script to the image
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/2.9.0/wait /wait
RUN chmod +x /wait

WORKDIR /app
CMD /wait && ./yoki-webhook-executor

# CMD ["sleep", "infinity"]
