package main

import (
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
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

func changename(path string) string {
	if !exists(path) {
		return path
	}
	fileExt := filepath.Ext(path)
	pathlen := len(path)
	subpath := path[:pathlen-len(fileExt)]
	randStr := createRandomString(4)
	newpath := fmt.Sprintf("%s-%s%s", subpath, randStr, fileExt)
	return changename(newpath)
}

func filecopy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	dst = changename(dst)

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
			files, _ := listAll(path + info.Name() + Sep)
			infos = append(infos, files...)
		} else {
			var tmp = fileinfo{}
			tmp.Path = path + Sep
			tmp.Info = info
			infos = append(infos, tmp)
		}
	}
	return
}

func createRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}
