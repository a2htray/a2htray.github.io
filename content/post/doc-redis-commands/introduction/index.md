---
title: "介绍"
date: 2022-03-27T05:56:48+08:00
draft: false
reward: false
toc: true
pinned: false
categories:
  - 后端技术
  - Redis
tags:
  - "Redis"
series:
  - "Redis 命令手册"
type: "docs"
images:
  - images/redis.png
---

Redis 是开源的、高性能的数据结构存储系统，在框架设计中常常被当作缓存服务器。不同于传统的关系型数据库（如 MySQL、PostgreSQL），Redis 将数据以键值对的方式存储于内存并且支持数据持久化。尽管 Redis 采用了单线程模型来处理请求，但其通过 I/O 多路复用技术做到了应用级别的异步，运行的性能也十分良好。

<!--more-->

根据操作对象的不同，可将 Redis 中的命令分成以下几类：

1. 键相关命令
2. 字符串值相关命令
3. 列表值相关命令
4. 集合值相关命令
5. 有序集合值相关命令
6. 流类型值相关命令
7. 集群相关命令
8. 服务相关命令
