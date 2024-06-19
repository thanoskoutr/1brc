BINARY=1brc

.PHONY:all
all: build

.PHONY:build
build:
	go version
	go build -ldflags "-w -s" -o ${BINARY}

.PHONY:clean
clean:
	rm -f ./${BINARY}
