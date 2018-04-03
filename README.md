# go_spider 

*Created by Dash Wang, 2015-11-25*

目录说明(`cd $GOPATH/src/github.com/wusuopubupt/go_spider && tree`)：

``` shell

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
├── readme.txt
└── src(核心代码目录)
    ├── conf
    │   ├── config.go
    │   └── config_test.go
    ├── downloader
    │   ├── downloader.go
    │   └── downloader_test.go
    ├── main
    │   ├── main
    │   └── main.go
    ├── spider
    │   ├── spider.go
    │   └── spider_test.go
    ├── utils
    │   └── utils.go
    └── vendor(第三方包目录)
        ├── github.com
        │   └── alecthomas
        │       └── log4go
        │           ├── LICENSE
        │           ├── README
        │           ├── config.go
        │           ├── filelog.go
        │           ├── log4go.go
        │           ├── pattlog.go
        │           ├── socklog.go
        │           ├── termlog.go
        │           └── wrapper.go
        ├── golang.org
        │   └── x
        │       └── net
        │           ├── LICENSE
        │           ├── PATENTS
        │           └── html
        │               ├── atom
        │               │   ├── atom.go
        │               │   ├── gen.go
        │               │   └── table.go
        │               ├── const.go
        │               ├── doc.go
        │               ├── doctype.go
        │               ├── entity.go
        │               ├── escape.go
        │               ├── foreign.go
        │               ├── node.go
        │               ├── parse.go
        │               ├── render.go
        │               └── token.go
        ├── gopkg.in
        │   ├── gcfg.v1
        │   │   ├── LICENSE
        │   │   ├── README
        │   │   ├── doc.go
        │   │   ├── errors.go
        │   │   ├── go1_0.go
        │   │   ├── go1_2.go
        │   │   ├── read.go
        │   │   ├── scanner
        │   │   │   ├── errors.go
        │   │   │   └── scanner.go
        │   │   ├── set.go
        │   │   ├── token
        │   │   │   ├── position.go
        │   │   │   ├── serialize.go
        │   │   │   └── token.go
        │   │   └── types
        │   │       ├── bool.go
        │   │       ├── doc.go
        │   │       ├── enum.go
        │   │       ├── int.go
        │   │       └── scan.go
        │   └── warnings.v0
        │       ├── LICENSE
        │       ├── README
        │       └── warnings.go
        └── vendor.json
        
```

运行：

``` shell
cd $GOPATH/src/github.com/wusuopubupt/go_spider/src/main && go run main.go
```
