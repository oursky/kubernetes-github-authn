FROM alpine:3.11.5

COPY _output/main /boot

CMD ["/boot"]
