---
title: "集合相关"
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
  - images/redis-set.png
---

Redis 服务器中与集合相关的命令。

<!--more-->

## SADD：向集合中添加元素，若元素存在则忽略本次操作

`SADD` 返回成功添加到集合的元素个数，集合中已存在的元素不作添加处理。

格式：`SADD key member [member ...]`

```bash
127.0.0.1:6379> SADD odds 1 3 5 7
(integer) 4
127.0.0.1:6379> SADD odds 7 9
(integer) 1
# 此时，只有元素 9 被成功添加
```

## SREM：删除集合中的元素

`SREM` 返回成功删除的元素个数。

格式：`SREM key member [member ...]`

```bash
127.0.0.1:6379> SREM odds 9 7
(integer) 2
127.0.0.1:6379> SMEMBERS odds
1) "1"
2) "3"
3) "5"
```

## SMEMBERS：获取集合中的所有元素

格式：`SMEMBERS key`

```bash
127.0.0.1:6379> SMEMBERS odds
1) "1"
2) "3"
3) "5"
```

## SCARD：获取集合中元素的个数

格式：`SCARD key`

```bash
127.0.0.1:6379> SCARD odds
(integer) 3
```

## SRANDMEMBER：随机获取集合中指定个数的元素

格式：`SRNADMEMBER key [count]`

```bash
127.0.0.1:6379> SMEMBERS odds
1) "1"
2) "3"
3) "5"
127.0.0.1:6379> SRANDMEMBER odds
"1"
127.0.0.1:6379> SRANDMEMBER odds 2
1) "1"
2) "5"
```

## SPOP：从集合中随机弹出一个或多个元素

格式：`SPOP key [count]`

```bash
127.0.0.1:6379> SPOP odds 1
1) "3"
127.0.0.1:6379> SMEMBERS odds
1) "1"
2) "5"
127.0.0.1:6379> SPOP odds 2
1) "1"
2) "5"
127.0.0.1:6379> SMEMBERS odds
(empty array)
```

## SDIFF：集合运算，取差集

`SDIFF` 可以对多个集合进行取差集操作。

格式：`SDIFF key [key ...]`

```bash
127.0.0.1:6379> SADD integers 1 2 3 4 5 6 7 8 9 10
(integer) 10
127.0.0.1:6379> SADD odds 1 3 5 7 9
(integer) 5
127.0.0.1:6379> SDIFF integers odds
1) "2"
2) "4"
3) "6"
4) "8"
5) "10"
```

## SINTER：集合运算，取交集

`SINTER` 可以对多个集合进行取交集操作。

格式：`SINTER key [key ...]`

```bash
127.0.0.1:6379> SINTER integers odds
1) "1"
2) "3"
3) "5"
4) "7"
5) "9"
```

## SUNION：集合运算，取并集

`SUNION` 可以对多个集合进行取并集操作。

格式：`SUNION key [key ...]`

```bash
127.0.0.1:6379> SADD odds 1 3 5 7 9
(integer) 0
127.0.0.1:6379> SADD evens 2 4 6 8 10
(integer) 5
127.0.0.1:6379> SUNION odds evens
 1) "1"
 2) "2"
 3) "3"
 4) "4"
 5) "5"
 6) "6"
 7) "7"
 8) "8"
 9) "9"
10) "10"
```

## SDIFFSTORE：集合运算，存储差集

`SDIFFSTORE` 可以将多个集合的取差集的运算结果保存到另一个键中，返回值为差集的元素个数。

格式：`SDIFFSTORE destination key [key ...]`

```bash
127.0.0.1:6379> SDIFFSTORE new:odds integers evens
(integer) 5
```

## SINTERSTORE：集合运算，存储交集

`SINTERSTORE` 可以将多个集合的取交集的运算结果保存到另一个键中，返回值为交集的元素个数。

格式：`SINTERSTORE destination key [key ...]`

```bash
127.0.0.1:6379> SINTERSTORE new:evens integers evens
(integer) 5
```

## SUNIONSTORE：集合运算，存储并集

`SUNIONSTORE` 可以将多个集合的取并集的运算结果保存到另一个键中，返回值为并集的元素个数。

格式：`SUNIONSTORE destination key [key ...]`

```bash
127.0.0.1:6379> SUNIONSTORE new:integers odds evens
(integer) 10
127.0.0.1:6379> SMEMBERS new:integers
 1) "1"
 2) "2"
 3) "3"
 4) "4"
 5) "5"
 6) "6"
 7) "7"
 8) "8"
 9) "9"
10) "10"
```
