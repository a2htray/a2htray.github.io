---
title: "Go 1.16 运行 Revel 项目"
date: 2022-03-31T10:40:58+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - Revel
---

[Revel](http://revel.github.io/) 是一个以高效率、高性能著称的 Go Web 框架，提供了路由、参数解析和验证、会话机制、模板机制、缓存和任务管理等诸多常用的 Web 开发功能。同时作为一个全栈的 MVC 框架， Revel 通过模块实现了组件的复用，因此可以大大提高开发者的效率。其高性能则是依托 Go 语言的性能，相信这个不必多说。但相较于其它职责相对单一的 Web 框架（如 Gin、go-restful），Revel 只能说是在保证性能的基础上尽可能地对开发者友好。

<!--more-->

## 问题重现

**环境**

1. Go 的版本：go1.16.9 windows/amd64

2. Revel：v1.0.0



今天在试验 Revel 项目时，运行新建的项目会报错，如下：

```bash
$ revel run -a .
ERROR 10:46:39   file.go:372: Error                                     seeking=github.com/revel/revel count=1 App Import Path=github.com/revel/revel filesystem path=github.com/revel/revel errors="[-: no required module provides package github.com/revel/revel; to add it:\n\tgo get github.com/revel/revel]"
Downloading related packages ... completed.
Revel executing: run a Revel application
WARN  10:46:41 harness.go:175: No http.addr specified in the app.conf listening on localhost interface only. This will not allow external access to your application
Changed detected, recompiling
Parsing packages, (may require download if not cached)...Changed detected, recompiling
 Completed
ERROR 10:46:44  build.go:406: Build errors                             errors="C:\\Users\\a2htray\\go\\pkg\\mod\\github.com\\revel\\revel@v1.0.0\\cache\\memcached.go:11:2: no required module provides package github.com/bradfitz/gomemcache/memcache; to add it:\n\tgo get github.com/bradfitz/gomemcache/memcache\nC:\\Users\\a2htray\\go\\pkg\\mod\\github.com\\revel\\revel@v1.0.0\\cache\\redis.go:10:2: no required module provides package github.com/garyburd/redigo/redis; to add it:\n\tgo get github.com/garyburd/redigo/redis\nC:\\Users\\a2htray\\go\\pkg\\mod\\github.com\\revel\\revel@v1.0.0\\cache\\inmemory.go:12:2: no required module provides package github.com/patrickmn/go-cache; to add it:\n\tgo get github.com/patrickmn/go-cache\n"
C:\Users\a2htray\go\src\omics-framework\C:\Users\a2htray\go\pkg\mod\github.com\revel\revel@v1.0.0\cache\memcached.go:11
WARN  10:46:44  build.go:420: Could not find in GO path                file=C:\\Users\\a2htray\\go\\pkg\\mod\\github.com\\revel\\revel@v1.0.0\\cache\\memcached.go:11
ERROR 10:46:44 harness.go:239: Build detected an error                  error="Go Compilation Error (in C:\\Users\\a2htray\\go\\pkg\\mod\\github.com\\revel\\revel@v1.0.0\\cache\\memcached.go:11:2): no required module provides package github.com/bradfitz/gomemcache/memcache; to add it:"

Error compiling code, to view error details see proxy running on http://:9000


Time to recompile 2.4257731s
```

新建的 Revel 项目使用 `go.mod` 对包进行管理，初始的包如下：

```bash
require (
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/revel/modules v1.0.0
	github.com/revel/revel v1.0.0
)
```

通过错误输出，可知当前项目缺少了 3 个包：

```bash
github.com/bradfitz/gomemcache/memcache
github.com/garyburd/redigo/redis
github.com/patrickmn/go-cache
```

## 解决办法

既然是项目缺少包，那就只要把缺失的包 `go get` 一下即可。

```bash
$ go get github.com/bradfitz/gomemcache/memcache
$ go get github.com/garyburd/redigo/redis
$ go get github.com/patrickmn/go-cache
```

再次运行：

```bash
$ revel run -a .
Revel executing: run a Revel application
WARN  11:33:47 harness.go:175: No http.addr specified in the app.conf listening on localhost interface only. This will not allow external access to your application
Changed detected, recompiling
Parsing packages, (may require download if not cached)... Completed
Changed detected, recompiling
INFO  11:33:54    app     run.go:34: Running revel server
INFO  11:33:54    app   plugin.go:9: Go to /@tests to run the tests.
Revel proxy is listening, point your browser to : Revel engine is listening on.. localhost:52469
9000

Time to recompile 7.0696399s
```

## 其它

在这个 [`issue 1528 `](https://github.com/revel/revel/issues/1528) 里，有人说是 Go 版本问题，即在 Go 1.15 中是可以运行的。我想解决上面的问题，就把缺失包补上就可以了，而且也猜应该不是 Go 版本问题，毕竟 Revel 的 `cache` 实现中也确实使用了这 3 个包。再深入想一想，如果 Revel 项目也是通过 Go Module 管理包的话，`revel run` 的时候就会自动下载这些包。