RUNTIME ?= podman
REPOOWNER ?= fromani
IMAGENAME ?= numalign
IMAGETAG ?= latest

all: dist

outdir:
	mkdir -p _output || :

.PHONY: dist
dist: binaries

.PHONY: binaries
binaries: outdir
	# go flags are set in here
	./hack/build-binaries.sh

clean:
	rm -rf _output

.PHONY: image
image: binaries
	@echo "building image"
	$(RUNTIME) build -f Dockerfile -t quay.io/$(REPOOWNER)/$(IMAGENAME):$(IMAGETAG) .

.PHONY: push
push: image
	@echo "pushing image"
	$(RUNTIME) push quay.io/$(REPOOWNER)/$(IMAGENAME):$(IMAGETAG)
