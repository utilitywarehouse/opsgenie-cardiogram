FROM golang:1.14-alpine AS build
WORKDIR /go/src/github.com/utilitywarehouse/opsgenie-cardiogram
COPY . /go/src/github.com/utilitywarehouse/opsgenie-cardiogram
ENV CGO_ENABLED 0
RUN apk --no-cache add git \
  && go get -t ./... \
  && go test -v \
  && go build -o /opsgenie-cardiogram .

FROM alpine:3.12
RUN apk add --no-cache ca-certificates
COPY --from=build /opsgenie-cardiogram /opsgenie-cardiogram

VOLUME /data

ENTRYPOINT ["/opsgenie-cardiogram"]
CMD ["-config", "/data/config.yml"]
