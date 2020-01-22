FROM registry.access.redhat.com/ubi8/ubi-minimal 
COPY numalign /bin/numalign
ENTRYPOINT ["/bin/numalign"]
