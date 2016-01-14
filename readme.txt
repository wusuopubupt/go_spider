# go_spider by Dash Wang, 2015-11-25

目录说明(cd $HOME/work/src/github.com/wusuopubupt/go_spider && tree)：

├── readme.txt
├── conf(配置文件目录)
│   ├── example.xml
│   ├── logconf.xml
│   └── spider.conf
├── data(数据目录)
│   └── url.data
├── log(日志目录)
│   ├── mini_spider.log
│   └── mini_spider.wf.log
├── output(输出文件目录)
└── src(核心代码目录)
    ├── conf
    │   └── config.go
    ├── main
    │   └── main.go
    ├── downloader
    |   └── downloader.go
    ├── spider
    │   └── spider.go
    └── utils
        └── utils.go

运行：
cd $HOME/work/src/github.com/wusuopubupt/go_spider/src/main && go run main.go
