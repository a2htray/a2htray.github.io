---
title: "PAT 乙级 1043"
date: 2023-03-25T15:50:40+08:00
draft: false
reward: false
categories:
  - 算法
  - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1043.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/25

"""
PAT 乙级 1043
"""

if __name__ == '__main__':
    line = input()
    nums = [0] * 6
    chars = ['P', 'A', 'T', 'e', 's', 't']
    for c in line:
        if c == 'P':
            nums[0] += 1
        if c == 'A':
            nums[1] += 1
        if c == 'T':
            nums[2] += 1
        if c == 'e':
            nums[3] += 1
        if c == 's':
            nums[4] += 1
        if c == 't':
            nums[5] += 1

    total = sum(nums)

    while total > 0:
        for i in range(6):
            if nums[i] != 0:
                print(chars[i], end='')
                nums[i] -= 1
                total -= 1

    print('')
```

