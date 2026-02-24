---
title: "Redis with Python（二） 字符串操作"
date: 2023-02-26T14:47:07+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Python
tags:
  - Redis
  - Python
---

string 是 Redis 中最基础的数据类型，由于它是二进制安全的，所以可以存储图片的二进制信息。本文通过 3 个部分介绍在 Python 下如何操作 Redis 的 string 数据类型：

* 常规字符串操作
* 位运算操作
* 数值操作

<!--more-->

## 常规字符串操作

```python
import time
import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)

if __name__ == '__main__':
    # SET
    # 设置单个键值对
    # 返回是否设置成功
    print(rs.set('name', 'a2htray'))  # True

    # SETNX
    # 设置单个键值对，只有当键不存在时有效
    print(rs.setnx('name', 'a2htray2'))  # False
    print(rs.setnx('notExistName', 'a2htray'))  # True

    # SETEX
    # 设置具有生存时间的键值对，单位为秒
    print(rs.setex('course', 1, 'math'))  # True
    time.sleep(1.5)  # 设置 1.5 秒后访问
    print(rs.get('course'))  # None

    # PSETEX
    # 设置具有生存时间的键值对，单位为微秒
    # 1000 微秒 = 1 秒
    print(rs.psetex('store', 1000, '59'))  # True
    time.sleep(1.5)
    print(rs.get('store'))  # None

    # GET
    # 获取值
    print(rs.get('name'))  # b'a2htray'

    # GETSET
    # 设置值并返回旧值，若旧值不存在返回 None
    print(rs.getset('name', 'jimmy'))  # b'a2htray'
    print(rs.getset('otherName', 'Tonny'))  # None

    # GETRANGE
    # 根据给定的索引区间，以闭区间的方式返回字符串的子串
    # 若指定的 key 不存在，则返回空字符串
    print(rs.getrange('name', 0, 2))  # b'jim'
    print(rs.getrange('notExistName2', 0, 2))  # b''

    # SETRANGE
    # 根据给定的索引和值，原字符串中从该索引位置起用值进行替换
    # 返回新字符串的长度
    print(rs.setrange('name', 1, 'xxx'))  # 5
    print(rs.setrange('name', 4, 'xxx'))  # 7

    # MSET
    # 同时设置多个键值对
    print(rs.mset({'width': 10, 'height': 20}))

    # MSETNX
    # 设置多个键值对，若某个键存在，其余键值设置不会生效
    # 由于 width 已存在，所以 objName 设置末生效
    print(rs.msetnx({'width': 20, 'objName': 'AD00001'}))  # False

    # MGET
    # 同时取多个值
    print(rs.mget('width', 'height'))  # [b'10', b'20']

    # STRLEN
    # 取字符串长度
    print(rs.strlen('name'))  # 7

    # APPEND
    # 追回字符串，返回追加后字符串的长度
    print(rs.append('name', '0001'))  # 11
    # 若指定的键不存在，则会先创建一个值为空的键，并在其后进行追回
    print(rs.append('anotherName', 'Ann'))  # 3
```

代码在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/string_basic.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/string_basic.py)。

## 位运算操作

```python
import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)

if __name__ == '__main__':
    # GETBIT
    # 取字符串指定位置的 bit 表示
    title = 'Redis with Python'
    for c in title:
        print(bin(ord(c)))

    # 0b1010010
    # 0b1100101
    # 0b1100100
    # 0b1101001
    # 0b1110011
    # 0b100000
    # 0b1110111
    # 0b1101001
    # 0b1110100
    # 0b1101000
    # 0b100000
    # 0b1010000
    # 0b1111001
    # 0b1110100
    # 0b1101000
    # 0b1101111
    # 0b1101110

    rs.set('title', title)

    print(rs.getbit('title', 0))  # 0
    print(rs.getbit('title', 1))  # 1
    print(rs.getbit('title', 2))  # 0
    print(rs.getbit('title', 3))  # 1
    print(rs.getbit('title', 4))  # 0
    print(rs.getbit('title', 5))  # 0
    print(rs.getbit('title', 6))  # 1
    print(rs.getbit('title', 7))  # 0

    # SETBIT
    # 设置指定位置上的二进制值
    rs.setbit('title', 8 * 15 - 1, 1)  # Python 中 h 变成了 i
    print(rs.get('title'))  # b'Redis with Pytion'

    # BITCOUNT
    # 返回字符串的二进制表示中 1 的个数
    print(rs.bitcount('title', 0, 8 * 15))  # 64

    # BITOP
    # 对多个键执行位运算
    can_view = 0b0001
    can_add = 0b0010
    can_delete = 0b0100
    can_cancel = 0b1000

    rs.set('can_view', can_view)
    rs.set('can_add', can_add)
    rs.set('can_delete', can_delete)
    rs.set('can_cancel', can_cancel)

    rs.bitop('OR', 'permission', 'can_view', 'can_delete')  # 设置权限：可看、可删除
    permission = int(rs.get('permission'))
    print(
        permission,  # 5
        permission & can_view,  # 判断是否有"可看"的权限，输出 1
        permission & can_cancel,  # 判断是否有"可删除"的权限，输出 0
    )

    # BITPOS
    # 获取字符串表示的二进制中第一个为 0 或 1 的位置
    print(rs.bitpos('permission', 1))  # 2
```

代码在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/string_bit.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/string_bit.py)。

## 数值操作

```python
import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)

if __name__ == '__main__':
    # INCR
    # 数值的递增
    rs.set('score', 88)
    rs.incr('score', 2)  # +2
    print(rs.get('score'))  # b'90'

    # INCRBY
    rs.incrby('score', 3)  # +3
    print(rs.get('score'))  # b'93'

    # INCRBYFLOAT
    rs.incrbyfloat('score', 0.5)  # +0.5
    print(rs.get('score'))  # b'93.5'

    # DESR
    # 数值的递减
    rs.set('score', 10)
    rs.decr('score', 1)
    print(rs.get('score'))  # b'9'

    # DESRBY
    rs.decrby('score', 3)
    print(rs.get('score'))  # b'6'
```

代码在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/string_number.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/string_number.py)。

## 小结

通过代码来演示如何在 Python 环境下对 Redis 的 string 类型进行操作。