FROM golang:1.12-alpine as builder
WORKDIR $GOPATH/src/test_task
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /tmp/main .

FROM scratch
COPY response.json ./
COPY --from=builder /tmp/main ./
EXPOSE 9000
CMD ["/main"]