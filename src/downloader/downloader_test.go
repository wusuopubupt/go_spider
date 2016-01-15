/* downloader_test.go - downloader测试用例 */
/*
modification history
--------------------
2016-01-15, by wusuopubupt, create
*/
/*
DESCRIPTION
*/
package downloader

import (
	"testing"
)

// SaveAsFile()方法测试用例
func TestSaveAsFile(t *testing.T) {
    targetUrl := "http://pycm.baidu.com:8081/"
    outputDir := "../../output/"
    if ret := SaveAsFile(targetUrl, outputDir); ret != true {
        t.Error("failed!")
    } else {
        t.Log("passed!") 
    }
}
