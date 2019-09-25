package main

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rwcarlsen/goexif/exif"
	"io"
	"os"
	"time"
)

func insert(srcPath, dstPath, dbFile, name string) (err error) {
	srcFile := srcPath + "/" + name
	dstFile := dstPath + "/" + name

	srcMd5 , err := filemd5(srcFile)
	if err != nil {
		err = fmt.Errorf("%v, %s", err, srcFile)
		return
	}
	if !exists(dstPath) {
		err = os.MkdirAll(dstPath, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("创建文件夹失败, %s ", dstPath)
			return
		}
	}

	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		err = fmt.Errorf("创建数据库链接, %s/%s/%s ", dbFile)
		return
	}

	sqlTable := `
		CREATE TABLE IF NOT EXISTS album (
			md5 char(32) primary key,
			filename varchar(100),
			srcPath varchar(200)
		);
`
	_, err = db.Exec(sqlTable)
	if err != nil {
		return
	}

	var count int
	stat := `SELECT count(*) As count FROM album WHERE md5=$1;`
	row := db.QueryRow(stat, srcMd5)
	err = row.Scan(&count)
	if err != nil {
		return
	}

	if count > 0 {
		err = fmt.Errorf("已经存在相同MD5的文件, srcfile:%s\t dstfile:%s ", srcFile, dstFile)
		return
	}
	_, err = filecopy(srcFile, dstFile)
	if err != nil {
		err = fmt.Errorf("复制文件失败, srcfile:%s\t dstfile:%s ", srcFile, dstFile)
		return
	}

	dstMd5, err := filemd5(dstFile)
	if err != nil {
		err = fmt.Errorf("生成目标md5失败, %s ", dstFile)
		return
	}

	if dstMd5 != srcMd5 {
		err = fmt.Errorf("复制失败，目标md5和源md5不一致, srcfile:%s\t dstfile:%s ", srcFile, dstFile)
		return
	}

	stmt, err := db.Prepare("INSERT INTO album(md5, filename, srcPath) values(?,?,?)")
	_, err = stmt.Exec(srcMd5, name, dstFile)
	if err != nil {
		err = fmt.Errorf("插入数据库失败, %s, err=%v", dstFile, err)
		return
	}
	return
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
