/* spider.go - 爬虫主程序 */
/*
modification history
--------------------
2015-11-25, by wusuopubupt, create
*/
/*
DESCRIPTION
爬虫实现，页面请求和解析
*/
package spider

import (
	"io"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strings"
	"time"
)

import (
	l4g "code.google.com/p/log4go"
	"golang.org/x/net/html"
)

import (
	conf "github.com/wusuopubupt/go_spider/src/conf"
	downloader "github.com/wusuopubupt/go_spider/src/downloader"
)

// 单个任务结构
type Job struct {
	url   string
	depth int
}

// 爬虫结构
type Spider struct {
	// 多个goroutine共享的属性+任务队列
	outputDir     string
	maxDepth      int
	crawlInterval int
	crawlTimeout  int
	targetUrl     *regexp.Regexp
	jobs          chan Job
	visitedUrl    map[string]bool
}

/*
* getJob - 从队列取出任务
*
* RETURNS:
* 	-job
 */
func (s *Spider) getJob() (job Job) {
	return <-s.jobs
}

/*
* addJob - 新任务入队列
*
* PARAMS : - job
 */
func (s *Spider) addJob(job Job) {
	s.jobs <- job
}

/*
* getHref - 提取url
*
* PARAMS:
*   - t : html token
*
* RETURNS:
*   - ok, href
 */
func (s *Spider) getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

/*
* parseHtml - 解析html,以后针对不同的爬取任务，设定不同的parse方法

* PARAMS:
*   - b 	: html response body
*   - job	: current job
*
 */
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
			hasProto := strings.Index(link, "http") == 0
			u, _ := url.Parse(job.url)
			// 相对路径
			if !hasProto {
				link = u.Scheme + "://" + u.Host + "/" + link
			}
			// 检查url是否为需要存储的目标网页url格式
			if s.checkUrlRegexp(link) {
				// 保存为文件
				s.save(link)
			}
			if !s.visitedUrl[link] && job.depth < s.maxDepth {
				// 新任务入公共队列
				s.addJob(Job{link, job.depth + 1})
				l4g.Info("add job: %s, depth:%d", link, job.depth+1)
			}
		}
	}
}

/*
* save - 保存网页内容
*
* PARAMS:
*   - targetUrl : 目标网址
*
* RETURNS:
*   - true, if succeed
*   - false, if failed
*
 */
func (s *Spider) save(targetUrl string) bool {
	return downloader.SaveAsFile(targetUrl, s.outputDir)
}

/*
* checkUrlRegexp -  检查url是否为需要存储的目标网页url格式
*
* PARAMS:
*   - url
*
* RETURNS:
*   - true, if match
*   - false, if done't match
*
 */
func (s *Spider) checkUrlRegexp(url string) bool {
	return s.targetUrl.MatchString(url)
}

/*
* crawl  - 爬取和解析(getJob & addJob)
*
* PARAMS:
*   - chFinished : 当前goroutine完成时向chFinished信道发送消息,通知主goroutine
*
 */
func (s *Spider) crawl(chFinished chan bool) {
	// 通知主goroutine，当前goroutine已无任务可做
	defer func() {
		chFinished <- true
	}()
	timeout := make(chan bool, 1)
	go func() {
		// 等待队列中任务到达的超时时间，1秒
		time.Sleep(time.Second * 10)
		timeout <- true
	}()
CRAWL:
	for {
		var job Job
		select {
		case <-timeout:
			l4g.Info("get job timeout!, channel length:%d", len(s.jobs))
			break CRAWL
		case job = <-s.jobs:
			l4g.Info("get job url:%s, depth:%d, channel length:%d", job.url, job.depth, len(s.jobs))
			// 检查是否访问过
			if s.visitedUrl[job.url] {
				l4g.Info("visted job,continue. url:%s, depth:%d", job.url, job.depth)
				continue
			}
			if job.depth > s.maxDepth {
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
}

/*
* newSpider - 初始化爬虫
*
* PARAMS:
*   - config   : 爬虫配置文件
*   - jobs     : 任务队列
*   - confpath : 配置文件路径
*
* RETURNS:
*	*Spider 爬虫对象
 */
func newSpider(config conf.SpiderStruct, jobs chan Job, confpath string) *Spider {
	s := new(Spider)
	s.outputDir = path.Join(confpath, config.OutputDirectory)
	s.maxDepth = config.MaxDepth
	s.crawlInterval = config.CrawlInterval
	s.crawlTimeout = config.CrawlTimeout
	s.targetUrl = regexp.MustCompile(config.TargetUrl)
	s.jobs = jobs
	s.visitedUrl = make(map[string]bool)

	return s
}

/*
* Start - 启动爬虫
*
* PARAMS:
*   - seedUrls : 种子url切片
*   - config   : 爬虫配置文件
*   - confpath : 配置文件路径
*
* RETURNS:
 */
func Start(seedUrls []string, config conf.SpiderStruct, confpath string) {
	// 队列最多为100w个任务，否则阻塞
	jobs := make(chan Job, 1000000)
	chFinished := make(chan bool)
	// 初始化任务队列
	for _, url := range seedUrls {
		l4g.Info("url: %s", url)
		jobs <- Job{url, 0}
	}
	// 一个while(1)的循环，直到channel通知任务结束
WORKING:
	for {
		s := newSpider(config, jobs, confpath)
		// 开启threandCount个spider.crawl goroutine,等待通道中的任务到达
		for i := 0; i < config.ThreadCount; i++ {
			l4g.Info("spider #%d is running", i)
			go s.crawl(chFinished)
		}
		// 定时查看任务队列长度
		chTicker := time.NewTicker(time.Millisecond * 500).C
		for done := 0; done < config.ThreadCount; {
			select {
			case <-chTicker:
				l4g.Info("channel length:%d", len(s.jobs))
			// 阻塞,等待通知主goroutine任务结束
			case <-chFinished:
				l4g.Info("finiched #%d!", done)
				done++
			}
		}
		break WORKING
	}
}
