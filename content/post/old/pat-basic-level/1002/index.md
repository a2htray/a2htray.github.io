---
title: "PAT 乙级 1002"
date: 2023-02-28T20:52:15+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1002.png)

## 代码

```python
# basic_1002.py
chinese_pinyins = [
    'ling',
    'yi',
    'er',
    'san',
    'si',
    'wu',
    'liu',
    'qi',
    'ba',
    'jiu',
]

if __name__ == '__main__':
    num_chars = input()

    total = 0
    for num_char in num_chars:
        total += int(num_char)

    output = []
    for num_char in str(total):
        output.append(chinese_pinyins[int(num_char)])

    print(' '.join(output))


```

## 运行

```bash
input:
1928374
output:
san si
```

