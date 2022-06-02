FROM alpine:3.16 AS base

RUN apk add --no-cache --upgrade \
        ca-certificates \
        tzdata

RUN mkdir -v /tmp/rootfs \
    && tar -cf- /usr/share/zoneinfo | tar -xf- -C /tmp/rootfs \
    && mkdir -v /tmp/rootfs/etc \
    && echo 'nobody:x:65534:65534:nobody:/:' > /tmp/rootfs/etc/passwd \
    && echo 'nobody:x:65534:' > /tmp/rootfs/etc/group \
    && install -vDm 644 /etc/ssl/certs/ca-certificates.crt /tmp/rootfs/etc/ssl/certs/ca-certificates.crt

FROM scratch

COPY --from=base /tmp/rootfs /
COPY docker-reloader /docker-reloader

WORKDIR /

USER nobody:nobody

ENTRYPOINT ["/docker-reloader"]
