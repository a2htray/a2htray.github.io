---
title: "服务相关"
date: 2022-03-30T09:24:49+08:00
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
---

Redis 服务器中与服务相关的命令。

<!--more-->

## INFO：查看当前服务器信息

格式：`INFO [section]`

```bash
127.0.0.1:6380> INFO
# Server
redis_version:6.2.6
redis_git_sha1:00000000
redis_git_dirty:0
redis_build_id:9c9e426e2f96cc51
redis_mode:standalone
os:Linux 5.4.0-77-generic x86_64
arch_bits:64
# 还有很多 ...
```

## SHUTDOWN：客户端断开连接

格式：`SHUTDOWN [NOSAVE|SAVE]`

```bash
127.0.0.1:6379> SHUTDOWN SAVE
not connected>
```

## EXIT：退出客户端

```bash
127.0.0.1:6379> EXIT
```

