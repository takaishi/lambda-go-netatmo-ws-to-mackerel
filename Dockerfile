FROM golang:1.15.7 as build
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/app
COPY . .
RUN go mod download
RUN go build .


FROM alpine:3.11
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY --from=build /go/app/netatmo-ws-to-mackerel /netatmo-ws-to-mackerel
ENTRYPOINT ["/netatmo-ws-to-mackerel"]