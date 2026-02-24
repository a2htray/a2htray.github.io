---
title: "同源策略 Same-Origin Policy"
date: 2022-03-29T15:56:02+08:00
draft: false
reward: false
tags:
 - Web 安全
categories:
 - 前端技术
 - 基础原理
---

同源策略（Same-Origin Policy，SOP）是一种保护 Web 资源的安全机制，它限制了不同源之间的资源访问。需要说明的是，SOP 只作用于应用脚本，这意味着在 HTML 标签中可以引入不同源的图片、CSS 文件或动态加载的脚本文件（见[验证 1](#验证-1)）。

<!--more-->

## 同源 URL

统一资源标识符（Uniform Resource Locator，URL）标识了一个 Web 资源，其格式为：

 `schema://host[:port][/path ...]`

其中 `schema` 可为 http 或 https，port 默认为 80。如果两个 URL 的 schema、host、port 都相同时，则认为这两个 URL 是同源的。现有 URL 为 `http://foo.com/bar`，以下是其它 URL 是否同源的说明。

| URL                   | 是否同源 | 说明        |
| --------------------- | -------- | ----------- |
| https://foo.com       | 否       | schema 不同 |
| http://bar.com        | 否       | host 不同   |
| http://foo.com:81/bar | 否       | port 不同   |
| http://foo.com/zot    | 是       | 3 个都相同  |

### 访问规则

通常，直接读取跨域资源是不允许的，但仍然可以通过内嵌跨域资源进行访问。以下是允许跨域访问的规则：

| 方式       | 说明                                                         |
| ---------- | ------------------------------------------------------------ |
| iframes    | 响应头的 [`X-Frame-Options`](https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Headers/X-Frame-Options) 字段可以设置 `<frame>`、`<iframe>`、`<embed>` 或 `object` 标签可引用的页面，但跨域读 iframe 里的内容是不允许的 |
| CSS        | `<link>` 标签的 `href` 属性和 CSS 文件中的 `@import` 指令    |
| forms      | 此处不应该是读取，而是说 `<form>` 的 `action` 属性可以设置不同源的 URL，指的是目标服务可以接收不同源的数据 |
| images     | 通过 `<img>` 标签访问跨域图片，但在 `canvas` 元素里加载跨域图片是不允许的 |
| multimedia | 通过 `<video>` 和 `<audio>` 标签加载跨域的多媒体资源         |
| script     | 通过 `<script> ` 标签加载跨域的脚本，但请求跨域的 API 是不允许的 |

所以以上的规则容易变成 Web 服务攻击的入口，应当警惕。

## 验证

### 验证 1

示例目录结构

```bash
验证 1
  - static
    - images
      profile.jpg
  index.html
  main.go
```

```go
// main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("./static/"))
	http.Handle("/static/", http.StripPrefix("/static", fs))
	log.Fatal(http.ListenAndServe(":8001", nil))
}
```

```html
<!-- index.html -->
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>验证 1</title>
</head>
<body>
    <p>图片不受同源策略限制</p>
    <img width="100" height="100" alt="profile" src="http://localhost:8001/static/images/profile.jpg">
</body>
</html>
```

![Center](QQ截图20220329171230.png)

## 总结

1. URL 是否同源取决于 `schema`、`host`、`port` 是否一致；
2. 尽管跨域访问是不允许的，但仍然有一定的跨域可访问规则；