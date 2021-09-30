FROM golang:1.16.6-alpine3.14 AS builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR $GOPATH/src/github.com/bms/
ENV GO111MODULE=on
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /go/bin/bms

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/bin/bms /go/bin/bms
EXPOSE 8080
ENTRYPOINT ["/go/bin/bms"]
