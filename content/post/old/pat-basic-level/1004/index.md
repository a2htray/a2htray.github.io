---
title: "PAT 乙级 1004"
date: 2023-03-08T20:45:38+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1004.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/1

"""
PAT 乙级 1004
"""


def get_score(student):
    return student[2]


if __name__ == '__main__':
    n = int(input())

    students = []
    while n != 0:
        tokens = input().split(' ')
        students.append((tokens[0], tokens[1], int(tokens[2])))
        n -= 1

    students.sort(key=get_score)

    print(students[len(students) - 1][0], students[len(students) - 1][1])
    print(students[0][0], students[0][1])
```

