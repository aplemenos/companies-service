# Initial stage: download modules
FROM golang:1.20-alpine as builder

ENV config=docker

WORKDIR /app

COPY ./ /app

RUN go mod download


# Intermediate stage: Build the binary
FROM golang:1.20-alpine as runner

COPY --from=builder ./app ./app

RUN apk --no-cache add curl
RUN apk --no-cache add git

RUN go install github.com/githubnemo/CompileDaemon@v1.4.0

WORKDIR /app
ENV config=docker

EXPOSE 5000
EXPOSE 5555
EXPOSE 7070

ENTRYPOINT CompileDaemon --build="go build cmd/companies-api/main.go" --command=./main