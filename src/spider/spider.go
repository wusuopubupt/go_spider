package spider

import (
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
)

// abnormal exit
func AbnormalExit() {
	// http://stackoverflow.com/questions/14252766/abnormal-behavior-of-log4go
	// adding a time.Sleep(time.Second) to the end of the code snippeet will cause the log content flush
	time.Sleep(time.Second)
	os.Exit(1)
}

// get href attribute from a Token
func (s *Spider) getHref(t html.Token) (ok bool, href string) {
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
// which do current job and add new jobs to job queue
func (s *Spider) crawl(jobs JobQueue, chFinished chan bool) {
	job := s.getJob(jobs)
	l4g.Info("get job: %s", job.url)
	resp, err := http.Get(job.url)

	// Notify that we're done after this function
	defer func() {
		chFinished <- true
	}()
	if err != nil {
		l4g.Error("Failed to crawl %s, err[%s]", job.url, err)
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
			ok, url := s.getHref(t)
			if !ok {
				continue
			}
			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			if hasProto {
				jobs.url <- url
				jobs.depth <- job.depth + 1
				l4g.Info("add job: %s", url)
			}
		}
	}
	// 抓取间隔控制
	time.Sleep(time.Duration(s.crawlInterval) * time.Second)
}

// Crawler struct
type Spider struct {
	outputDir     string
	crawlInterval int
	crawlTimeout  int
	targetUrl     string
}

// one job
type Job struct {
	url   string
	depth int
}

// job queue
type JobQueue struct {
	url   chan string
	depth chan int
}

// Channels
/*
c := make(chan bool) //创建一个无缓冲的bool型Channel
c <- x //向一个Channel发送一个值
<- c //从一个Channel中接收一个值
x = <- c //从Channel c接收一个值并将其存储到x中
x, ok = <- c //从Channel接收一个值，如果channel关闭了或没有数据，那么ok将被置为false
*/

// get job from jobQueue
func (s *Spider) getJob(jobs JobQueue) (job Job) {
	job.url = <-jobs.url
	job.depth = <-jobs.depth
	return job
}

// add job to jobQueue
func (s *Spider) addJob(jobs JobQueue, job Job) {
	jobs.url <- job.url
	jobs.depth <- job.depth
}

// new spider
func newSpider(config conf.SpiderStruct) *Spider {
	s := new(Spider)
	s.outputDir = config.OutputDirectory
	s.crawlInterval = config.CrawlInterval
	s.crawlTimeout = config.CrawlTimeout
	s.targetUrl = config.TargetUrl

	return s
}

// 开启threandCount个spider goroutine,等待通道中的任务到达
func Start(seedUrls []string, config conf.SpiderStruct) {
	var jobs JobQueue
	var spiders []*Spider
	chFinished := make(chan bool)
	// 创建threadCount个工作goroutine
	for i := 0; i < config.ThreadCount; i++ {
		s := newSpider(config)
		spiders = append(spiders, s)
		l4g.Info("created new spider #%d", i)
		//go s.crawl(jobs, chFinished)
	}
	// 一个while(1)的循环，直到channel通知任务结束
	for {
		for i, s := range spiders {
			l4g.Info("spider #%d is running", i)
			go s.crawl(jobs, chFinished)
		}
		// 初始化任务队列
		for _, url := range seedUrls {
			l4g.Info("url: %s", url)
			jobs.url <- url
			jobs.depth <- 0
		}
		for done := 0; done < config.ThreadCount; {
			// 通知主goroutine任务结束
			<-chFinished
			done++
		}
	}
}
