RUNTIME ?= podman
REPOOWNER ?= fromani
IMAGENAME ?= numalign
IMAGETAG ?= latest

BUILDFLAGS=GO111MODULE=on GOPROXY=off GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 CGO_ENABLED=0

all: dist

outdir:
	mkdir -p _output || :

.PHONY: dist
dist: binaries

.PHONY: binaries
binaries: numalign sriovscan lsnt splitcpulist sriovctl

numalign: outdir
	$(BUILDFLAGS) go build -v -o _output/numalign ./cmd/numalign

sriovscan: outdir
	$(BUILDFLAGS) go build -v -o _output/sriovscan ./cmd/sriovscan

lsnt: outdir
	$(BUILDFLAGS) go build -v -o _output/lsnt ./cmd/lsnt

splitcpulist: outdir
	$(BUILDFLAGS) go build -v -o _output/splitcpulist ./cmd/splitcpulist

sriovctl: outdir
	$(BUILDFLAGS) go build -v -o _output/sriovctl ./cmd/sriovctl

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
