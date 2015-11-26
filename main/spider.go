package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

import (
	"code.google.com/p/gcfg"
	l4g "code.google.com/p/log4go"
)

import (
	utils "github.com/wusuopubupt/go_spider/utils"
)

type SpiderCfg struct {
	Spider struct {
		UrlListFile     string
		OutputDirectory string
		MaxDepth        int
		CrawlInterval   int
		CrawlTimeout    int
		TargetUrl       string
		ThreadCount     int
	}
}

// abnormal exit
func AbnormalExit() {
	// http://stackoverflow.com/questions/14252766/abnormal-behavior-of-log4go
	// adding a time.Sleep(time.Second) to the end of the code snippeet will cause the log content flush
	time.Sleep(time.Second)
	os.Exit(1)
}

func InitConf(confFile string) (*SpiderCfg, error) {
	l4g.Info("config file: %s", confFile)
	var cfg SpiderCfg
	err := gcfg.ReadFileInto(&cfg, confFile)

	if err != nil {
		l4g.Error("read config err [%s]", err)
		return nil, err
	}
	return &cfg, nil
}

func main() {
	l4g.LoadConfiguration("logconf.xml")

	// refer : http://www.01happy.com/golang-command-line-arguments/
	// 方法一： flag.StringVar(),传入指针，直接给confPath赋值
	var confPath string
	var printVer bool
	flag.StringVar(&confPath, "c", "../conf", "config file path")
	flag.BoolVar(&printVer, "v", false, "print version")
	// 方法二：flag.String()，把函数调用的返回值赋值给logPath
	//logPath := flag.String("l", "../log", "log file path")

	flag.Parse()

	if printVer {
		utils.PrintVersion()
	}

	l4g.Info("Hi, dash's %s is running...\n", "go_mini_spider")

	conf, err := InitConf(confPath + "/spider.conf")
	if err != nil {
		l4g.Error("rend spider config failed !")
		AbnormalExit()
	}

	fmt.Printf("urllistfile: %s", conf.Spider.UrlListFile)
}
