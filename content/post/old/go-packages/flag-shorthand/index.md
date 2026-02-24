---
title: "Go flag 支持选项简写"
date: 2023-03-09T20:17:17+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - Go
 - flags
---

在命令行工具的开发过程中，我们常常需要设置一个同时支持短名称和长名称的选项，如 `-d` 等价于 `--debug`、`-p` 等价于 `--password`。在使用 Go flag 包的情况下，该需求的实现相当简单，只需要定义两个不同的 flag 指向同一个变量即可。

<!--more-->

## 代码

```go
// flag_shorthand.go
package main

import (
	"flag"
	"fmt"
)

func main() {
	var debug bool
	var password string
	debugUsage := "run in debug mode"
	passwordUsage := "type password"

	flag.BoolVar(&debug, "d", false, debugUsage)
	flag.BoolVar(&debug, "debug", false, debugUsage)
	flag.StringVar(&password, "p", "", passwordUsage)
	flag.StringVar(&password, "password", "", passwordUsage)
	flag.Parse()

	fmt.Printf("debug: %v\n", debug)
	fmt.Printf("password: %s\n", password)
}
```

## 运行

**不带任何参数**

`debug` 和 `password` 都使用默认值，分别为 `false` 和 `""`。

```bash
$ go run flag_shorthand.go
debug: false
password: 
```

**使用短名称**

```bash
$ go run flag_shorthand.go -d -p 12345
debug: true
password: 12345
```

**使用长名称**

```bash
$ go run flag_shorthand.go --debug --password=12345
debug: true
password: 12345
```

## 小结

学习并掌握了长短名称选项的使用方式，[代码](https://github.com/a2htray/code-notebook/blob/main/BlogCode/Go/flag_shorthand/flag_shorthand.go)。