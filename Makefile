all: binary

dist: binary

binary: numalign

numalign:
	go build -v -o numalign

clean:
	rm numalign
