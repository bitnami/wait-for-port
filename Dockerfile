# wait-for-port in a container
#
# docker run --rm -it --net host bitnami/wait-for-port
#
FROM ryotakatsuki/godev as build

WORKDIR /go/src/app
COPY . .

RUN rm -rf out

RUN make

RUN upx --ultra-brute out/wait-for-port

FROM alpine:latest

COPY --from=build /go/src/app/out/wait-for-port /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/wait-for-port"]
