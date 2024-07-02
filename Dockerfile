# wait-for-port in a container
#
# docker run --rm -it --net host bitnami/wait-for-port
#
FROM golang:1.22-bullseye as build

RUN apt-get update && apt-get install -y --no-install-recommends \
    git make upx \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make

RUN upx --ultra-brute out/wait-for-port

FROM bitnami/minideb:bullseye

COPY --from=build /go/src/app/out/wait-for-port /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/wait-for-port"]
