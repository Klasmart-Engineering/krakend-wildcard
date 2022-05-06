FROM ghcr.io/kl-engineering/kl-krakend-builder:0.0.1 AS builder

WORKDIR /tmp/builder
COPY * plugins/tmp

WORKDIR /tmp/builder/plugins/tmp
RUN make all

FROM ghcr.io/kl-engineering/kl-krakend:0.0.1 as krakend

USER root

COPY krakend.json /etc/krakend/krakend.json

COPY --from=builder /tmp/builder/plugins/**/*.so /opt/krakend/plugins/
RUN if [ ! x$(find /opt/krakend/plugins -prune -empty) = x/opt/krakend/plugins ]; then chmod +x /opt/krakend/plugins/*.so; fi

USER krakend
