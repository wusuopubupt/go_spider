/* main.go 爬虫入口程序*/
/*
modification history
--------------------
2015-11-25, by wusuopubupt, create
2016-01-14, by wusuopubupt, modify GOMAXPROCS to NumCPU
*/
/*
DESCRIPTION
网页定向抓取爬虫,利用Golang的channel和goroutine实现并发
*/
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

import (
	l4g "github.com/alecthomas/log4go"
)

import (
	conf "github.com/wusuopubupt/go_spider/src/conf"
	spider "github.com/wusuopubupt/go_spider/src/spider"
	utils "github.com/wusuopubupt/go_spider/src/utils"
)

// 默认配置文件
var (
	SPIDER_CONFIG_FILE = "spider.conf"
	SPIDER_LOGCONF_XML = "../../conf/logconf.xml"
)

// l4g的bug,需要sleep一会再让主goroutine退出,log才能flush到文件
// http://stackoverflow.com/questions/14252766/abnormal-behavior-of-log4go
func SlowExit() {
	time.Sleep(time.Second)
	os.Exit(1)
}

// 等待信号
func waitSignal(sigChan chan os.Signal) {
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
}

// 爬虫主程序
func main() {
	// l4g的配置文件
	l4g.LoadConfiguration(SPIDER_LOGCONF_XML)

	// refer : http://www.01happy.com/golang-command-line-arguments/
	var confPath string
	var logPath string
	var printVer bool

	flag.StringVar(&confPath, "c", "../../conf", "config file path")
	flag.StringVar(&logPath, "l", "../../log", "log file path")
	flag.BoolVar(&printVer, "v", false, "print version")

	flag.Parse()

	if printVer {
		utils.PrintVersion()
		os.Exit(0)
	}

	l4g.Info("Hi, dash's %s is running...\n", "go_mini_spider")

	confFile := confPath + "/" + SPIDER_CONFIG_FILE
	conf, err := conf.InitConf(confFile)
	if err != nil {
		l4g.Error("read spider config failed, err [%s]", err)
		SlowExit()
	}

	var seedUrls []string
	// read and parse json,相对路径
	b, err := ioutil.ReadFile(confPath + "/" + conf.UrlListFile)
	if err != nil {
		l4g.Error("readfile err[%s]", err)
		SlowExit()
	}
	//json to []string
	if err := json.Unmarshal(b, &seedUrls); err != nil {
		l4g.Error("parse json err[%s]", err)
		SlowExit()
	}

	//GOMAXPROCS设置
    runtime.GOMAXPROCS(runtime.NumCPU())

	// 启动爬虫
	spider := spider.NewSpider(seedUrls, conf, confPath)
	spider.Start()
	// 等待任务完成
	spider.Wait()

	time.Sleep(1 * time.Second)
}
