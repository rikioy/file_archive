package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"
)

func add(name, dateStr string) (err error) {
	modTime, err := time.Parse("20060102", dateStr)
	if err != nil {
		err = fmt.Errorf("%s/%s\t formt date string failed, err=%v", SrcPath, name, err)
	}
	dstPath := fmt.Sprintf("%s/%s/%s/", DstPath, modTime.Format("2006"), modTime.Format("01-02"))
	dbFile := fmt.Sprintf("%s/%s/album.db", DstPath, modTime.Format("2006"))
	if err = insert(SrcPath, dstPath, dbFile, name); err != nil {
		return
	 }
	 fmt.Printf("插入文件成功: %s/%s => %s/%s\n", SrcPath, name, dstPath, name)
	return
}

func process() {
	files, err := ioutil.ReadDir(SrcPath)
	if err != nil {
		log.Fatalf("read dir %s failed, err=%v\n", SrcPath, err)
	}
	var errCount int
	for _, v := range files {
		filePath := SrcPath + "/" + v.Name()
		fileExt := filepath.Ext(filePath)
		if !Cfg.Section("type").HasKey(fileExt) {
			log.Printf("%s file type %s unsupport.\n", filePath, fileExt)
			errCount++
			continue
		}
		var modTime time.Time
		switch Cfg.Section("type").Key(fileExt).String() {
		case "exif":
			modTime, err = getexif(filePath)
			if err != nil {
				log.Printf("%s, get picture time failed, err=%v\n", filePath, err)
				errCount++
				continue
			}
			break
		case "mp4":
			tag, err := Probe(filePath)
			if err != nil {
				log.Printf("%s, get mp4 tag failed, err=%v\n", filePath, err)
				errCount++
				continue
			}
			modTime, err = time.Parse("2006-01-02T15:04:05.000000Z", tag.Format.Tags["creation_time"])
			if err != nil {
				log.Printf("%s, parse mp4 time failed, err=%v\n", filePath, err)
				errCount++
				continue
			}
			break
		}
		
		 if err := add(v.Name(), modTime.Format("20060102")); err != nil {
			log.Printf("%s, insert failed, err=%v\n", filePath, err)
			errCount++
		 }
	}
	fmt.Printf("批量插入完成，失败%d个，请查看log.txt\n", errCount)
}
