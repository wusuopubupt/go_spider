package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

import (
	"code.google.com/p/gcfg"
	l4g "code.google.com/p/log4go"
	"golang.org/x/net/html"
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

type Urls struct {
	url []string
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

// get href attribute from a Token
func GetHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
		l4g.Debug("get href: %s", href)
	}
	// return默认返回ok, href
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
		fmt.Println("ERROR: Failed to crawl \"" + url + "\"")
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

func main1() {
	//创建映射(数组), key type:string, value type:bool
	foundUrls := make(map[string]bool)
	seedUrls := os.Args[1:]

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

	// We're done! Print the results...

	fmt.Println("\nFound", len(foundUrls), "unique urls:\n")

	for url, _ := range foundUrls {
		fmt.Println(" - " + url)
	}

	close(chUrls)
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

	l4g.Debug("urllistfile: %s", *conf.Spider.UrlListFile)
	// read and parse json
	b, err := ioutil.ReadFile(*conf.Spider.UrlListFile)
	if err != nil {
		l4g.Error("readfile err[%s]", err)
		AbnormalExit()
	}
	//json 到 []string
	var urls []string
	if err := json.Unmarshal(b, &urls); err != nil {
		l4g.Error("parse json err[%s]", err)
		AbnormalExit()
	}
	l4g.Debug("urls: %s", urls)

	time.Sleep(1 * time.Second)
}
