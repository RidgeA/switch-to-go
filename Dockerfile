FROM golang:1.15-alpine as builder

RUN apk update && apk add --no-cache git && apk add build-base

WORKDIR /opt/switch2go
COPY . .

RUN go get -d -v

RUN go generate && go build -o /go/bin/switch2go

# todo: scratch, user
FROM alpine
ENV DATABASE_URL="set me"
WORKDIR /opt
COPY --from=builder /go/bin/switch2go switch2go
COPY --from=builder /opt/switch2go/_coverage _coverage
ENTRYPOINT ["/opt/switch2go"]
