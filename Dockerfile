FROM golang:alpine as builder

WORKDIR /build

ADD ./go.mod /build/go.mod
ADD ./go.sum /build/go.sum
RUN go mod download
RUN go mod verify

COPY main.go main.go
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o shallow-fetch-sha .

FROM alpine

LABEL maintainer="Rob Herley <robherley13@gmail.com>"

RUN apk --update --no-cache add git openssh ca-certificates
RUN update-ca-certificates

COPY --from=builder /build/shallow-fetch-sha /usr/local/bin/shallow-fetch-sha

RUN chmod a+rx /usr/local/bin/shallow-fetch-sha

USER guest

WORKDIR /usr/src/repo

ENTRYPOINT ["/usr/local/bin/shallow-fetch-sha"]
CMD ["--help"]