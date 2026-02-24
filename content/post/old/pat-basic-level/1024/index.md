---
title: "PAT 乙级 1024"
date: 2023-03-08T22:51:25+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1024.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/8

"""
PAT 乙级 1024
"""
from typing import List


class SciNumber:
    def __init__(self, s: str):
        self.sign = '-' if s[0] == '-' else ''
        nums, exp = s[1:].split('E')
        self.nums = [v for v in nums]
        self.left_move = exp[0] == '-'
        self.exp = int(exp[1:])

        # print(self.sign, self.nums, self.left_move, self.exp)

    def to_number(self) -> List[str]:
        if self.left_move:
            ret = ['0'] * self.exp + self.nums
            del ret[ret.index('.')]
            return ret[0:1] + ['.'] + ret[1:]
        else:
            n = len(self.nums[2:])
            if self.exp >= n:
                ret = self.nums + ['0'] * (self.exp - n)
                del ret[1]
            else:
                ret = self.nums
                for i in range(1, 1 + self.exp):
                    ret[i], ret[i + 1] = ret[i + 1], ret[i]

            if ret[0] == '0':
                ret = ret[1:]
            return ret


if __name__ == '__main__':
    sci_number = SciNumber(input())
    print(sci_number.sign + ''.join(sci_number.to_number()))
```

