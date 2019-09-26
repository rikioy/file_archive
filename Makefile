PROJECTNAME=$(shell basename "$(PWD)")
BASE=$(shell pwd)
BUILD_VERSION:=1.0
SOURCE:= $(wildcard *.go)
DIST_DIR:=$(BASE)/dist
LINUX_DIST_NAME:=file_archive_linux_amd64
WINDOWS_DIST_NAME:=file_archive_windows_amd64

TMP_DIR:=$(BASE)/tmp
LINUX_TMP_DIR:=$(TMP_DIR)/${LINUX_DIST_NAME}
WIN_TMP_DIR:=$(TMP_DIR)/${WINDOWS_DIST_NAME}

build: linux



all: linux windows

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${BUILD_NAME} ${SOURCE}
	cp config.ini.example config.ini

windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o ${BUILD_NAME}.exe ${SOURCE}
	cp config.ini.example config.ini

clean:
	@rm -rf ${DIST_DIR}
	@rm -rf ${TMP_DIR}

dist: clean
	@echo " > Building binary..."
	@mkdir ${DIST_DIR}
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${LINUX_TMP_DIR}/${PROJECTNAME} ${SOURCE}
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o ${WIN_TMP_DIR}/${PROJECTNAME}.exe ${SOURCE}
	@cp ${BASE}/config.ini.example ${LINUX_TMP_DIR}/config.ini
	@cp ${BASE}/config.ini.example ${WIN_TMP_DIR}/config.ini
	@tar zcf ${DIST_DIR}/${LINUX_DIST_NAME}.tar.gz -C ${LINUX_TMP_DIR}/ .
	@zip -rj ${DIST_DIR}/${WINDOWS_DIST_NAME}.zip ${WIN_TMP_DIR}/*

.PHONY: all build linux windows clean dist