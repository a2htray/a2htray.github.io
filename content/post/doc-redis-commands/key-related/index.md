---
title: "键相关"
date: 2022-03-27T05:57:48+08:00
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
  - images/redis-key.jpeg
---

Redis 服务器中与 key 相关的命令。

<!--more-->

## KEYS：获取数据库中匹配规则的键名

`KEYS` 命令遍历数据库中的所有键，支持 glob 风格通配符格式，在存在大量键值对的 Redis 服务器上应谨慎使用。

格式：`KEYS patten`

```bash
127.0.0.1:6379> KEYS *
(empty array)
127.0.0.1:6379> SET config:logLevel Fatal
OK
127.0.0.1:6379> KEYS config:*
1) "config:logLevel"
```

> glob 风格通配符格式
>
> | 符号   | 含义         |
> | ---- | ---------- |
> | ?    | 匹配一个字符     |
> | *    | 匹配任意多个字符   |
> | []   | 匹配括号间的任一字符 |
> | \    | 转义         |

## EXISTS：判断键名是否存在

`EXISTS` 用于判断键名是否存在，返回值为存在键名的个数。

格式：`EXISTS key [key ...]`

```bash
127.0.0.1:6379> EXISTS config:logLevel config:pagination
(integer) 1
127.0.0.1:6379> EXISTS config:pagination
(integer) 0
```

## EXPIRE：给键设置过期时间

`EXPIRE` 可以给一个键设置一个以秒为单位的过期时间。

格式：`EXPIRE key seconds`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> EXPIRE user:1:name 5
(integer) 1
127.0.0.1:6379> GET user:1:name
"xiaoming"
# 5 秒内访问
127.0.0.1:6379> GET user:1:name
(nil)
# 5 秒后访问
```

## EXPIREAT：给键设置过期时间

`EXPIREAT` 通过指定一个 UNIX 时间戳为键设置一个过期时间。

格式：`EXPIREAT key timestamp`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> EXPIREAT user:1:name 1648470000
(integer) 1
# 1648470000 为 2022-03-28 20:20:00
127.0.0.1:6379> GET user:1:name
"xiaoming"
```

## PEXPIRE：给键设置过期时间

`PEXPIRE` 与 `EXPIRE` 类似，不同之处在于 `PEXPIRE` 的时间单位是微秒。

格式：`PEXPIRE key milliseconds`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> PEXPIRE user:1:name 10000
(integer) 1
127.0.0.1:6379> GET user:1:name
"xiaoming"
# 10 秒内访问
127.0.0.1:6379> GET user:1:name
(nil)
# 10 秒后访问
```

## PEXPIREAT：给键设置过期时间

`PEXPIREAT` 与 `EXPIREAT` 类似，不同之外在于 `PEXPIREAT` 的时间单位是微秒。

格式：`PEXPIREAT key milliseconds-timestamp`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> TTL user:1:name
(integer) -1
127.0.0.1:6379> PEXPIREAT user:1:name 1648470000000
(integer) 1
```

## PERSIST：移除键的过期时间

`PERSIST` 可以移除键的过期，使其永不失效。

格式：`PERSIST key`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> EXPIRE user:1:name 100
(integer) 1
127.0.0.1:6379> TTL user:1:name
(integer) 93
127.0.0.1:6379> PERSIST user:1:name
(integer) 1
127.0.0.1:6379> TTL user:1:name
(integer) -1
```

## TTL：返回键的剩余生存时间

`TTL` 返回以秒为单位的键的剩余生存时间，对长期有效的键使用会返回 -1。

格式：`TTL key`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> EXPIRE user:1:name 100
(integer) 1
127.0.0.1:6379> TTL user:1:name
(integer) 93
127.0.0.1:6379> SET user:1:name xiaoming
(integer) 1
127.0.0.1:6379> TTL user:1:name
(integer) -1
```

## PTTL：返回键的剩余生存时间

`PTTL` 与 `TTL` 类似，不同之外在于返回的剩余生存时间的单位为微秒。

格式：`PTTL key`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> PTTL user:1:name
(integer) -1
127.0.0.1:6379> EXPIRE user:1:name 100
(integer) 1
127.0.0.1:6379> PTTL user:1:name
(integer) 93978
```

## RENAME：修改键名

`RENAME` 可用于修改键名。

格式：`RENAME key newkey`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> RENAME user:1:name user:2:name
OK
127.0.0.1:6379> GET user:2:name
"xiaoming"
```

## RENAMENX：修改键名

`RENAMENX` 命令只有在给定的新键名不存在时，才会起作用。

格式：`RENAMENX key newkey`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> KEYS user:newname*
(empty array)
127.0.0.1:6379> RENAMENX user:1:name user:newname:1:name
(integer) 1
127.0.0.1:6379> GET user:1:name
(nil)
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> RENAMENX user:1:name user:1:name
(integer) 0
# 新键名不变，但执行不成功
```

## DEL：删除一个或多个键

`DEL` 用于删除一个或多个键，返回值为删除键的个数。

格式：`DEL key [key ...]`

```bash
127.0.0.1:6379> DEL config:logLevel config:pagination
(integer) 1
# config:pagination 此时并不存在，故返回值为 1
```

## RANDOMKEY：随机返回一个键

格式：`RANDOMKEY key`

```bash
127.0.0.1:6379> RANDOMKEY
"student:weights"
```

## DUMP：序列化给定的键

`DUMP` 可以序列化指定的键并返回序列化的值。

格式：`DUMP key`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> DUMP user:1:name
"\x00\bxiaoming\t\x00\xe6u\x97\x84\x19\x1c\x01\x81"
```

## TYPE：获取指定键对应值的类型

`TYPE` 用于获取指定键对应值的类型，返回值包括 `string | hash | list | set | zset | stream`

格式：`TYPE key`

```bash
127.0.0.1:6379> SET user:1:name xiaoming
OK
127.0.0.1:6379> TYPE user:1:name
string
```

## DBSIZE：返回数据库中 key 的数量

格式：`DBSIZE`

```bash
127.0.0.1:6379> DBSIZE
(integer) 26
```
