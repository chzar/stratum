BIN_DIR=bin
BIN_NAME=stratum

build:
	go build -o ${BIN_DIR}/${BIN_NAME} main.go

run:
	go run main.go

clean:
	rm -r ${BIN_DIR}