---
title: "列表值相关"
date: 2022-03-27T05:32:39+08:00
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
  - images/redis-list.png
---

Redis 服务器中与列表值相关的命令。

<!--more-->

## LPUSH：向列表左端添加元素

`LPUSH` 返回添加元素后列表的长度。

格式：`LPUSH key element [element ...]`

```bash
127.0.0.1:6379> LPUSH colors green yellow red blue gray
(integer) 5
# 此时列表为 [gray blue yellow red green]
```

## LPUSHX：向列表左端添加元素

`LPUSHX` 与 `LPUSH` 类似，但只有在 key 存在的情况，操作才有效。

格式：`LPUSHX key element [element ...]`

```bash
127.0.0.1:6379> KEYS colors
(empty array)
127.0.0.1:6379> LPUSHX colors red green blue
(integer) 0
```

## RPUSH：向列表右端添加元素

`RPUSH` 返回添加元素后列表的长度。

格式：`RPUSH key element [element ...]`

```bash
127.0.0.1:6379> RPUSH colors lightgreen lightyellow lightred lightblue
(integer) 9
# 此时列表为 [gray blue yellow red green
# lightgreen lightyellow lightred lightblue]
```

## RPUSHX：向列表右端添加元素

`RPUSHX` 与 `RPUSH` 类似，但只能存在的键有效。

格式：`RPUSHX key element [element ...]`

```bash
127.0.0.1:6379> KEYS colors
(empty array)
127.0.0.1:6379> RPUSHX colors red green blue
(integer) 0
127.0.0.1:6379> KEYS colors
(empty array)
```

## LPOP：从列表左端弹出元素

`LPOP` 返回弹出的元素。

格式：`LPOP key [count]`

```bash
127.0.0.1:6379> LPOP colors 2
1) "gray"
2) "blue"
# 此时列表为 [yellow red green
# lightgreen lightyellow lightred lightblue]
```

## RPOP：从列表右端弹出元素

`RPOP` 返回弹出的元素。

格式：`RPOP key [count]`

```bash
127.0.0.1:6379> RPOP colors 2
1) "lightblue"
2) "lightred"
# 此时列表为 [yellow red green
# lightgreen lightyellow]
```



## LLEN：获取列表中元素的个数

格式：`LLEN key`

```bash
127.0.0.1:6379> LLEN colors
(integer) 5
```

## LRANGE：获取列表指定区间上的元素

`LRANGE` 指定的区间包括两端。

格式：`LRANGE key start stop`

```bash
127.0.0.1:6379> LRANGE colors 2 -1
1) "green"
2) "lightgreen"
3) "lightyellow"
```

## LREM：删除列表中前 count 个指定的元素

格式：`LREM key count element`

```bash
127.0.0.1:6379> LPUSH colors yellow yellow yellow
(integer) 8
# 此时列表为 [yellow yellow yellow yellow red green
# lightgreen lightyellow]
127.0.0.1:6379> LREM colors 3 yellow
(integer) 3
# 此时列表为 [yellow red green
# lightgreen lightyellow]
```

## LINDEX：获取指定位置上的元素

格式：`LINDEX key index`

```bash
127.0.0.1:6379> LINDEX colors 2
"green"
```

## LSET：设置列表中指定位置上元素的值

格式：`LSET key index element`

```bash
127.0.0.1:6379> LSET colors 0 blue
OK
# 此时列表为 [blue red green
# lightgreen lightyellow]
```

## LTRIM：对列表进行裁剪

`LTRIM` 裁剪列表并保存到原有列表中。

格式：`LTRIM key start stop`

```bash
127.0.0.1:6379> LTRIM colors 0 2
OK
# 列表只保留了前 3 个元素
127.0.0.1:6379> LRANGE colors 0 9
1) "blue"
2) "yellow"
3) "green"
```

## LINSERT：向列表插入元素

`LINSERT` 用于在列表元素前或后插入指定元素。

格式：`LINSERT key BEFORE|AFTER pivot element`

```bash
127.0.0.1:6379> LINSERT colors BEFORE blue red
(integer) 4
127.0.0.1:6379> LRANGE colors 0 -1
1) "red"
2) "blue"
3) "yellow"
4) "green"
```

## RPOPLPUSH：操作两个列表，对元素进行弹出再推入

`RPOPLPUSH` 返回值为第 1 个列表弹出的元素。

格式：`RPOPLPUSH source destination`

```bash
127.0.0.1:6379> RPOPLPUSH colors other:colors
"green"
127.0.0.1:6379> LRANGE colors 0 -1
1) "red"
2) "blue"
3) "yellow"
127.0.0.1:6379> LRANGE other:colors 0 -1
1) "green"
```

## BLPOP：阻塞式从列表左端弹出元素

`BLPOP` 同样用于从列表左端弹出元素，但是当列表为空，该命令会阻塞列表直到超时或列表有元素可弹出，超时时间单位为秒。

格式：`BLPOP key [key ...] timeout`

```bash
127.0.0.1:6379> BLPOP mock:list 2
(nil)
(2.06s)
```

## BRPOP：阻塞式从列表右端弹出元素

`BRPOP` 同样用于从列表右端弹出元素，但是当列表为空，该命令会阻塞列表直到超时或列表有元素可弹出，超时时间单位为秒。

格式：`BRPOP key [key ...] timeout`

```bash
127.0.0.1:6379> BRPOP mock:list 2`
(nil)
(2.05s)
```

## BRPOPLPUSH：操作两个列表，对元素进行弹出再推入

`BRPOPLPUSH` 与 `RPOPLPUSH` 类似，但如果列表中没有元素会阻塞直到等待超时或有元素弹出，超时时间的单位为秒。

格式：`BRPOPLPUSH source destination timeout`

```bash
127.0.0.1:6379> RPUSH colors red green yellow
(integer) 3
127.0.0.1:6379> BRPOPLPUSH colors dest:colors 10
"yellow"
127.0.0.1:6379> LRANGE colors 0 -1
1) "red"
2) "green"
127.0.0.1:6379> LRANGE dest:colors 0 -1
1) "yellow"
```

