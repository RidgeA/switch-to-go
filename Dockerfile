FROM golang:1.15-alpine as builder

RUN apk update && apk add --no-cache git && apk add build-base

WORKDIR $GOPATH/src/mypackage/myapp/
COPY . .

RUN go get -d -v

RUN tree && go generate && go build -o /go/bin/switch2go

# todo: scratch, user
FROM alpine
ENV DATABASE_URL="set me"
COPY --from=builder /go/bin/switch2go /switch2go
COPY ./_coverage /_coverage
ENTRYPOINT ["/switch2go"]
