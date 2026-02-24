---
title: "PAT 乙级 1037"
date: 2023-03-23T23:52:23+08:00
draft: false
reward: false
categories:
  - 算法
  - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1037.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/13

"""
PAT 乙级 1037
"""


def get_gsk(s):
    return list(map(int, s.split('.')))


if __name__ == '__main__':
    gsk1, gsk2 = map(get_gsk, input().split(' '))

    price1 = gsk1[0] * 17 * 29 + gsk1[1] * 29 + gsk1[2]
    price2 = gsk2[0] * 17 * 29 + gsk2[1] * 29 + gsk2[2]

    if price2 >= price1:
        flag = ''
        diff = price2 - price1
    else:
        flag = '-'
        diff = price1 - price2

    print(f'{flag}{diff // 29 // 17}.{diff // 29 % 17}.{diff % 29}')
```

