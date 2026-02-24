---
title: "哈希表相关"
date: 2023-02-26T20:19:59+08:00
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

Redis 服务器中与哈希表相关的命令。

<!--more-->

## HSET：设置哈希表的值

`HSET` 用于设置散列一个或多个的值，返回设置成功的 field-value 对数。

格式：`HSET key field value [field value ...]`

```bash
127.0.0.1:6379> HSET person name a2htray height 179
(integer) 2
127.0.0.1:6379> HSET person weight 73.5
(integer) 1
```

## HGET：取哈希表中相应 field 对应的值

`HGET` 用于获取哈希表中单个 field 所对应的值，若 key 不存在或 field 不存在，返回 nil。

格式：`HGET key field`

```bash
127.0.0.1:6379> HGET person name
"a2htray"
127.0.0.1:6379> HGET otherPerson name
(nil)
127.0.0.1:6379> HGET person age
(nil)
```

## HINCRBY：对哈希表中整数值进行自增

`HINCRBY` 作用于哈希表可解析为整数的值，并为其加上对应整数的增量，若需要减少，将整数设置为负数即可。

格式：`HINCRBY key field increment`

```bash
127.0.0.1:6379> HSET score math 100 english 100
(integer) 2
127.0.0.1:6379> HINCRBY score math -40
(integer) 60
127.0.0.1:6379> HINCRBY score english 20
(integer) 120
```

## HINCRBYFLOAT：对哈希表中浮点数进行自增

与 `HINCRBY` 类似，不过 `HINCRBYFLOAT` 作用于可解析成浮点数的值。

格式：`HINCRBYFLOAT key field increment`

```bash
127.0.0.1:6379> HGET person weight
"73.5"
127.0.0.1:6379> HINCRBYFLOAT person weight 0.5
"74"
```

## HSTRLEN：获取哈希表中某个字段值的长度

`HSTRLEN` 用于获取哈希表中某个字段值的长度。

格式：`HSTRLEN key field`

```bash
127.0.0.1:6379> HSTRLEN person name
(integer) 7
```

## HDEL：删除一个或多个哈希表中的字段

`HDEL` 可用于邮件一个或多个哈希表中的字段，返回删除成功的字段个数。与 `DEL` 指令类似，但 `DEL` 用于删除 Redis 中存储的键。

格式：`HDEL key field [field ...]`

```bash
127.0.0.1:6379> HKEYS person
1) "name"
2) "height"
3) "weight"
127.0.0.1:6379> HDEL person height weight
(integer) 2
127.0.0.1:6379> HKEYS person
1) "name"
```

## HKEYS：列出哈希表中的字段列表

`HKEYS` 可用于列表哈希表中字段的列表，若指定的 key 不存在，返回空数组。

格式：`HKEYS key`

```bash
127.0.0.1:6379> HKEYS person
1) "name"
127.0.0.1:6379> HKEYS notExistPerson
(empty array)
```

## HEXISTS：判断哈希表中是否存在某个字段

`HEXISTS` 用于检查哈希表中是否存在某个特定的字段。

格式：`HEXISTS key field`

```bash
127.0.0.1:6379> HEXISTS person name
(integer) 1
127.0.0.1:6379> HEXISTS person notExistName
(integer) 0
```

## HGETALL：获取哈希表中所有的字段和值

`HGETALL` 用于获取哈希表中所有字段和值。

格式：`HGETALL key`

```bash
127.0.0.1:6379> HGETALL person
1) "name"
2) "a2htray"
```

## HLEN：取哈希表中字段的个数

`HLEN` 可用于获取哈希表中字段的个数。

格式：`HLEN key`

```bash
127.0.0.1:6379> HKEYS person
1) "name"
2) "weight"
3) "height"
127.0.0.1:6379> HLEN person
(integer) 3
```

## HVALS：取哈希表中所有值

`HVALS` 可用于获取哈希表中所有的值。

格式：`HVALS key`

```bash
127.0.0.1:6379> HVALS person
1) "a2htray"
2) "68.06"
3) "178.5"
```

## HMSET：同时设置哈希表中的多个字段和值

`HMSET` 可同时设置哈希表中的多个字段和值。

格式：`HMSET key field value [field value ...]`

```bash
127.0.0.1:6379> HMSET house size 78 price 150 location Wuhan
OK
```

## HMGET：获取哈希表中多个字段对应的值

`HMGET` 可同时获取哈希表中的多个字段对应的值。

格式：`HMGET key field [field ...]`

```bash
127.0.0.1:6379> HMGET house size price
1) "78"
2) "150"
```

## HSETNX：设置哈希表的字段和值，字段不存在时有效

`HSETNX` 只有在字段不存在时，设置才会有效。

格式：`HSETNX key field value`

```bash
127.0.0.1:6379> HSETNX house size 90
(integer) 0
127.0.0.1:6379> HSETNX house floor 16
(integer) 1
```

