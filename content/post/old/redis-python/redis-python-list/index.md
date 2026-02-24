---
title: "Redis with Python（一） 列表操作"
date: 2023-02-25T21:58:12+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Python
tags:
  - Redis
  - Python
---

list 是 Redis 中操作比较频繁的数据结构，本文将结合 Python 对其进行相关操作并按“是否发生阻塞”分成两个部分：

* 非阻塞指令
* 阻塞指令

<!--more-->

## 非阻塞操作

Redis 的 list 操作中，非阻塞指令包括：

* LPUSH 和 LPUSHX
* RPUSH 和 RPUSHX
* LPOP 和 RPOP
* LLEN
* LRANGE
* LREM
* LINDEX
* LSET
* LTRIM
* LINSERT
* RPOPLPUSH

以下是各指令的 Python 代码示例：

```python
import redis

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)


if __name__ == '__main__':
    # LPUSH
    # 创建 colors，并向左添加两个元素
    # 返回添加后列表的长度
    print(rs.lpush('colors', 'green', 'yellow'))  # 2

    # LPUSHX
    # 向 colors 添加新的元素
    # 返回添加后列表的长度
    print(rs.lpushx('colors', 'red'))  # 3

    # LPUSHX
    # 向不存在的列表添加元素，不会创建相应列表并返回 0
    print(rs.lpushx('notExistColors', 'green'))  # 0

    # RPUSH
    # 向右添加一个元素
    # 返回添加后列表的长度
    print(rs.rpush('colors', 'gray'))  # 4

    # RPUSHX
    # 向右添加一个元素
    # 返回添加后列表的长度
    print(rs.rpushx('colors', 'orange'))  # 5

    # RPUSHX
    # 向不存在的列表添加元素，不会创建相应列表并返回 0
    print(rs.rpushx('notExistColors', 'green'))  # 0

    # LPOP
    # 取列表左边的第一个元素并从列表中移除该元素
    print(rs.lpop('colors'))  # b'red'

    # RPOP
    # 取列表右边的第一个元素并从列表中移除该元素
    color = rs.rpop('colors')
    print(color)  # b'orange'

    # LLEN
    # 获取列表中元素的个数
    # 如果指定列表不存在，返回 0
    print(rs.llen('colors'))  # 3
    print(rs.llen('notExistColors'))  # 0

    # LRANGE
    # 取列表中指定区间内连续的子集

    # 取全部
    print(rs.lrange('colors', 0, -1))  # [b'yellow', b'green', b'gray']
    # 取闭区间 [0, 1] 的元素
    print(rs.lrange('colors', 0, 1))  # [b'yellow', b'green']

    # LREM
    # 从左端删除指定个数的元素
    # 第三个参数指定要删除的值
    # 返回被删除元素的个数
    print(rs.lrem('colors', 2, 'green'))  # 1

    # LINDEX
    # 取指定位置上的元素
    print(rs.lindex('colors', 0))  # b'yellow'
    # 超过列表索引最大值返回 None
    print(rs.lindex('colors', 3))  # None

    # LSET
    # 设置指定位置的值
    # 返回是否设置成功
    print(rs.lset('colors', 0, 'red'))  # True

    # LTRIM
    # 截取原列表并保存至原列表
    # 返回是否截取成功
    print(rs.ltrim('colors', 0, 0))  # True

    # LINSERT
    # 在指定元素前或后插入元素
    # 返回插入成功后的列表长度
    print(rs.linsert('colors', 'AFTER', 'red', 'green'))  # 2
    # 向不存在的列表插入元素返回 0
    print(rs.linsert('notExistColors', 'BEFORE', 'red', 'green'))  # 0

    # RPOPLPUSH
    # 操作两个列表，从一个列表右端弹出元素并向左端添加到另一个列表
    # 返回被操作的元素
    rs.lpush('otherColors', 'yellow')
    print(rs.rpoplpush('colors', 'otherColors'))  # b'green'

    rs.delete('colors')
    rs.delete('otherColors')
```

代码保存在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/list.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/list.py)。

## 阻塞操作

阻塞操作是指各别指令在执行过程中会发生阻塞，并且在满足某特定条件下才返回相应值，阻塞指令包括：

* BLPOP
* BRPOP
* BRPOPLPUSH

比如，在执行 BLPOP 指令并且指定列表为空时，BLPOP 指令会阻塞当前的进程，只有当列表非空才会有返回，如以下代码：

```python
import redis
from threading import Thread

redis_config = {
    'host': 'redis-testing',
    'port': 6379,
    'db': 1,
}

rs = redis.Redis(**redis_config)


def lpush(key, *months):
    for month in months:
        rs.lpush(key, month)


def blpop(key, timeout):
    while True:
        ret = rs.blpop(key, timeout)
        if ret is None:
            break
        key, month = ret
        print(month)


if __name__ == '__main__':
    producer = Thread(target=lpush, args=('months', 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12))
    consumer = Thread(target=blpop, args=('months', 2))

    producer.start()
    consumer.start()

    producer.join()
    consumer.join()
```

输出：

```bash
# 输出是无序的，结果非一致
b'1'
b'2'
b'3'
b'4'
b'5'
b'8'
b'7'
b'9'
b'10'
b'12'
b'11'
b'6'
```

代码保存在：[https://github.com/a2htray/code-notebook/blob/main/Redis/app/list_block_ops.py](https://github.com/a2htray/code-notebook/blob/main/Redis/app/list_block_ops.py)。

上述代码以 BLPOP 为例，实现的功能为：

* 启动两个线程，分别为生产者（执行 lpush）和消费者（执行 blpop）
* lpush 会持续向 `months` 列表添加 12 个元素
* blpop 会持续向 `months` 列表取 1 个元素

## 小结

Python redis 包中定义的方法与 Redis 指令名称一致，熟悉 Redis 指令就等于会使用 redis 包。