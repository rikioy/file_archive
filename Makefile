PROJECTNAME=$(shell basename "$(PWD)")
BASE=$(shell pwd)
BUILD_VERSION:=1.0
DIST_DIR:=$(BASE)/dist
BINARY_NAME=file_archive
LINUX_DIST_NAME:=$(BINARY_NAME)_linux_amd64
WINDOWS_DIST_NAME:=$(BINARY_NAME)_windows_amd64.exe

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean

TMP_DIR:=$(BASE)/tmp
LINUX_TMP_DIR:=$(TMP_DIR)/${LINUX_DIST_NAME}
WIN_TMP_DIR:=$(TMP_DIR)/${WINDOWS_DIST_NAME}

build: linux

all: linux windows

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o ${LINUX_DIST_NAME}

windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc $(GOBUILD) -o ${WINDOWS_DIST_NAME}

clean:
	@rm -rf ${DIST_DIR}
	@rm -rf ${TMP_DIR}

dist: clean
	@echo " > Building binary..."
	@mkdir ${DIST_DIR}
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ${LINUX_TMP_DIR}/${PROJECTNAME}
	@CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o ${WIN_TMP_DIR}/${PROJECTNAME}.exe
	@tar zcf ${DIST_DIR}/${LINUX_DIST_NAME}.tar.gz -C ${LINUX_TMP_DIR}/ .
	@zip -rj ${DIST_DIR}/${WINDOWS_DIST_NAME}.zip ${WIN_TMP_DIR}/*

.PHONY: all build linux windows clean dist