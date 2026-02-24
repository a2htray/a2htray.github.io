---
title: "go module 使用本地开发的包"
date: 2022-04-25T11:04:33+08:00
draft: false
toc: true
reward: false
categories:
 - 后端技术
 - Go
tags:
 - stackoverflow
 - go module
---

在 `go.mod` 文件中新增 `replace` 信息，内容如下：

```bash
module github.com/userName/mainModule

require "github.com/userName/otherModule" v0.0.0
replace "github.com/userName/otherModule" v0.0.0 => "本地包路径"
```

<!--more-->

[**Accessing local packages within a go module (go 1.11)**](https://stackoverflow.com/questions/52026284/accessing-local-packages-within-a-go-module-go-1-11)

