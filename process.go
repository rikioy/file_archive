package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
)

func add(path, name, dateStr string) (err error) {
	modTime, err := time.Parse("20060102", dateStr)
	if err != nil {
		err = fmt.Errorf("%s/%s\t formt date string failed, err=%v", path, name, err)
	}
	dstPath := fmt.Sprintf("%s/%s/%s/", DstPath, modTime.Format("2006"), modTime.Format("01-02"))
	dbFile := fmt.Sprintf("%s/%s/album.db", DstPath, modTime.Format("2006"))
	if err = insert(path, dstPath, dbFile, name); err != nil {
		return
	}
	log.Infof("插入文件成功: %s%s => %s%s", path, name, dstPath, name)
	return
}

func singleAdd(filepath, date string) error {
	info, err := os.Stat(filepath)
	if err != nil {
		log.Fatalf("read file %s failed, err=%v", filepath, err)
	}
	pathEnd := len(filepath) - len(info.Name())
	path := filepath[:pathEnd]
	return add(path, info.Name(), date)
}

func deal(v fileinfo) bool {
	singleLogger := fileLog.WithFields(log.Fields{
		"file:": v.Info.Name(),
	})
	filePath := v.Path + v.Info.Name()
	fileExt := filepath.Ext(filePath)

	// 判断是否支持的文件类型
	if !Cfg.Section("type").HasKey(fileExt) {
		singleLogger.Printf("%s file type %s unsupport.", filePath, fileExt)
		return false
	}
	// 获取modtime
	var modTime time.Time
	var err error
	switch Cfg.Section("type").Key(fileExt).String() {
	case "exif":
		modTime, err = getexif(filePath)
		if err != nil {
			singleLogger.Printf("%s, get picture time failed, err=%v", filePath, err)
			return false
		}
		break
	case "mp4":
		tag, err := Probe(filePath)
		if err != nil {
			singleLogger.Printf("%s, get mp4 tag failed, err=%v", filePath, err)
			return false
		}
		modTime, err = time.Parse("2006-01-02T15:04:05.000000Z", tag.Format.Tags["creation_time"])
		if err != nil {
			singleLogger.Printf("%s, parse mp4 time failed, err=%v", filePath, err)
			return false
		}
		break
	}

	if err := add(v.Path, v.Info.Name(), modTime.Format("20060102")); err != nil {
		singleLogger.Printf("%s, insert failed, copy failed file to %s, err=%v", filePath, FailPath+"/"+v.Info.Name(), err)
		return false
	}

	return true
}

func process(path string) {
	if path[len(path)-1:] != "/" {
		path = path + "/"
	}
	files, err := listAll(path)
	if err != nil {
		log.Fatalf("read dir %s failed, err=%v", path, err)
	}
	var errCount, succCount int
	for _, v := range files {
		if deal(v) {
			succCount++
		} else {
			errCount++
			filecopy(v.Path+v.Info.Name(), FailPath+v.Info.Name())
		}
	}
	log.Infof("批量插入完成，成功%d个，失败%d个，失败记录请查看failed.log", succCount, errCount)
}
