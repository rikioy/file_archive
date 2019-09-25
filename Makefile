BUILD_NAME:=file_archive
BUILD_VERSION:=1.0
SOURCE:= ./*.go

build: linux

all: linux windows

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_NAME} ${SOURCE}
	cp config.ini.example config.ini

windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o ${BUILD_NAME}.exe ${SOURCE}
	cp config.ini.example config.ini

clean:
	rm ${BUILD_NAME}
	rm ${BUILD_NAME}.exe
	rm config.ini

.PHONY: all build linux windows clean