---
title: "PAT 乙级 1026"
date: 2023-03-09T23:50:37+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1026.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/9

"""
PAT 乙级 1026
"""

if __name__ == '__main__':
    start_clock, end_clock = map(int, input().split(' '))
    duration = 1.0 * (end_clock - start_clock) / 100

    h = int(duration // 3600)
    m = int(duration % 3600 // 60)
    s = str(duration % 3600 % 60)

    tokens = s.split('.')

    add = 1 if int(tokens[1][0]) >= 5 else 0
    s = int(tokens[0]) + add

    print('%02d:%02d:%02d' % (h, m, s))
```

