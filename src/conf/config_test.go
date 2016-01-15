/* config_test.go - config测试用例*/
/*
modification history
--------------------
2016-01-15, by wusuopubupt, create
*/
/*
DESCRIPTION
*/
package conf

import (
	"testing"
)

// InitConf()方法测试用例
func TestInitConf(t *testing.T){
    confFile := "../../conf/spider.conf"
    _,err := InitConf(confFile);if err != nil {
        t.Error("config.InitConf failes")
    } else {
        t.Log("config.InitConf passed.")
    }
}
