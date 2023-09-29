FROM golang:1.21.1 AS build

WORKDIR /go/src/github.com/ebauman/widgetfactory
COPY . .

ENV GOOS=linux
ENV CGO_ENABLED=0

RUN go get -d -v ./...
RUN go install -v ./...

FROM node:20-alpine AS app-build

WORKDIR /app

COPY app /app

RUN npm install
RUN npm run build

FROM alpine:latest

COPY --from=build /go/bin/widgetfactory /usr/local/bin
COPY --from=app-build /app/dist /web

ENTRYPOINT ["widgetfactory", "--static-content-path", "/web"]