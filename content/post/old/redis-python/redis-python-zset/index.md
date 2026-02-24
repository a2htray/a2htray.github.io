---
title: "Redis with Python（四） 有序集合操作"
date: 2023-02-26T20:10:40+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Python
tags:
  - Redis
  - Python
---

zset 有序集合是 set 的补充，zset 中的元素都带有一个用于排序的分数，以下通过代码示例进行说明。

<!--more-->

## 常规操作

```python
import random

import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)

if __name__ == '__main__':
    # ZADD
    # 设置一个或多个带分数的值
    # 返回设置成功的个数
    # 若值已存在，则更新其分数
    print(rs.zadd('heights', {'King': 146, 'Jimmy': 160, 'Ann': 155}))  # 3
    # 更新分数不计入设置成功的个数
    print(rs.zadd('heights', {'Jimmy': 165}))  # 0

    # ZSCORE
    # 获取值的分数
    print(rs.zscore('heights', 'King'))  # 146.0

    # ZRANGE
    # 取指定索引区间上的值
    print(rs.zrange('heights', 1, 2))  # [b'Ann', b'Jimmy']

    # ZRANGEBYSCORE
    # 取分数符合区间内的元素
    print(rs.zrangebyscore('heights', 140, 159))  # [b'King', b'Ann']
    # 第 3 个参数指定开始查找的索引
    # 第 4 个参数指定返回值的个数
    # withscores 为 True 时，返回值要带上分数
    print(rs.zrangebyscore('heights', 140, 159, 1, 1, withscores=True))  # [(b'Ann', 155.0)]

    # ZINCRBY
    # 增加某个值的分数
    # 返回更新后的分数
    print(rs.zincrby('heights', 5, 'Ann'))  # 160.0

    # ZCARD
    # 返回集合中元素的个数
    print(rs.zcard('heights'))  # 3

    # ZCOUNT
    # 返回指定分数区间上值的个数
    print(rs.zcount('heights', 155, 170))  # 2

    # ZREM
    # 删除一个或多个元素
    # 返回删除成功的个数
    print(rs.zrem('heights', 'Jimmy', 'Ann'))  # 2

    # ZREMRANGEBYRANK
    # 删除指定区间上的元素
    # 返回成功删除的个数
    rs.zadd('floats', {'0.1': 0.1, '0.2': 0.2, '0.3': 0.3})
    print(rs.zremrangebyrank('floats', 0, 1))  # 2
    print(rs.zrange('floats', 0, rs.zcard('floats')))  # [b'0.3']

    # ZREMRANGEBYSCORE
    # 删除符合指定分数区间的元素
    # 返回成功删除的个数
    rand_values = {}
    s = 'abcdefg'
    for v in range(len(s)):
        rand_values[s[v]] = random.random()
    rs.zadd('rand_values', rand_values)
    print(rs.zrange('rand_values', 0, len(s),
                    withscores=True))  # [(b'd', 0.200721567466151), (b'f', 0.3296169544640849), (b'c', 0.3677595274693486), (b'b', 0.512665164306821), (b'g', 0.5353782879780641), (b'a', 0.7277933908083912), (b'e', 0.942418290434278)]

    print(rs.zremrangebyscore('rand_values', 0.3, 0.4))  # 2
    print(rs.zcard('rand_values'))  # 5

    # ZRANK
    # 获取元素在集合中的位置
    print(rs.zrank('rand_values', 'a'))  # 3

    # ZREVRANK
    # 获取元素在降序集合中的位置
    print(rs.zrevrank('rand_values', 'a'))  # 1
```

代码在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/zset.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/zset.py)。

## 小结

在 Pyhon 中用字典作为 zset 的操作值，其中字典的 key 为相应的元素，字典的 value 为对应的分数。常规操作包括：

* 设置或删除 zset 的值
* 按索引或分数进行取值
* 取索引信息

