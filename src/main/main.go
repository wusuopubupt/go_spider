package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"time"
)

import (
	l4g "code.google.com/p/log4go"
)

import (
	conf "github.com/wusuopubupt/go_spider/src/conf"
	spider "github.com/wusuopubupt/go_spider/src/spider"
	utils "github.com/wusuopubupt/go_spider/src/utils"
)

var (
	SPIDER_CONFIG_FILE = "spider.conf"
	SPIDER_LOGCONF_XML = "../../conf/logconf.xml"
)

// abnormal exit
func AbnormalExit() {
	// http://stackoverflow.com/questions/14252766/abnormal-behavior-of-log4go
	// adding a time.Sleep(time.Second) to the end of the code snippeet will cause the log content flush
	time.Sleep(time.Second)
	os.Exit(1)
}

// 爬虫主程序
func main() {
	l4g.LoadConfiguration(SPIDER_LOGCONF_XML)

	// refer : http://www.01happy.com/golang-command-line-arguments/
	// 方法一： flag.StringVar(),传入指针，直接给confPath赋值
	var confPath string
	var printVer bool
	flag.StringVar(&confPath, "c", "../../conf", "config file path")
	flag.BoolVar(&printVer, "v", false, "print version")
	// 方法二：flag.String()，把函数调用的返回值赋值给logPath
	//logPath := flag.String("l", "../log", "log file path")

	flag.Parse()

	if printVer {
		utils.PrintVersion()
	}

	l4g.Info("Hi, dash's %s is running...\n", "go_mini_spider")

	confFile := confPath + "/" + SPIDER_CONFIG_FILE
	conf, err := conf.InitConf(confFile)
	if err != nil {
		l4g.Error("read spider config failed, err [%s]", err)
		AbnormalExit()
	}

	var seedUrls []string
	// read and parse json,相对路径
	b, err := ioutil.ReadFile(confPath + "/" + conf.UrlListFile)
	if err != nil {
		l4g.Error("readfile err[%s]", err)
		AbnormalExit()
	}
	//json to []string
	if err := json.Unmarshal(b, &seedUrls); err != nil {
		l4g.Error("parse json err[%s]", err)
		AbnormalExit()
	}
	seedUrls = []string{"www.baidu.com", "www.sina.com.cn"}
	l4g.Debug("seedUrls: %s", seedUrls)

	// start miniSpider
	spider.Start(seedUrls, conf)

	time.Sleep(1 * time.Second)
}
