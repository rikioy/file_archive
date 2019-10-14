package main

import (
	"fmt"
	"os"

	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func insert(srcPath, dstPath, dbFile, name string) (realpath string, err error) {
	srcFile := srcPath + Sep + name
	dstFile := dstPath + Sep + name

	srcMd5, err := filemd5(srcFile)
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
		err = fmt.Errorf("创建数据库链接, %s", dbFile)
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
		//fmt.Printf("已经存在相同MD5的文件, srcfile:%s dstfile:%s\n", srcFile, dstFile)
		return
	}

	realpath, err = filecopy(srcFile, dstFile)
	if err != nil {
		err = fmt.Errorf("复制文件失败, srcfile:%s dstfile:%s, err=%v", srcFile, realpath, err)
		return
	}

	dstMd5, err := filemd5(realpath)
	if err != nil {
		err = fmt.Errorf("生成目标md5失败, %s, err=%v", realpath, err)
		return
	}

	if dstMd5 != srcMd5 {
		err = fmt.Errorf("复制失败，目标md5和源md5不一致, srcfile:%s\t dstfile:%s ", srcFile, realpath)
		return
	}

	stmt, err := db.Prepare("INSERT INTO album(md5, filename, srcPath) values(?,?,?)")
	_, err = stmt.Exec(srcMd5, name, realpath)
	if err != nil {
		err = fmt.Errorf("插入数据库失败, %s, err=%v", realpath, err)
		return
	}
	return
}
