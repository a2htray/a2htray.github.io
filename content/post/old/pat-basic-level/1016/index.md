---
title: "PAT 乙级 1016"
date: 2023-03-05T23:28:55+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1016.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/5

"""
PAT 乙级 1016
"""


def count(string, char):
    c = 0
    for v in string:
        if v == char:
            c += 1
    return c


def pad(char, num):
    if num == 0:
        return '0'
    ret = ''
    while num != 0:
        ret += char
        num -= 1
    return ret


if __name__ == '__main__':
    strings = input().split(' ')
    a_count = count(strings[0], strings[1])
    b_count = count(strings[2], strings[3])

    print(int(pad(strings[1], a_count)) + int(pad(strings[3], b_count)))

```

