---
title: "PAT 乙级 1031"
date: 2023-03-12T23:10:04+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1031.png)

## 代码

```bash
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/12

"""
PAT 乙级 1031
"""
weights = [7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2]
m = ['1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2']


def weight_sum(digits):
    total = 0
    for i, digit in enumerate(digits):
        total += digit * weights[i]
    return total


if __name__ == '__main__':
    n = int(input())

    not_passed = []
    for _ in range(n):
        line = input()
        if len(line) != 18:
            not_passed.append(line)
            continue

        all_digit = True
        for c in line[:17]:
            if not c.isdigit():
                all_digit = False
        if not all_digit:
            not_passed.append(line)
            continue

        ws = weight_sum(list(map(int, line[:17])))
        z = ws % 11
        if m[z] != line[17]:
            not_passed.append(line)

    if len(not_passed) == 0:
        print('All passed')
    else:
        for item in not_passed:
            print(item)
```

