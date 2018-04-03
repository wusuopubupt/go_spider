/* utils.go - 爬虫公用组件 */
/*
modification history
--------------------
2015-11-25, by wusuopubupt, create
2016-04-03, by wusuopubupt, update version to 1.0.1
*/
/*
DESCRIPTION
*/
package utils

import (
	"fmt"
)

/*
* PrintVersion - print spider's version
 */
func PrintVersion() {
	const version = "1.0.1"
	fmt.Println("go_spider version", version)
}
