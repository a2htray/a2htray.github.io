---
title: "字符串相关"
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
  - images/redis-medium.png
---

Redis 服务器中与字符串值相关的命令。

<!--more-->

## SET：设置键值对

`SET` 用于设置单个值为字符串的健值对。

格式：`SET key value`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
```

## SETNX：键不存在时才设置值

`SETNX` 只有在键不存在的情况下，才可以设置值。

格式：`SETNX key value`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> SETNX user:1:name xiaohong
(integer) 0
# 此处设置不成功
127.0.0.1:6379> GET user:1:name
"xiaoming"
127.0.0.1:6379> EXISTS user:1:email
(integer) 0
# 键 user:1:email 不存在
127.0.0.1:6379> SETNX user:1:email foo@bar.com
(integer) 1
# 设置成功
```

## SETEX：设置具有生存时间的键值对

`SETEX` 可以设置一个具有生存时间的键值对，过期时间的单位为秒。

格式：`SETEX key seconds value`

```bash
127.0.0.1:6379> SETEX user:1:name 10 xiaoming
OK
127.0.0.1:6379> GET user:1:name
"xiaoming"
# 10 秒内访问
127.0.0.1:6379> GET user:1:name
(nil)
# 10 秒外访问
```

## PSETEX：设置具有过期时间的键值对

`PSETEX` 与 `SETEX` 类似，但其生存时间的单位为微秒。

格式：`PSETEX key millIseconds value`

```bash
127.0.0.1:6379> PSETEX user:1:name 10000 xiaoming
OK
127.0.0.1:6379> GET user:1:name
"xiaoming"
# 10 秒内访问
127.0.0.1:6379> GET user:1:name
(nil)
# 10 秒外访问
```


## GET：取特定键对应的值

`GET` 用于获取获取特定键对应的值。

格式：`GET key`

```bash
127.0.0.1:6379> GET user:1:name
"xiaoming"
```

## GETSET：设置值并返回旧值

`GETSET` 可以为设置一个值并返回 key 的旧值，若 key 不存在则返回 `nil`。

格式：`GETSET key value`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> GETSET user:1:name xiaohong
"xiaoming"
```

## GETRANGE：返回字符串值的字符子串

`GETRANGE` 用于返回字符串值的字符子串，包括两端。

格式：`GETRANGE key start end`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> GETRANGE user:1:name 2 3
"ao"
```

## SETRANGE：给字符串值的指定位置设置新字符

`SETRANGE` 给字符串值的指定位置设置新字符并返回新字符串的长度。

格式：`SETRANGE key offset value`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> SETRANGE user:1:name 4 hong
(integer) 8
127.0.0.1:6379> GET user:1:name
"xiaohong"
```

## MSET：同时设置多个键值对

`MSET` 用于设置多个键值对。

格式：`MSET key value [key value ...]`

```bash
127.0.0.1:6379> MSET user:1:age 18 user:1:gender male
OK
```

## MSETNX：同时设置多个键值对

`MSETNX` 用于同时设置多个键值对，当有一个键是存在时，该命令其它键的赋值也不会生效。

格式：`MSETNX key value [key value ...]`

```bash
127.0.0.1:6379> KEYS user*
(empty array)
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> MSETNX user:1:name xiaoming user:2:name xiaohong
(integer) 0
# 因为 user:1:name 存在，所以 user:2:name 没能被赋值
127.0.0.1:6379> EXISTS user:1:name user:2:name
(integer) 1
127.0.0.1:6379> EXISTS user:1:name
(integer) 1
127.0.0.1:6379> EXISTS user:2:name
(integer) 0
127.0.0.1:6379> GET user:2:name
(nil)
```

## MGET：同时取多个值

`MGET` 用于同时取多个值。

格式：`MGET key [key ...]`

```bash
127.0.0.1:6379> MGET user:1:age user:1:gender
1) "18"
2) "male"
```

## STRLEN：获取字符串的长度

`STRLEN` 可以获取指定键对应字符串值的长度。

格式：`STRLEN key`

```bash
127.0.0.1:6379> STRLEN user:1:name
(integer) 8
127.0.0.1:6379> SET welcome 你好
OK
127.0.0.1:6379> STRLEN welcome
(integer) 6
```

## APPEND：追加字符串，返回追加后字符串的长度

`APPEND` 可以向字符串值追加新的字符，返回值为追加后新字符串的长度。

格式：`APPEND key value`

```bash
127.0.0.1:6379> APPEND user:1:name Li
(integer) 10
127.0.0.1:6379> GET user:1:name
"xiaomingLi"
```

## GETBIT：获取字符串值指定位置上的二进制值

指定位置指的是字符串的二进制存储形式中从左至右的偏移量，以 0 开始计数。

格式：`GETBIT key offset`

```bash
127.0.0.1:6379> SET user:2:name abc
OK
# abc 的二进制存储形式为
# a: 01100001
# b: 01100010
# c: 01100011
127.0.0.1:6379> BITCOUNT user:2:name
(integer) 10
GETBIT user:2:name 1
(integer) 1
```

## SETBIT：设置字符串值指定位置上的二进制值

返回值为原指定位置上的二进制值

格式：`SETBIT key offset value`

```bash
127.0.0.1:6379> SETBIT user:2:name 6 1
(integer) 0
# 此时 a 已变成 c
127.0.0.1:6379> GET user:2:name
"cbc"
```

## BITCOUNT：返回字符串类型值的二进制值中 1 的个数

下标指的字符在字符串中的位置，而不是字符串二进制形式的比特位置。

格式：`BITCOUNT key [start end]`

```bash
127.0.0.1:6379> SET user:2:name abc
OK
127.0.0.1:6379> BITCOUNT user:2:name
(integer) 10
127.0.0.1:6379> BITCOUNT user:2:name 0 2
(integer) 10
127.0.0.1:6379> BITCOUNT user:2:name 1 1
(integer) 3
# 实际取的字符 b 的二进制形式中 1 的个数
```

## BITOP：对多个字符串类型值进行位运算

`BITOP` 支持的运算符包括 `AND`、`OR`、`XOR` 和 `NOT`，并将运算结果保存到目标 key 中。

格式：`BITOP operation destkey key [key ...]`

```bash
127.0.0.1:6379> SET char:1 a
OK
127.0.0.1:6379> SET char:2 z
OK
127.0.0.1:6379> BITOP AND char:1:2:and:result char:1 char:2
(integer) 1
127.0.0.1:6379> GET char:1:2:and:result
"`"
127.0.0.1:6379> BITOP XOR char:1:2:xor:result char:1 char:2
(integer) 1
127.0.0.1:6379> GET char:1:2:xor:result
"\x1b"
```

## BITPOS：获取字符串类型值的二进制中第一个是 0 或 1 的位置

格式：`BITPOS key bit [start] [end]`

```bash
127.0.0.1:6379> BITPOS char:1 0 0 7
(integer) 0
# 字符 a 的二进制形式的第 0 位为 0
127.0.0.1:6379> BITPOS char:1 1 0
(integer) 1
# 字符 a 的二进制形式的第 1 位为 1
```

## INCR：数值递增

`INCR` 返回递增后的数值，若指定的 key 不存在，则先对其初始化为 0，再进行数值递增。

格式：`INCR key`

```bash
127.0.0.1:6379> SET job:crash:count 100
OK
127.0.0.1:6379> INCR job:crash:count
(integer) 101
127.0.0.1:6379> INCR job:running:count
(integer) 1
# 此时键 job:running:count 并不存在
```

## INCRBY：以指定数值对键值进行递增

`INCRBY` 返回递增后的数值，若指定的 key 不存在，则先对其初始化为 0，再进行数值递增。支持负值操作。

格式：`INCRBY key increment`

```bash
127.0.0.1:6379> INCRBY job:crash:count 2
(integer) 103
```

## DECR：数值递减

格式：`DECR key`

```bash
127.0.0.1:6379> DECR job:crash:count
(integer) 102
```

## DECRBY：以指定数值对键值进行递减

`DECRBY` 支持负值操作。

格式：`DECRBY key decrement`

```bash
127.0.0.1:6379> DECRBY job:crash:count 2
(integer) 100
```

## INCRBYFLOAT：以指定浮点数值对键值进行递增

`INCRBYFLOAT` 可通过指定浮点数对键值进行加减，若指定 key 不存在，则先进行创建。

```bash
127.0.0.1:6379> INCRBYFLOAT job:1:weight 0.4
"0.4"
127.0.0.1:6379> INCRBYFLOAT job:1:weight 0.4
"0.4"
127.0.0.1:6379> INCRBYFLOAT job:1:weight -0.2
"0.2"
```
