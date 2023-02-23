FROM golang:1.19-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/gitlab.com/distributed_lab/acs/unverified-svc
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/unverified-svc /go/src/gitlab.com/distributed_lab/acs/unverified-svc


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/unverified-svc /usr/local/bin/unverified-svc
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["unverified-svc"]
