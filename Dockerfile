FROM alpine:3.9
COPY numalign /bin/numalign
ENTRYPOINT ["/bin/numalign"]
