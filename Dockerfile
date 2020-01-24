FROM alpine:3.8
COPY numalign /bin/numalign
ENTRYPOINT ["/bin/numalign"]
