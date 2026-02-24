---
title: "PAT 乙级 1013"
date: 2023-03-08T20:47:11+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1013.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/3

"""
PAT 乙级 1013
"""
import math


def is_prime(num):
    for i in range(2, int(math.sqrt(num)) + 1):
        if num % i == 0:
            return False
    return True


if __name__ == '__main__':
    nums_str = input().split(' ')
    m, n = int(nums_str[0]), int(nums_str[1])
    prime_nums = [0] * n

    i = 0
    num = 2
    while i != n:
        if is_prime(num):
            prime_nums[i] = num
            i += 1
        num += 1

    output_nums = prime_nums[m - 1:n]
    for i in range(0, len(output_nums), 10):
        print(' '.join([str(v) for v in output_nums[i:i + 10]]))
```

