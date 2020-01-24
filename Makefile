all: binary

dist: binary

binary: numalign

numalign:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -o numalign

clean:
	rm numalign
