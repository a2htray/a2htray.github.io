---
title: "有序集合相关"
date: 2022-03-27T05:45:52+08:00
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
- images/redis-zset.png
---

Redis 服务器中与有序集合相关的命令。

<!--more-->

## ZADD：添加元素

`ZADD` 用于将一个或多个带分数的元素添加到有序集合中，返回成功添加到有序集合的元素个数。

当添加元素已在有序集合中时，更新元素的分数使其在有序集合中保持正确的位置。

格式：`ZADD key score member [score member ...]`

```bash
127.0.0.1:6379> ZADD student:weights 63.2 xiaoming 67.5 xiaolei
(integer) 2
127.0.0.1:6379> ZADD student:weights 64.2 xiaoming
(integer) 0
```

## ZSCORE：获取元素的分数

格式：`ZSCORE key member`

```bash
127.0.0.1:6379> ZSCORE student:weights xiaoming
"64.200000000000003"
127.0.0.1:6379> ZSCORE student:weights xiaolei
"67.5"
```

## ZRANGE：获取指定位置区间上的元素

`ZRANGE` 可以获取指定位置区间上的元素，包括区间的两端。

格式：`ZRANGE key min max`

```bash
127.0.0.1:6379> ZADD student:weights 81.5 xiaopang
(integer) 1
127.0.0.1:6379> ZRANGE student:weights 1 2
1) "xiaolei"
2) "xiaopang"
127.0.0.1:6379> ZRANGE student:weights 2 2
1) "xiaopang"
```

## ZRANGEBYSCORE：获取指定分数区间上的元素

`ZRANGEBYSCORE` 可指定分数区间获取元素。

`+inf` 表示正无穷，`-inf` 表示负无穷。

格式：`ZRANGEBYSCORE key min max [WITHSCORES] [LIMIT offset count]`

```bash
127.0.0.1:6379> ZRANGEBYSCORE student:weights 65 85 WITHSCORES
1) "xiaolei"
2) "67.5"
3) "xiaopang"
4) "81.5"
127.0.0.1:6379> ZRANGEBYSCORE student:weights 65 85 WITHSCORES LIMIT 1 1
1) "xiaopang"
2) "81.5"
127.0.0.1:6379> ZRANGEBYSCORE student:weights 65 85 WITHSCORES LIMIT 1 2
1) "xiaopang"
2) "81.5"
127.0.0.1:6379> ZRANGEBYSCORE student:weights 65 85 WITHSCORES LIMIT 0 2
1) "xiaolei"
2) "67.5"
3) "xiaopang"
4) "81.5"
127.0.0.1:6379> ZRANGEBYSCORE student:weights -inf +inf WITHSCORES
1) "xiaoming"
2) "64.200000000000003"
3) "xiaolei"
4) "67.5"
5) "xiaopang"
6) "81.5"
```

## ZINCRBY：增加某个元素的分数

`ZINCRBY` 的返回值为修改后的分数。

格式：`ZINCRBY key increment member`

```bash
127.0.0.1:6379> ZINCRBY student:weights 0.5 xiaoming
"64.700000000000003"
127.0.0.1:6379> ZINCRBY student:weights -0.5 xiaoming
"64.200000000000003"
```

## ZCARD：获取集合中元素的个数

格式：`ZCARD key`

```bash
127.0.0.1:6379> ZCARD student:weights
(integer) 3
```

## ZCOUNT：获取指定分数范围内的元素个数

格式：`ZCOUNT key min max`

```bash
127.0.0.1:6379> ZCOUNT student:weights 65 85
(integer) 2
```

## ZREM：删除一个或多个元素

`ZREM` 返回删除成功的元素个数。

格式：`ZREM key member [member ...]`

```bash
127.0.0.1:6379> ZREM student:weights xiaoming xiaolei
(integer) 2
127.0.0.1:6379> ZREM student:weights notaname
(integer) 0
```

## ZREMRANGEBYRANK：通过指定位置区间删除集合元素

格式：`ZREMRANGEBYRANK key start stop`

```bash
127.0.0.1:6379> ZADD student:weights 63.2 xiaoming 67.5 xiaolei
(integer) 2
127.0.0.1:6379> ZRANGEBYSCORE student:weights -inf +inf WITHSCORES
1) "xiaoming"
2) "63.200000000000003"
3) "xiaolei"
4) "67.5"
5) "xiaopang"
6) "81.5"
127.0.0.1:6379> ZREMRANGEBYRANK student:weights 0 1
(integer) 2
127.0.0.1:6379> ZRANGEBYSCORE student:weights -inf +inf WITHSCORES
1) "xiaopang"
2) "81.5"
```

## ZREMRANGEBYSCORE：通过指定分数区间删除集合元素

格式：`ZREMRANGEBYSCORE key min max`

```bash
127.0.0.1:6379> ZADD student:weights 63.2 xiaoming 67.5 xiaolei 81.5 xiaopang
(integer) 3
127.0.0.1:6379> ZRANGEBYSCORE student:weights -inf +inf WITHSCORES
1) "xiaoming"
2) "63.200000000000003"
3) "xiaolei"
4) "67.5"
5) "xiaopang"
6) "81.5"
127.0.0.1:6379> ZREMRANGEBYSCORE student:weights 80 85
(integer) 1
127.0.0.1:6379> ZRANGEBYSCORE student:weights -inf +inf WITHSCORES
1) "xiaoming"
2) "63.200000000000003"
3) "xiaolei"
4) "67.5"
```

## ZRANK：获取元素的排序

格式：`ZRANK key member`

```bash
127.0.0.1:6379> ZRANK student:weights xiaolei
(integer) 1
127.0.0.1:6379> ZRANK student:weights notaname
(nil)
```

## ZREVRANK：降序获取元素的排序

格式：`ZREVRANK key member`

```bash
127.0.0.1:6379> ZREVRANK student:weights xiaolei
(integer) 0
```