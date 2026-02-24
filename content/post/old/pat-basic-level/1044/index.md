---
title: "PAT 乙级 1044"
date: 2023-03-25T18:46:30+08:00
draft: false
reward: false
categories:
  - 算法
  - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1044.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/25

"""
PAT 乙级 1044
"""

mars_digits = [
    'tret',
    'jan', 'feb', 'mar', 'apr',
    'may', 'jun', 'jly', 'aug',
    'sep', 'oct', 'nov', 'dec',
]

mars_carries = [
    'tam', 'hel', 'maa', 'huh',
    'tou', 'kes', 'hei', 'elo',
    'syy', 'lok', 'mer', 'jou',
]


def is_earth(s: str):
    return s.isdigit()


def is_mars(s: str):
    return not is_earth(s)


def to_earth(s: str):
    tokens = s.split(' ')
    if len(tokens) == 1:
        if tokens[0] in mars_digits:
            return mars_digits.index(tokens[0])
        if tokens[0] in mars_carries:
            return (mars_carries.index(tokens[0]) + 1) * 13
    else:
        return (mars_carries.index(tokens[0]) + 1) * 13 + mars_digits.index(tokens[1])


def to_mars(s: str):
    num = int(s)
    if num < 13:
        return mars_digits[num]

    if num % 13 == 0:
        return f'{mars_carries[num // 13 - 1]}'
    else:
        return f'{mars_carries[num // 13 - 1]} {mars_digits[num % 13]}'


if __name__ == '__main__':
    n = int(input())
    lines = []
    for _ in range(n):
        line = input()
        lines.append(line)

    for line in lines:
        if is_earth(line):
            print(to_mars(line))
        else:
            print(to_earth(line))
```

