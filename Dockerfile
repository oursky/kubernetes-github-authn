FROM gliderlabs/alpine:3.4

RUN apk --update add ca-certificates

COPY _output/main /main

CMD ["/main"]
