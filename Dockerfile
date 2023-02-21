FROM --platform=$TARGETPLATFORM scratch as runner

LABEL maintainer="hedongshu <hedongshu@foxmail.com>"

ARG TARGETOS
ARG TARGETARCH

COPY ./build/markdown-book-${TARGETOS}-${TARGETARCH}/markdown-book /data/app/

EXPOSE 5006

ENTRYPOINT ["/data/app/markdown-book", "web"]
