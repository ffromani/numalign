all: dist

dist: binaries

binaries: numalign sriovscan

numalign:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o numalign ./cmd/numalign

sriovscan:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o sriovscan ./cmd/sriovscan

clean:
	rm numalign sriovscan
