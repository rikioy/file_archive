package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/rwcarlsen/goexif/exif"
)

type fileinfo struct {
	Path string
	Info os.FileInfo
}

func filemd5(path string) (md5Str string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		err = fmt.Errorf("md5计算失败, err=%v", err)
		return
	}
	md5Str = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func filecopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func getexif(path string) (t time.Time, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	defer file.Close()

	x, err := exif.Decode(file)
	if err != nil {
		fmt.Println("failed to decode file, err: ", err)
		return
	}
	t, err = x.DateTime()
	return
}

func listAll(path string) (infos []fileinfo, err error) {
	readInfos, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, info := range readInfos {
		if info.IsDir() {
			files, _ := listAll(path + info.Name() + "/")
			infos = append(infos, files...)
		} else {
			var tmp = fileinfo{}
			tmp.Path = path
			tmp.Info = info
			infos = append(infos, tmp)
		}
	}
	return
}
