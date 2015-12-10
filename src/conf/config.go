package conf

import (
	"fmt"
)

import (
	"code.google.com/p/gcfg"
	l4g "code.google.com/p/log4go"
)

// spider config
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

// initialize config
func InitConf(confFile string) (*SpiderCfg, error) {
	l4g.Info("config file: %s", confFile)
	var cfg SpiderCfg

	// read conf
	err := gcfg.ReadFileInto(&cfg, confFile)
	if err != nil {
		l4g.Error("read config err [%s]", err)
		return nil, err
	}

	// check conf
	if err = checkConf(&cfg); err != nil {
		l4g.Error("read config err [%s]", err)
		return nil, err
	}

	return &cfg, nil
}

// check conf
func checkConf(s *SpiderCfg) error {
	c := s.Spider
	if c.UrlListFile == nil {
		return fmt.Errorf("Spider conf item: UrlListFile is not configured")
	}
	if c.OutputDirectory == nil {
		return fmt.Errorf("Spider conf item: OutputDirectory is not configured")
	}
	if *c.MaxDepth < 0 {
		return fmt.Errorf("Spider conf item: MaxDepth should be greater than 0")
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
