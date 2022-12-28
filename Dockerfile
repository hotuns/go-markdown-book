FROM --platform=$TARGETPLATFORM scratch as runner

MAINTAINER willgao <will-gao@hotmail.com>

ARG TARGETOS
ARG TARGETARCH

COPY ./build/markdown-book-${TARGETOS}-${TARGETARCH}/markdown-book /data/app/

EXPOSE 5006

ENTRYPOINT ["/data/app/markdown-book", "web"]
