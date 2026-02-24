---
title: "Redis with Python（五） 哈希表操作"
date: 2023-02-27T22:46:29+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Python
tags:
 - Redis
 - Python
---

哈希表（hash）是 Redis 中重要的数据结构，本文通过示例演示如何使用 Python 完成对哈希表的操作，各方法调用分别对应着 Redis 的一个指令：

* HSET
* HGET
* HINCRBY
* HINCRBYFLOAT
* HSTRLEN
* HVALS
* HMSET
* HMGET
* HSETNX

<!--more-->

## 常规操作

```python
# -*- coding:utf-8 -*-
# Date: 2023/2/27
# Created by: a2htray
# Description: 哈希表操作
import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)

if __name__ == '__main__':
    # HSET
    # 设置哈希表的字段与值
    # 返回设置成功的字段个数
    print(rs.hset('person', 'name', 'a2htray', {
        'height': 179,
    }))  # 2

    # HGET
    # 取哈希表中某个字段的值
    print(rs.hget('person', 'name'))  # b'a2htray'

    # HINCRBY
    # 对哈希表中某个字段值（整数）进行自增
    # 返回字段的新值
    rs.hset('score', 'math', 100, {'english': 100})
    print(rs.hincrby('score', 'math', -40))  # 60
    print(rs.hincrby('score', 'english', 20))  # 120

    # HINCRBYFLOAT
    # 对哈希表中的某个字段值（浮点数）进行自增
    # 返回字段的新值
    rs.hset('person', 'weight', 73.5)
    print(rs.hincrbyfloat('person', 'weight', 0.5))  # 74.0

    # HSTRLEN
    # 获取哈希表中某个字段值的长度
    print(rs.hstrlen('person', 'name'))  # 7

    # HDEL
    # 删除一个或多个哈希表中的字段
    # 返回删除成功的字段个数
    print(rs.hdel('person', 'height', 'weight'))  # 2
    print(rs.hgetall('person'))  # {b'name': b'a2htray'}

    # HKEYS
    # 返回哈希表中的字段列表
    print(rs.hkeys('person'))  # [b'name']
    print(rs.hkeys('notExistPerson'))  # []

    # HEXISTS
    # 判断哈希表中是否存在某个字段
    print(rs.hexists('person', 'name'))  # True
    print(rs.hexists('person', 'notExistName'))  # False

    # HGETALL
    # 取哈希表所有信息
    print(rs.hgetall('person'))  # {b'name': b'a2htray'}

    # HLEN
    # 取哈希表中字段的个数
    print(rs.hlen('person'))  # 1

    # HVALS
    # 返回哈希表中值列表
    print(rs.hvals('person'))  # [b'a2htray']

    # HMSET
    # 设置一个或多个哈希表的字段和值
    # hmset 已作废，应使用 hset
    print(rs.hmset('person', {
        'height': 179,
        'weight': 73.5,
    }))

    # HMGET
    # 获取哈希表中多个值信息
    print(rs.hmget('person', *['name', 'height', 'weight']))  # [b'a2htray', b'179', b'73.5']

    # HSETNX
    # 设置哈希表的字段和值，只有当字段不存在时才生效
    print(rs.hsetnx('person', 'location', 'WuHan'))  # 1
```

## 小结

简明易用。