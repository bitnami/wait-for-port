# wait-for-port in a container
#
# docker run --rm -it --net host bitnami/wait-for-port
#
FROM bitnami/golang:1.25 as build

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make

FROM bitnami/minideb:bookworm

COPY --from=build /go/src/app/out/wait-for-port /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/wait-for-port"]
