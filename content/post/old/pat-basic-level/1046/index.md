---
title: "PAT 乙级 1046"
date: 2023-03-26T11:32:42+08:00
draft: false
reward: false
categories:
  - 算法
  - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1046.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/26

"""
PAT 乙级 1046
"""

if __name__ == '__main__':
    n = int(input())
    a = 0
    b = 0

    for i in range(n):
        tokens = list(map(int, input().split(' ')))
        if tokens[1] == tokens[3]:
            continue
        total = tokens[0] + tokens[2]
        if tokens[1] == total:
            b += 1
            continue
        if tokens[3] == total:
            a += 1
            continue

    print(a, b)
```

