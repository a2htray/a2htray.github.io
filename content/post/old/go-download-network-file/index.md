---
title: "URL 下载网络文件"
date: 2022-04-24T10:00:25+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - HTTP
 - I/O
---

从网络上下载文件是开发过程中常用的需求，常规流程是：（1）发送请求；（2）接收响应并读取响应体内容；（3）保存到本地文件。本文包含的两个例子分别来自于**参考 [1]** 和**参考 [2]**，在此基础上做了少量的修改。

<!--more-->

## 例 1 普通下载

```go
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	fileUrl := "https://golangcode.com/logo.svg"
	filename := "logo.svg"
	resp, err := http.Get(fileUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("Downloaded: " + fileUrl)
}
```

例 1 中使用了 `io.Copy` 方法将响应体内容复制到目标文件。`io.Copy` 是带缓冲的复制，可以避免在内存中堆积大量的数据，类似的方法还有 `io.CopyBuffer`。

```go
func Copy(dst Writer, src Reader) (written int64, err error) {
	return copyBuffer(dst, src, nil)
}
func CopyBuffer(dst Writer, src Reader, buf []byte) (written int64, err error) {
	if buf != nil && len(buf) == 0 {
		panic("empty buffer in CopyBuffer")
	}
	return copyBuffer(dst, src, buf)
}
```

## 例 2 带进度的下载

```go
package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"io"
	"net/http"
	"os"
	"strings"
)

type Progress struct {
	total uint64
}

func (p *Progress) Write(bs []byte) (n int, err error) {
	n = len(bs)
	p.total += uint64(n)
	p.Show()
	return n, nil
}

func (p Progress) Show() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))
	fmt.Printf("\rDownloading... %s complete", humanize.Bytes(p.total))
}

func main() {
	fileUrl := "https://golangcode.com/logo.svg"
	filename := "logo.svg"
	resp, err := http.Get(fileUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	out, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	if _, err = io.Copy(out, io.TeeReader(resp.Body, &Progress{})); err != nil {
		panic(err)
	}
}
```

与例 1 的不同之处在于 `io.Copy` 的第 2 参数换成了一个 `io.teeReader` 的对象。下面是 `io.teeReader` 的定义：

```go
type teeReader struct {
	r Reader
	w Writer
}

func (t *teeReader) Read(p []byte) (n int, err error) {
	n, err = t.r.Read(p)
	if n > 0 {
		if n, err := t.w.Write(p[:n]); err != nil {
			return n, err
		}
	}
	return
}
```

`teeReader` 实现了 `Reader` 接口，而在 `Read` 方法中保留了原先读取数据的操作，新增了一个写数据的操作。`Progress` 实现了 `Writer` 接口，正好可以作为 `teeReader` 的 `w` 字段的值，所以在执行 `Read` 的过程中会调用 `Progress.Write` 方法，从而可以知道已经读取数据的大小。最后用 `Progress.Show` 方法将 `total` 字段的值输出到终端。

## 参考

1. [GolangCode: Download a File (from a URL)](https://golangcode.com/download-a-file-from-a-url/)
2. [GolangCode: Download Large Files with Progress Reports](https://golangcode.com/download-a-file-with-progress/)

