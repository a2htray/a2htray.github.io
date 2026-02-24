---
title: "PAT 乙级 1001"
date: 2023-02-28T20:34:14+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1001.png)

## 代码

```python
# basic_1001.py

if __name__ == '__main__':
    n = int(input())
    step = 0

    while n != 1:
        if n % 2 == 0:
            n = n // 2
        else:
            n = (3 * n + 1) // 2
        step += 1

    print(step)
```

## 运行

```bash
input:
100
output:
18
```

