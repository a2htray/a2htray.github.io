---
title: "Redis with Python（三） 集合操作及运算"
date: 2023-02-26T16:59:30+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Python
tags:
  - Redis
  - Python
---

set 数据类型对应元素不重复的数据结构，在 Redis 中，set 数据类型是无序的，与之相对的有序集合 zset。本文内容分两部分：

* 集合操作
* 集合运算

<!--more-->

## 集合操作

```python
import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)

if __name__ == '__main__':
    # SADD 向集合中添加元素
    # 返回添加成功元素的个数
    print(rs.sadd('numbers', 1, 2, 3, 4, 5))  # 5
    # 由于 4, 5 已存在于 numbers 集合，所以只需要添加 6, 7, 8，返回成功添加的个数 3
    print(rs.sadd('numbers', 4, 5, 6, 7, 8))  # 3

    # SREM
    # 删除集合中的元素，返回成功删除的元素个数
    print(rs.srem('numbers', *[1, 2, 3]))  # 3
    print(rs.srem('numbers', *[1, 2, 3, 4]))  # 1

    # SMEMBERS
    # 获取集合中的所有元素，返回类型对应 Python 中的 set 类型
    print(rs.sadd('alphabet', *['a', 'b', 'c']))  # 3
    print(rs.smembers('alphabet'))  # {b'b', b'a', b'c'}

    # SCARD
    # 取集合中元素的个数
    print(rs.scard('alphabet'))  # 3

    # SRANDMEMBER
    # 随机取集合中指定个数的元素
    # 随机输出，结果不保证一致
    print(rs.srandmember('alphabet', 1))  # [b'b']
    # 指定个数大于集合元素个数时，返回全部元素
    print(rs.srandmember('alphabet', 4))  # [b'c', b'a', b'b']

    # SPOP
    # 从集合中弹出指定个数的元素
    print(rs.spop('alphabet', 2))  # [b'b', b'c']
    print(rs.smembers('alphabet'))  # {b'a'}
```

代码在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/set_basic.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/set_basic.py)。

## 集合运算

```python
import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)

if __name__ == '__main__':
    rs.sadd('projectA', *['A001', 'A002', 'A003'])
    rs.sadd('projectAB', *['A001', 'A002', 'A003', 'B001', 'B002'])

    # SDIFF
    # 取集合的差集
    # 若参数为 S1, S2, S3, ...，则返回结果为 S1 - S2 - S3 - ...
    print(rs.sdiff('projectA', 'projectAB'))  # set()
    print(rs.sdiff('projectAB', 'projectA'))  # {b'B001', b'B002'}

    # SINTER
    # 取集合的交集
    print(rs.sinter('projectA', 'projectAB'))  # {b'A001', b'A003', b'A002'}

    # SUNION
    # 取集合的并集
    rs.sadd('projectB', *['B001', 'B002'])
    print(rs.sunion('projectA', 'projectB'))  # {b'A002', b'A001', b'B001', b'A003', b'B002'}

    # SDIFFSTORE
    # 计算集合的差集并将结果保存至另一键的集合中
    # 返回差集中元素的个数
    rs.delete('projectB')
    print(rs.sdiffstore('projectB', 'projectAB', 'projectA'))  # 2
    print(rs.smembers('projectB'))  # {b'B001', b'B002'}

    # SINTERSTORE
    # 计算集合的交集并将结果保存至另一键的集合中
    # 返回交集中元素的个数
    print(rs.sinterstore('newProjectB', 'projectAB', 'projectB'))  # 2
    print(rs.smembers('newProjectB'))  # {b'B001', b'B002'}

    # SUNIONSTORE
    # 计算集合的并集并将结果保存至另一键的集合中
    # 返回并集中元素的个数
    rs.delete('projectAB')
    print(rs.sunionstore('projectAB', 'projectA', 'projectB'))  # 5
    print(rs.smembers('projectAB'))  # {b'A002', b'A001', b'B001', b'A003', b'B002'}
```

代码在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/set_ops.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/set_ops.py)。

## 小结

Redis 中的 set 数据类型保存的是无序的元素集合，Python 调用方式与 Redis 指令高度一致。