package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

import (
	l4g "code.google.com/p/log4go"
	"golang.org/x/net/html"
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

// get href attribute from a Token
func GetHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
		l4g.Debug("get href: %s", href)
	}
	// 空的return默认返回ok, href
	return
}

// Extract all http** links from a given webpage
func crawl(url string, ch chan string, chFinished chan bool) {
	resp, err := http.Get(url)

	defer func() {
		// Notify that we're done after this function
		chFinished <- true
	}()
	if err != nil {
		l4g.Error("Failed to crawl %s, err[%s]", url, err)
		return
	}

	b := resp.Body
	defer b.Close() // close Body when the function returns

	z := html.NewTokenizer(b)

	for {
		tt := z.Next()
		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			return
		case tt == html.StartTagToken:
			t := z.Token()
			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}
			// Extract the href value, if there is one
			ok, url := GetHref(t)
			if !ok {
				continue
			}
			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				ch <- url
			}
		}
	}
}

/**
* @brief 爬取url
* @param seedUrls 种子urls数组
*
 */

func GetUrls(seedUrls []string) {
	// Channels
	/*
		c := make(chan bool) //创建一个无缓冲的bool型Channel
		c <- x //向一个Channel发送一个值
		<- c //从一个Channel中接收一个值
		x = <- c //从Channel c接收一个值并将其存储到x中
		x, ok = <- c //从Channel接收一个值，如果channel关闭了或没有数据，那么ok将被置为false
	*/
	chUrls := make(chan string)
	chFinished := make(chan bool)
	foundUrls := make(map[string]bool)

	// Kick off the crawl process (concurrently)
	for _, url := range seedUrls {
		go crawl(url, chUrls, chFinished)
	}

	// Subscribe to both channels
	for c := 0; c < len(seedUrls); {
		// 监听 IO 操作
		select {
		case url := <-chUrls:
			foundUrls[url] = true
		case <-chFinished:
			c++
		}
	}

	l4g.Info("Found %d unique urls", len(foundUrls))
	for url, _ := range foundUrls {
		l4g.Info(" - " + url)
	}
	// close channel
	//chUrls <- "www.baidu.com"
	//url := <-chUrls
	close(chUrls)
}

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
	l4g.Debug("seedUrls: %s", seedUrls)

	// start miniSpider
	spider.Start(seedUrls, conf)

	// get urls
	//GetUrls(seedUrls)

	time.Sleep(1 * time.Second)
}
