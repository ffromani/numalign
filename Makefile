
all: dist

outdir:
	mkdir -p _output || :

dist: binaries

binaries: numalign sriovscan lsnt

numalign: outdir
	GO111MODULE=on GOPROXY=off GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o _output/numalign ./cmd/numalign

sriovscan: outdir
	GO111MODULE=on GOPROXY=off GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o _output/sriovscan ./cmd/sriovscan

lsnt: outdir
	GO111MODULE=on GOPROXY=off GOFLAGS=-mod=vendor GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o _output/lsnt ./cmd/lsnt

clean:
	rm -rf _output
