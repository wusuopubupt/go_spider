/* config.go - 爬虫配置文件解析程序 */
/*
modification history
--------------------
2015-11-25, by wusuopubupt, create
*/
/*
DESCRIPTION
爬虫配置文件解析
*/
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

/*
* InitConf - 初始化配置
*
* PARAMS:
*   - confFile : 配置文件名(全路径)
*
* RETURNS:
*   - (SpiderStruct ,nil), if succeed
*   - (err,SpiderStruct) if failed
 */
func InitConf(confFile string) (SpiderStruct, error) {
	l4g.Info("config file: %s", confFile)
	var cfg SpiderCfg

	// read conf
	err := gcfg.ReadFileInto(&cfg, confFile)
	if err != nil {
		l4g.Error("read config err [%s]", err)
		return cfg.Spider, err
	}

	// check conf
	if err = checkConf(&cfg); err != nil {
		l4g.Error("read config err [%s]", err)
		return cfg.Spider, err
	}

	return cfg.Spider, nil
}

/*
* checkConf - 检查配置文件合法性
*
* PARAMS:
*   - cfg: SpiderCfg结构指针
*
* RETURNS:
*   - nil, if succeed
*   - error, if failed
 */
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
