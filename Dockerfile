FROM golang:1.17 as builder

WORKDIR /build

ADD ./go.mod /build/go.mod
ADD ./go.sum /build/go.sum
RUN go mod download
RUN go mod verify

COPY main.go main.go
COPY internal/ internal/
COPY script/ script/

RUN script/build -o sfs
RUN script/generate_known_hosts

FROM alpine

LABEL maintainer="Rob Herley <robherley13@gmail.com>"

RUN apk --update --no-cache add bash git openssh ca-certificates
RUN update-ca-certificates

COPY --from=builder --chown=1001:1001 /build/sfs /usr/local/bin/sfs
COPY --from=builder --chown=1001:1001 /build/ssh_known_hosts /etc/ssh/ssh_known_hosts

RUN adduser -Du 1001 sfs-user
USER sfs-user

WORKDIR /usr/src/repo

ENTRYPOINT ["/usr/local/bin/sfs"]
CMD ["--help"]