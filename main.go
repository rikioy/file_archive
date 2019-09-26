package main

import (
	"bufio"
	log "github.com/sirupsen/logrus"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/ini.v1"
)

var (
	//SrcPath 原始目录
	SrcPath string //
	//DstPath 目标目录
	DstPath string
	//FfprobeExe ffprobe路径
	FfprobeExe string
	//Cfg 配置文件
	Cfg *ini.File
	//CopyMode 是否复制文件
	CopyMode bool
	fileLog *log.Logger

)

func init() {
	var err error
	f, err := os.OpenFile("failed.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("fail to open log file, err=%v", err)
	}
	fileLog = log.New()
	fileLog.SetOutput(f)
	Cfg, err = ini.Load("config.ini")
	if err != nil {
		log.Fatalf("fail to read config file, err=%v", err)
	}
}

func main() {
	SrcPath = Cfg.Section("").Key("src_path").String()
	DstPath = Cfg.Section("").Key("dst_path").String()
	FfprobeExe = Cfg.Section("").Key("ffprobe_exe").String()
	CopyMode, _ = Cfg.Section("").Key("copy_mode").Bool()
	log.Printf("src path:%s", SrcPath)
	log.Printf("dst path:%s", DstPath)
	log.Printf("ffprobe_exe:%s", FfprobeExe)
	log.Printf("copy_mode:%t", CopyMode)

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
