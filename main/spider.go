package main

import (
	// log4go install: http://wuzhuti.cn/2411.html
	l4g "code.google.com/p/log4go"
	"flag"
	"fmt"
	utils "github.com/wusuopubupt/go_spider/utils"
	"time"
)

func main() {
	l4g.LoadConfiguration("logconf.xml")

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
	l4g.Info("Hi, dash's %s is running...", "go_mini_spider")
	l4g.Error("Unable to open file: %s", "xxx")
	// And now we're ready!
	l4g.Finest("This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("About that time, eh chaps?")
	if printVer {
		utils.PrintVersion()
	}

	// http://stackoverflow.com/questions/14252766/abnormal-behavior-of-log4go
	// adding a time.Sleep(time.Second) to the end of the code snippeet will cause the log content flush
	time.Sleep(time.Second)
}
