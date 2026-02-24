---
title: "PAT 乙级 1005"
date: 2023-03-08T20:46:39+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1005.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/2

"""
PAT 乙级 1005
"""


def compute_sequence(num):
    sequence = []
    while num != 1:
        if num % 2 == 0:
            num = num // 2
        else:
            num = (num * 3 + 1) // 2
        sequence.append(num)

    return sequence


if __name__ == '__main__':
    sequences = []
    n = int(input())
    nums = [int(num) for num in input().split(' ')]
    i = 0
    while i < n:
        sequence = compute_sequence(nums[i])
        sequences.append(sequence[:-1])
        i += 1

    key_nums = []
    for i, num in enumerate(nums):
        is_key = True
        for j, sequence in enumerate(sequences):
            if i == j:
                continue
            if num in sequence:
                is_key = False
                break

        if is_key:
            key_nums.append(num)

    key_nums.sort(reverse=True)
    print(' '.join([str(num) for num in key_nums]))
```

