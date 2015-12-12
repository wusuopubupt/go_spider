package conf

import (
	"fmt"
)

import (
	"code.google.com/p/gcfg"
	l4g "code.google.com/p/log4go"
)

// spider config,和../../conf/spider.conf一一对应
type SpiderStruct struct {
	UrlListFile     string
	OutputDirectory string
	MaxDepth        int
	CrawlInterval   int
	CrawlTimeout    int
	TargetUrl       string
	ThreadCount     int
}

// 嵌套结构体分开声明，用起来更灵活
type SpiderCfg struct {
	Spider SpiderStruct
}

// initialize config
func InitConf(confFile string) (SpiderStruct, error) {
	l4g.Info("config file: %s", confFile)
	var cfg SpiderCfg

	// read conf
	err := gcfg.ReadFileInto(&cfg, confFile)
	if err != nil {
		l4g.Error("read config err [%s]", err)
		return cfg.Spider, nil
	}

	// check conf
	if err = checkConf(&cfg); err != nil {
		l4g.Error("read config err [%s]", err)
		return cfg.Spider, nil
	}

	return cfg.Spider, nil
}

// check conf
func checkConf(cfg *SpiderCfg) error {
	s := cfg.Spider
	if s.UrlListFile == "" {
		return fmt.Errorf("Spider conf item: UrlListFile is not configured")
	}
	if s.OutputDirectory == "" {
		return fmt.Errorf("Spider conf item: OutputDirectory is not configured")
	}
	if s.MaxDepth < 0 {
		return fmt.Errorf("Spider conf item: MaxDepth should be greater than 0")
	}
	if s.CrawlInterval == 0 {
		return fmt.Errorf("Spider conf item: CrawlInterval is not configured")
	}
	if s.CrawlTimeout == 0 {
		return fmt.Errorf("Spider conf item: CrawlTimeout is not configured")
	}
	if s.TargetUrl == "" {
		return fmt.Errorf("Spider conf item: TargetUrl is not configured")
	}
	if s.ThreadCount == 0 {
		return fmt.Errorf("Spider conf item: ThreadCount is not configured")
	}
	return nil
}
