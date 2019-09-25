BUILD_NAME:=file_archive
BUILD_VERSION:=1.0
SOURCE:=./src/*.go

all: build_linux build_windows

build:
	go build -o ${BUILD_NAME} ${SOURCE}

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_NAME} ${SOURCE}
	cp ${SOURCE}/config.ini.example .

build_windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o ${BUILD_NAME}.exe ${SOURCE}
	cp ${SOURCE}/config.ini.example .

clean:
	rm ${BUILD_NAME}.exe
	rm ${BUILD_NAME}.exe
	rm config.ini


.PHONY: all build build_linux build_win clean