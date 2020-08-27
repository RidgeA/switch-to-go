FROM golang:1.15-alpine as builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

RUN go get -d -v

RUN go build -o /go/bin/switch2go

# todo: scratch
FROM alpine
COPY --from=builder /go/bin/switch2go /switch2go
COPY ./data /data
COPY ./views /views
ENTRYPOINT ["/switch2go"]
