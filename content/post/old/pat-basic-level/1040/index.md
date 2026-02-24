---
title: "PAT 乙级 1040"
date: 2023-03-25T14:56:07+08:00
draft: false
reward: false
categories:
  - 算法
  - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1040.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/25

"""
PAT 乙级 1039
"""

if __name__ == '__main__':
    line = input()

    num_t = 0
    num_at = 0
    num_pat = 0

    for i in range(len(line) - 1, -1, -1):
        if line[i] == 'T':
            num_t += 1
        if num_t != 0 and line[i] == 'A':
            num_at += num_t
        if num_at != 0 and line[i] == 'P':
            num_pat += num_at

    print(num_pat % 1000000007)
```

