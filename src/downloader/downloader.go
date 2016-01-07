/* downloader.go - 爬虫下载器 */
/*
modification history
--------------------
2015-12-20, by wusuopubupt, create
*/
/*
DESCRIPTION
从spider.go中分离出来
*/
package downloader

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

import (
	l4g "code.google.com/p/log4go"
)

/*
* SaveAsFile - 保存网页内容
*
* PARAMS:
*   - targetUrl : 目标网址
*   - outputDir : 存储路径
*
* RETURNS:
*   - true, if succeed
*   - false, if failed
*
 */
func SaveAsFile(targetUrl string, outputDir string) bool {
	res, err := http.Get(targetUrl)
	defer res.Body.Close()
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		l4g.Error("read url content%s, err:%s", targetUrl, err)
		return false
	}
	filename := path.Join(outputDir, url.QueryEscape(targetUrl))
	f, err := os.Create(filename)
	if err != nil {
		l4g.Error("create file %s, err:%s", filename, err)
		return false
	}
	defer f.Close()
	f.Write(content)
	return true
}
