FROM alpine:3.9
COPY _output /usr/local/bin
ENTRYPOINT ["/bin/sh"]
