---
title: "Redis 的数据持久化"
date: 2022-03-27T09:17:03+08:00
draft: true
comment: false
featured: true
reward: false
toc: true
pinned: false
categories:
  - "数据库"
tags:
  - "Redis"
weight: 1
---

数据的持久化是指将数据以某种方式存储起来，待服务重启后，数据可以正常恢复访问。由于 Redis 的数据存储在内存中，一旦宕机或关机，数据将不复存在。所以 Redis 提供两种数据持久化方式将内存中的数据保存到磁盘文件中，分别为 RDB（Redis Database）和 AOF（Append Only File）。

<!--more-->

## RDB



## AOF



## 混合模式



## 总结

