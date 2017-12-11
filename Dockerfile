FROM centos:7

COPY _output/main /boot

CMD ["/boot"]
