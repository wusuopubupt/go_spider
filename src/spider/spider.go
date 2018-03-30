/* spider.go - 爬虫主程序 */
/*
modification history
--------------------
2015-11-25, by wusuopubupt, create
2016-01-11, by wusuopubupt, 修改同步方式为sync.waitGroup
2016-01-14, by wusuopubupt, 下载方法改为异步
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
	"sync"
	"time"
)

import (
	l4g "github.com/alecthomas/log4go"
	"golang.org/x/net/html"
)

import (
	conf "github.com/wusuopubupt/go_spider/src/conf"
	downloader "github.com/wusuopubupt/go_spider/src/downloader"
)

// 任务队列最大长度
var MAX_JOBS = 1000000

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
	threadCount   int
	targetUrl     *regexp.Regexp
	jobs          chan Job
	visitedUrl    map[string]bool
	wg            sync.WaitGroup
	stop          chan bool
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
func (s *Spider) parseHtml(b io.Reader, job Job) []string {
	var urls = make([]string, 0)
	z := html.NewTokenizer(b)
	for {
		tokenType := z.Next()
		switch {
		case tokenType == html.ErrorToken:
			return urls
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
				link = u.Scheme + "://" + u.Host + "/"+ link
			}
			// 检查url是否为需要存储的目标网页url格式
			if s.checkUrlRegexp(link) {
				// 保存为文件
				go s.save(link)
			}
			if !s.visitedUrl[link] && job.depth < s.maxDepth {
				urls = append(urls, link)
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
*   - id
 */
func (s *Spider) crawl(id int) {
CRAWL:
	for {
		select {
		case <-s.stop:
			l4g.Info("spider stop...")
			break CRAWL
		case job := <-s.jobs:
	        l4g.Info("goroutine#%d parsing url:%s, depth:%d",id, job.url, job.depth)
			s.work(job)
			// 抓取间隔控制
			time.Sleep(time.Duration(s.crawlInterval) * time.Second)
		}
	}
}

/*
* work - 处理单个任务
*
* PARAMS:
*   - job : 一个待处理任务
*
 */
func (s *Spider) work(job Job) {
	defer s.wg.Done()
	// 检查是否访问过
	if s.visitedUrl[job.url] {
		l4g.Info("visted job,continue. url:%s, depth:%d", job.url, job.depth)
		return
	}
	// 判断是否超出最大爬取深度
	if job.depth > s.maxDepth {
		l4g.Info("visted job,continue. url:%s, depth:%d", job.url, job.depth)
		return
	}
	// 标记为访问过
	s.visitedUrl[job.url] = true
	resp, err := http.Get(job.url)
	if err != nil {
		l4g.Error("Failed to crawl %s, err[%s]", job.url, err)
		return
	} else {
		//l4g.Info("http response:%s", resp)
	}
	defer resp.Body.Close()
	// 解析Html, 获取新的url并入任务队列
	urls := s.parseHtml(resp.Body, job)
	for _, url := range urls {
		//新任务入公共队列
		s.wg.Add(1)
		s.addJob(Job{url, job.depth + 1})
		l4g.Info("add job: %s, depth:%d", url, job.depth+1)
	}
}

/*
* NewSpider - 初始化爬虫
*
* PARAMS:
*   - seedUrls : 种子url切片
*   - config   : 爬虫配置文件
*   - confpath : 配置文件路径
*
* RETURNS:
*	*Spider 爬虫对象
 */
func NewSpider(seedUrls []string, config conf.SpiderStruct, confpath string) *Spider {
	s := new(Spider)
	// 队列最多为100w个任务，否则阻塞
	jobs := make(chan Job, MAX_JOBS)
	// 初始化任务队列
	for _, url := range seedUrls {
		l4g.Info("url: %s", url)
		// waitgroup+1, 比自己用channel length或者timeout去判断更准确
		s.wg.Add(1)
		jobs <- Job{url, 0}
	}
	s.outputDir = path.Join(confpath, config.OutputDirectory)
	s.maxDepth = config.MaxDepth
	s.crawlInterval = config.CrawlInterval
	s.crawlTimeout = config.CrawlTimeout
	s.targetUrl = regexp.MustCompile(config.TargetUrl)
	s.jobs = jobs
	s.visitedUrl = make(map[string]bool)
	s.threadCount = config.ThreadCount

	return s
}

/*
* Start - 启动爬虫
*
*
* RETURNS:
 */
func (s *Spider) Start() {
	for i := 0; i < s.threadCount; i++ {
		l4g.Info("spider #%d is running", i)
		go s.crawl(i)
	}
}

/*
* Wait - 等待所有goroutine完成任务
*
* PARAMS:
*
* RETURNS:
 */
func (s *Spider) Wait() {
	s.wg.Wait()
	l4g.Info("finiched !, channel len: %d!", len(s.jobs))
}

/*
* Stop - 停止爬虫
*
* PARAMS:
*
* RETURNS:
 */
func (s *Spider) Stop() {
	time.Sleep(time.Duration(s.crawlInterval) * time.Second)
	for i := 0; i < s.threadCount; i++ {
		s.stop <- true
	}
}
