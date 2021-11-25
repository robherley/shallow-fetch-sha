FROM golang:1.17 as builder

WORKDIR /build

ADD ./go.mod /build/go.mod
ADD ./go.sum /build/go.sum
RUN go mod download
RUN go mod verify

COPY main.go main.go
COPY internal/ internal/

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o sfs .

FROM alpine

LABEL maintainer="Rob Herley <robherley13@gmail.com>"

RUN apk --update --no-cache add bash git openssh ca-certificates
RUN update-ca-certificates

COPY --from=builder --chown=1001:1001 /build/sfs /usr/local/bin/sfs

COPY --chown=1001:1001 ./script/generate_known_hosts /usr/local/bin/generate_known_hosts
RUN cd /tmp && generate_known_hosts && mv ssh_known_hosts /etc/ssh/ssh_known_hosts

RUN adduser -Du 1001 sfs-user
USER sfs-user

WORKDIR /usr/src/repo

ENTRYPOINT ["/usr/local/bin/sfs"]
CMD ["--help"]