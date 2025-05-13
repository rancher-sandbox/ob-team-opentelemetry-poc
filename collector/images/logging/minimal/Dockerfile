ARG TARGET_ARCH=amd64

FROM registry.suse.com/bci/golang:1.23 AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o collector -ldflags="-extldflags -static -s -w" .

FROM registry.suse.com/bci/golang:1.23
RUN zypper -n install net-tools iproute2 systemd
RUN which ss
COPY --from=builder /usr/src/app/collector /usr/local/bin/collector
ENTRYPOINT ["collector"]

