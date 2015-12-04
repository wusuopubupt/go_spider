package main

import (
	"flag"
	"fmt"
	"net/http"
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
		UrlListFile     *string
		OutputDirectory *string
		MaxDepth        *int
		CrawlInterval   *int
		CrawlTimeout    *int
		TargetUrl       *string
		ThreadCount     *int
	}
}

var (
	SPIDER_CONFIG_FILE = "spider.conf"
	SPIDER_LOGCONF_XML = "../conf/logconf.xml"
)

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

// check conf
func CheckConf(s *SpiderCfg) error {
	c := s.Spider
	if c.UrlListFile == nil {
		return fmt.Errorf("Spider conf item: UrlListFile is not configured")
	}
	if c.OutputDirectory == nil {
		return fmt.Errorf("Spider conf item: OutputDirectory is not configured")
	}
	if c.MaxDepth == nil {
		return fmt.Errorf("Spider conf item: MaxDepth is not configured")
	}
	if c.CrawlInterval == nil {
		return fmt.Errorf("Spider conf item: CrawlInterval is not configured")
	}
	if c.CrawlTimeout == nil {
		return fmt.Errorf("Spider conf item: CrawlTimeout is not configured")
	}
	if c.TargetUrl == nil {
		return fmt.Errorf("Spider conf item: TargetUrl is not configured")
	}
	if c.ThreadCount == nil {
		return fmt.Errorf("Spider conf item: ThreadCount is not configured")
	}
	return nil
}

func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		chFinished <- true
	}()
}

func main() {
	l4g.LoadConfiguration(SPIDER_LOGCONF_XML)

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

	confFile := confPath + "/" + SPIDER_CONFIG_FILE
	conf, err := InitConf(confFile)
	if err != nil {
		l4g.Error("read spider config failed, err [%s]", err)
		AbnormalExit()
	}
	// check conf
	if err := CheckConf(conf); err != nil {
		l4g.Error("check spider config failed, err [%s]", err)
		AbnormalExit()
	}

	fmt.Printf("urllistfile: %s", *conf.Spider.UrlListFile)

	time.Sleep(time.Second)
}
