---
title: "PAT 乙级 1006"
date: 2023-03-08T20:46:43+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1006.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/3

"""
PAT 乙级 1006
"""

if __name__ == '__main__':
    chars = [*input()]
    chars.reverse()

    tail = []
    for i in range(1, int(chars[0]) + 1):
        tail.append(str(i))

    bs_list = [[], []]
    for i, char in enumerate(chars[1:]):
        if i == 0:
            bs_list[0] = ['S'] * int(char)
        else:
            bs_list[1] = ['B'] * int(char)

    print(''.join(bs_list[1] + bs_list[0] + tail))
```

