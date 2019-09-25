package main

import (
	"log"
	"os"
	"bufio"

	"github.com/urfave/cli"
	"gopkg.in/ini.v1"
)

var (
	//SrcPath 原始目录
	SrcPath    string // 
	//DstPath 目标目录
	DstPath    string
	//FfprobeExe ffprobe路径
	FfprobeExe string
	//Cfg 配置文件
	Cfg         *ini.File
	//MoveFile 是否移动文件
	MoveFile bool
)

func init() {
	f, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE, 0755)
	log.SetOutput(f)
	Cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Fatalf("fail to read file, err=%v\n", err)
	}
}

func main() {
	SrcPath = Cfg.Section("").Key("src_path").String()
	DstPath = Cfg.Section("").Key("dst_path").String()
	FfprobeExe = Cfg.Section("").Key("ffprobe_exe").String()
	MoveFile, _= Cfg.Section("").Key("move_file").Bool()
	log.Printf("src path:%s\n", SrcPath)
	log.Printf("dst path:%s\n", DstPath)
	log.Printf("ffprobe_exe:%s\n", FfprobeExe)
	log.Printf("move_file:%t\n", MoveFile)

	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a single file to dst",
			Flags: []cli.Flag{
				cli.StringFlag{Name: "file, f"},
				cli.StringFlag{Name: "date, d"},
			},
			Action: func(c *cli.Context) error {
				err := add(c.String("file"), c.String("date"))
				return err
			},
		},
		{
			Name:    "process",
			Aliases: []string{"p"},
			Usage:   "process the folder file to dst",
			Action: func(c *cli.Context) error {
				process()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	inputReader := bufio.NewReader(os.Stdin)
	inputReader.ReadString('\n')
}
