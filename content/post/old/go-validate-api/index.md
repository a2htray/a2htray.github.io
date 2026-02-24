---
title: "API 参数验证 - foo"
date: 2022-05-11T11:11:30+08:00
draft: true
reward: false
categories:
  - 后端技术
  - Go
tags:
 - parameter validation
---

一个 API 请求代表了一次用户的操作，当请求到达服务器后台程序时，程序应对请求的各项参数进行验证。各项参数包括：

1. 查询参数（Query Parameter，QP）；
2. 路径参数（Path Parameter，PP）；
3. cookie 参数（Cookie Parameter，CP）；
4. HTTP 请求头（HTTP Request Header，HRH）；

参数验证的好处在于：

1. 过滤无效请求。当程序察觉到 API 请求的参数明显不合法或不支持，则可以直接返回并防止对数据库的查询；