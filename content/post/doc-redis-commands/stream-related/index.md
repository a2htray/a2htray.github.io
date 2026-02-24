---
title: "流相关"
date: 2022-03-27T05:47:22+08:00
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
- images/redis-stream.png
---

Redis 服务器中与 stream 相关的命令。

<!--more-->

## XADD：向 stream 添加消息

`XADD` 可以向 stream 添加消息，返回实体 entry 的 ID。

格式：`XADD key *|ID field value [field value ...]`

```bash
127.0.0.1:6379> XADD chat:1:messages * msg "hello world" date 2020
"1648184286632-0"
127.0.0.1:6379> XADD chat:2:messages 1 msg "hello world" date 2020
"1-0"
```

## XRANGE：返回 stream 记录的列表

`XRANGE` 用于获取指定 ID 范围内的 entry，其中 `-` 和 `+` 为特征 ID，分别表示最小 ID 和最大 ID。

格式：`XRANGE key start stop [COUNT count]`

```bash
127.0.0.1:6379> XADD chat:1:messages * msg "hello GO" date 2021
"1648184498673-0"
127.0.0.1:6379> XADD chat:1:messages * msg "hello TypeScript" date 2022
"1648184539697-0"
127.0.0.1:6379> XRANGE chat:1:messages - +
1) 1) "1648184286632-0"
   2) 1) "msg"
      2) "hello world"
      3) "date"
      4) "2020"
2) 1) "1648184498673-0"
   2) 1) "msg"
      2) "hello GO"
      3) "date"
      4) "2021"
3) 1) "1648184539697-0"
   2) 1) "msg"
      2) "hello TypeScript"
      3) "date"
      4) "2022"
127.0.0.1:6379> XRANGE chat:1:messages 1648184286632-0 1648184286632-1
1) 1) "1648184286632-0"
   2) 1) "msg"
      2) "hello world"
      3) "date"
      4) "2020"
127.0.0.1:6379> XRANGE chat:1:messages 1648184286632-0 1648184498673-1
1) 1) "1648184286632-0"
   2) 1) "msg"
      2) "hello world"
      3) "date"
      4) "2020"
2) 1) "1648184498673-0"
   2) 1) "msg"
      2) "hello GO"
      3) "date"
      4) "2021"
```

## XREVRANGE

`XREVRANGE` 与 `XRANGE` 用途相近，但该命令会以倒序的方式返回 entry。

格式：`XREVRANGE key end start [COUNT count]`

```bash
127.0.0.1:6379> XREVRANGE chat:1:messages 1648184498673-1 1648184286632-0
1) 1) "1648184498673-0"
   2) 1) "msg"
      2) "hello GO"
      3) "date"
      4) "2021"
2) 1) "1648184286632-0"
   2) 1) "msg"
      2) "hello world"
      3) "date"
      4) "2020"
127.0.0.1:6379> XREVRANGE chat:1:messages 1648184498673-1 1648184286632-0 COUNT 1
1) 1) "1648184498673-0"
   2) 1) "msg"
      2) "hello GO"
      3) "date"
      4) "2021"
```

## XTRIM：裁剪 stream

格式：`XTRIM key MAXLEN|MINID [=|~] threshold [LIMIT count]`

`MAXLEN`：用于保留最近的 entry

`MINID`：用于裁剪低于某一 ID 的 entry

```bash
127.0.0.1:6379> XTRIM chat:1:messages MAXLEN 2
(integer) 1
127.0.0.1:6379> XRANGE chat:1:messages - +
1) 1) "1648184498673-0"
   2) 1) "msg"
      2) "hello GO"
      3) "date"
      4) "2021"
2) 1) "1648184539697-0"
   2) 1) "msg"
      2) "hello TypeScript"
      3) "date"
      4) "2022"
# 此时最早加入的 entry 已经被裁剪，stream 中保留两条 entry
127.0.0.1:6379> XTRIM chat:1:messages MINID 1648184498674
(integer) 1
# ID 低于 1648184498674 的 entry 会被删除
```

## XDEL：删除 stream 中 entry

格式：`XDEL key ID [ID ...]`

```bash
127.0.0.1:6379> XRANGE chat:1:messages - +
1) 1) "1648184539697-0"
   2) 1) "msg"
      2) "hello TypeScript"
      3) "date"
      4) "2022"
127.0.0.1:6379> XDEL chat:1:messages 1648184539697-0
(integer) 1
```

## XLEN：返回 stream 中 entry 的数目

格式：`XLEN key`

```bash
127.0.0.1:6379> XLEN chat:1:messages
(integer) 0
127.0.0.1:6379> DEL chat:1:messages
(integer) 1
```

## XREAD：从一个或多个 stream 中读取数据

格式：`XREAD [COUNT count] [BLOCK milliseconds] STREAMS key [key ...] id [id ...]`

[XREAD](https://redis.io/commands/xread/) 可以以阻塞或非阻塞的方式读取 stream 的数据（指定 BLOCK）。在获取 stream 记录时，需要指定记录的 ID。

```bash
127.0.0.1:6379> XRANGE api-request-log - +
1) 1) "1650702336219-0"
   2)  1) "remote_addr"
       2) "[::1]:54058"
       3) "url"
       4) "/api/users"
       5) "access_time"
       6) "1650702336"
       7) "time_executed"
       8) "0"
       9) "body_bytes_sent"
      10) "96"
2) 1) "1650702505299-0"
   2)  1) "remote_addr"
       2) "[::1]:54112"
       3) "url"
       4) "/api/users"
       5) "access_time"
       6) "1650702505"
       7) "time_executed"
       8) "0"
       9) "body_bytes_sent"
      10) "96"
127.0.0.1:6379> XREAD COUNT 1 BLOCK 1000 STREAMS api-request-log 1650702336219-0
1) 1) "api-request-log"
   2) 1) 1) "1650702505299-0"
         2)  1) "remote_addr"
             2) "[::1]:54112"
             3) "url"
             4) "/api/users"
             5) "access_time"
             6) "1650702505"
             7) "time_executed"
             8) "0"
             9) "body_bytes_sent"
            10) "96"
```

