package main

import (
	"flag"
	"fmt"
	utils "github.com/wusuopubupt/go_spider/utils"
)

func main() {
	// refer : http://www.01happy.com/golang-command-line-arguments/
	// 方法一： flag.StringVar(),传入指针，直接给confPath赋值
	var confPath string
	var printVer bool
	flag.StringVar(&confPath, "c", "../conf", "config file path")
	flag.BoolVar(&printVer, "v", false, "print version")
	// 方法二：flag.String()，把函数调用的返回值赋值给logPath
	//logPath := flag.String("l", "../log", "log file path")

	flag.Parse()

	fmt.Println("Hi, dash's go_mini_spider is running...")
	if printVer {
		utils.PrintVersion()
	}
}
