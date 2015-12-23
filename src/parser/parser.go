/**
 * @author : wusuopubupt
 * @date   : 2015-11-15
 * @brief  : 网页解析器
 */
package parser

import (
	"net/url"
	"strings"
)

import (
	l4g "code.google.com/p/log4go"
)

// 提取url,注意相对路径
func (s *Spider) getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

// 解析html
// 以后针对不同的爬取任务，设定不同的parse方法
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
			if !hasProto {
				link = url + "/" + link
			}
			u, _ := url.Parse(link)
			realUrl := u.Scheme + "://" + u.Host + u.Path
			if !s.visitedUrl[realUrl] && hasProto && job.depth < s.maxDepth {
				// 新任务入公共队列
				s.addJob(Job{realUrl, job.depth + 1})
				l4g.Info("add job: %s, depth:%d", realUrl, job.depth+1)
			}
		}
	}
}

// 爬取和解析(getJob & addJob)
func (s *Spider) crawl(chFinished chan bool) {
	// 通知主goroutine，当前goroutine已无任务可做
	defer func() {
		chFinished <- true
	}()
	// 等待队列中任务到达的超时时间，3秒
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Second * 3)
		timeout <- true
	}()
CRAWL:
	for {
		var job Job
		select {
		case <-timeout:
			l4g.Info("get job timeout!")
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
			/////////////////////////////////////
			///  网络请求和解析单独设计package实现
			/////////////////////////////////////
			s.visitedUrl[job.url] = true
			resp, err := http.Get(job.url)
			if err != nil {
				l4g.Error("Failed to crawl %s, err[%s]", job.url, err)
				return
			}
			defer resp.Body.Close()
			/////////////////////////////////////
			// 以后针对不同的爬取任务，设定不同的parse方法
			s.parseHtml(resp.Body, job)
			/////////////////////////////////////
			/////////////////////////////////////
			// 抓取间隔控制
			time.Sleep(time.Duration(s.crawlInterval) * time.Second)
		}
	}
}

// 初始化爬虫
func newSpider(config conf.SpiderStruct, jobs chan Job) *Spider {
	s := new(Spider)
	s.outputDir = config.OutputDirectory
	s.maxDepth = config.MaxDepth
	s.crawlInterval = config.CrawlInterval
	s.crawlTimeout = config.CrawlTimeout
	s.targetUrl = config.TargetUrl
	s.jobs = jobs
	s.visitedUrl = make(map[string]bool)

	return s
}

// 启动爬虫
func Start(seedUrls []string, config conf.SpiderStruct) {
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
		s := newSpider(config, jobs)
		// 开启threandCount个spider.crawl goroutine,等待通道中的任务到达
		for i := 0; i < config.ThreadCount; i++ {
			l4g.Info("spider #%d is running", i)
			go s.crawl(chFinished)
		}
		for done := 0; done < config.ThreadCount; {
			// 阻塞,等待通知主goroutine任务结束
			<-chFinished
			l4g.Info("finiched #%d!", done)
			done++
		}
		break WORKING
	}
}
