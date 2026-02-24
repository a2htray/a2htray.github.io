---
title: "PAT 乙级 1007"
date: 2023-03-08T20:46:54+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1007.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/3

"""
PAT 乙级 1007
"""
import math


def is_prime(n):
    prime = True
    for i in range(3, int(math.sqrt(n)) + 1):
        if n % i == 0:
            prime = False
            break
    return prime


if __name__ == '__main__':
    nums = [n for n in range(3, int(input()) + 1)]
    nums = [n for n in nums if n % 2 != 0]

    nums = list(filter(is_prime, nums))
    if len(nums) <= 1:
        print(0)
        exit(0)

    count = 0
    i = 0
    j = 1
    while j != len(nums):
        if nums[j] - nums[i] == 2:
            count += 1
        i += 1
        j += 1

    print(count)
```

