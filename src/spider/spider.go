package spider

import (
	"io"
	"net/http"
	"net/url"
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

// Crawler struct
type Spider struct {
	// 多个goroutine共享的属性
	outputDir     string
	crawlInterval int
	crawlTimeout  int
	targetUrl     string
	jobs          JobQueue
	visitedUrl    map[string]bool
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
		//l4g.Debug("get href: %s", href)
	}
	// 空的return默认返回ok, href
	return
}

// parse html
func (s *Spider) parseHtml(b io.Reader, job Job) {
	z := html.NewTokenizer(b)
	for {
		tokenType := z.Next()
		switch {
		case tokenType == html.ErrorToken:
			l4g.Debug("end of page: %s,\tstart to get next job", job.url)
			return
		case tokenType == html.StartTagToken:
			token := z.Token()
			if !(token.Data == "a") {
				continue
			}
			ok, link := s.getHref(token)
			if !ok {
				continue
			}
			// Make sure the url begines in http**
			hasProto := strings.Index(link, "http") == 0
			u, _ := url.Parse(link)
			realUrl := u.Scheme + "://" + u.Host + u.Path
			if !s.visitedUrl[realUrl] && hasProto {
				s.jobs.url <- realUrl
				s.jobs.depth <- job.depth + 1
				l4g.Info("add job: %s, depth:%d", realUrl, job.depth+1)
			}
		}
	}
}

// crawl
func (s *Spider) crawl(chFinished chan bool) {
	// Notify that we're done after this function
	defer func() {
		chFinished <- true
	}()
	for {
		job := s.getJob()
		l4g.Info("get job url:%s, depth:%d", job.url, job.depth)
		// 检查是否访问过
		if s.visitedUrl[job.url] {
			l4g.Info("visted job,continue. url:%s, depth:%d", job.url, job.depth)
			continue
		}
		s.visitedUrl[job.url] = true
		resp, err := http.Get(job.url)
		if err != nil {
			l4g.Error("Failed to crawl %s, err[%s]", job.url, err)
			return
		}
		defer resp.Body.Close()
		s.parseHtml(resp.Body, job)
		// 抓取间隔控制
		time.Sleep(time.Duration(s.crawlInterval) * time.Second)
	}
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
func (s *Spider) getJob() (job Job) {
	job.url = <-s.jobs.url
	job.depth = <-s.jobs.depth
	return job
}

// add job to jobQueue
func (s *Spider) addJob(jobs JobQueue, job Job) {
	jobs.url <- job.url
	jobs.depth <- job.depth
}

// new spider
func newSpider(config conf.SpiderStruct, jobs JobQueue) *Spider {
	s := new(Spider)
	s.outputDir = config.OutputDirectory
	s.crawlInterval = config.CrawlInterval
	s.crawlTimeout = config.CrawlTimeout
	s.targetUrl = config.TargetUrl
	s.jobs = jobs
	s.visitedUrl = make(map[string]bool)

	return s
}

func Start(seedUrls []string, config conf.SpiderStruct) {
	//var spiders []*Spider
	var jobs JobQueue
	jobs.url = make(chan string, 1000000)
	jobs.depth = make(chan int, 1000000)
	chFinished := make(chan bool)
	// 初始化任务队列
	for _, url := range seedUrls {
		l4g.Info("url: %s", url)
		jobs.url <- url
		jobs.depth <- 0
	}
	// 一个while(1)的循环，直到channel通知任务结束
	for {
		s := newSpider(config, jobs)
		// 开启threandCount个spider.crawl goroutine,等待通道中的任务到达
		for i := 0; i < config.ThreadCount; i++ {
			l4g.Info("spider #%d is running", i)
			go s.crawl(chFinished)
		}
		for done := 0; done < config.ThreadCount; {
			// 通知主goroutine任务结束
			<-chFinished
			l4g.Info("finiched one !")
			done++
		}
		break
	}
}
