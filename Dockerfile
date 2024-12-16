FROM ubuntu:22.04 AS isolate

RUN apt-get update \
    && apt-get install -y git make gcc libcap-dev libsystemd-dev pkg-config

ENV ISOLATE_VERSION=v2.0

RUN git clone --branch $ISOLATE_VERSION --depth 1 https://github.com/ioi/isolate; \
    cd isolate; \
    make install

RUN sed -i 's|^cg_root = .*|cg_root = /sys/fs/cgroup/isolate.slice/isolate.service|' /usr/local/etc/isolate

FROM golang:1.23.4-alpine3.21 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /firefly/firefly /app/cmd/app/main.go

FROM ubuntu:22.04 AS firefly

ENV ISOLATE=/firefly/isolate/bin/isolate

RUN apt-get update \
    && apt-get install -y systemd

COPY --from=isolate /isolate/isolate /firefly/isolate/bin/
COPY --from=isolate /usr/local/etc/isolate /usr/local/etc/

COPY --from=isolate /isolate/systemd/* /etc/systemd/system/

COPY --from=builder /firefly/firefly /firefly/firefly

WORKDIR /firefly

RUN apt-get install -y libcap2-bin acl

RUN apt-get install -y python3 python3-pip

RUN apt-get install -y gcc g++ make

RUN apt-get install -y libseccomp-dev

EXPOSE 8080

CMD ["/firefly/firefly"]