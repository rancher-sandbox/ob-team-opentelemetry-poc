ARG TARGET_ARCH=amd64
ARG OCB_VERSION=0.120.0

# FROM registry.suse.com/bci/golang:1.23 AS ocb-getter
# ARG OCB_VERSION
# ARG TARGET_ARCH
# RUN zypper -n install wget
# RUN wget https://github.com/open-telemetry/opentelemetry-collector-releases/releases/download/cmd%2Fbuilder%2Fv${OCB_VERSION}/ocb_${OCB_VERSION}_linux_${TARGET_ARCH}
# RUN mv -v ocb_${OCB_VERSION}_linux_amd64 /usr/local/bin/ocb
# RUN chmod +x /usr/local/bin/ocb
# RUN /usr/local/bin/ocb version

FROM registry.suse.com/bci/golang:1.23 AS builder
WORKDIR /usr/src/app
# COPY --from=ocb-getter /usr/local/bin/ocb /usr/local/bin/ocb
# RUN ocb version
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o collector -ldflags="-extldflags -static -s -w" .

FROM registry.suse.com/bci/golang:1.23
RUN zypper -n install net-tools iproute2
RUN which ss
COPY --from=builder /usr/src/app/collector /usr/local/bin/collector
ENTRYPOINT ["collector"]

