FROM alpine:3.9
COPY numalign /bin/numalign
COPY sriovscan /bin/sriovscan
ENTRYPOINT ["/bin/sh"]
