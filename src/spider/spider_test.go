/* spider_test.go - spider测试用例 */
/*
modification history
--------------------
2016-01-14, by wusuopubupt, create
*/
/*
DESCRIPTION
测试用例
*/
package spider

import (
	"testing"
)

import (
    conf "github.com/wusuopubupt/go_spider/src/conf"
)


// NewSpider()方法单元测试
func TestNewSpider(t *testing.T) {
    confPath := "../../conf"
    confFile := confPath + "/spider.conf"
    seedUrls := []string{"http://www.baidu.com"}
    conf, _ := conf.InitConf(confFile)
    if s := NewSpider(seedUrls, conf, confPath); s == nil {
        t.Error("spider.NewSpider failes")
    } else {
        t.Log("spider.NewSpider passed.")
    }
}
