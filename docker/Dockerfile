# would use alpine, but it's not officially supported yet (https://hub.docker.com/_/golang)
FROM golang:1.16

WORKDIR /go/src/app

RUN useradd --no-log-init tp

COPY . .
RUN go build -v

ENV GIN_MODE=release
EXPOSE 8080

USER tp

CMD ["./teleport_interview"]
