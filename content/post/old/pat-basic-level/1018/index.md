---
title: "PAT 乙级 1018"
date: 2023-03-08T21:08:06+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1018.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/6

"""
PAT 乙级 1018
"""

if __name__ == '__main__':
    n = int(input())
    first = [0, 0, 0]
    second = [0, 0, 0]
    first_win = 0
    second_win = 0

    alphabet = ['B', 'C', 'J']

    for _ in range(n):
        p1, p2 = input().split()
        if p1 == 'C' and p2 == 'J':
            first[1] += 1
            first_win += 1
        elif p1 == 'C' and p2 == 'B':
            second[0] += 1
            second_win += 1
        elif p1 == 'B' and p2 == 'J':
            second[2] += 1
            second_win += 1
        elif p1 == 'B' and p2 == 'C':
            first[0] += 1
            first_win += 1
        elif p1 == 'J' and p2 == 'B':
            first[2] += 1
            first_win += 1
        elif p1 == 'J' and p2 == 'C':
            second[1] += 1
            second_win += 1

    print(first_win, n - first_win - second_win, second_win)
    print(second_win, n - first_win - second_win, first_win)

    print(alphabet[first.index(max(first))], alphabet[second.index(max(second))])
```

