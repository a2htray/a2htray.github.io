---
title: "PAT 乙级 1033"
date: 2023-03-23T23:49:46+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1033.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/13

"""
PAT 乙级 1033
"""

if __name__ == '__main__':
    bad_keys = input()
    shift_bad = '+' in bad_keys

    for v in input():
        if v in bad_keys or v.upper() in bad_keys:
            continue

        if shift_bad and ord('A') <= ord(v) <= ord('Z'):
            continue

        print(v, end='')
```

