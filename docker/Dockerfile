FROM node:16 AS builder
WORKDIR /app
COPY frontend/package.json frontend/package-lock.json ./
RUN npm install
COPY frontend ./
RUN npm run build

# would use alpine, but it's not officially supported yet (https://hub.docker.com/_/golang)
FROM golang:1.16

WORKDIR /go/src/app

RUN useradd --no-log-init tp

COPY go.mod go.sum main.go ./
COPY internal ./internal
RUN go install
RUN go build

COPY --from=builder /app/build internal/server/web/

ENV GIN_MODE=release
EXPOSE 8080

USER tp

CMD ["./teleport_interview"]
