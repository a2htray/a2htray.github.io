---
title: "PAT 乙级 1014"
date: 2023-03-05T23:28:48+08:00
draft: false
reward: false
categories:
 - 算法
 - PAT
tags:
 - PAT Basic Level
 - PAT
---

![](1014.png)

## 代码

```python
#!/usr/bin/env python
# -*- coding: utf-8 -*-
# author: a2htray
# create date: 2023/3/4

"""
PAT 乙级 1014
"""

week_dict = {
    'A': 'MON',
    'B': 'TUE',
    'C': 'WED',
    'D': 'THU',
    'E': 'FRI',
    'F': 'SAT',
    'G': 'SUN',
}

hour_dict = {
    '0': '00',
    '1': '01',
    '2': '02',
    '3': '03',
    '4': '04',
    '5': '05',
    '6': '06',
    '7': '07',
    '8': '08',
    '9': '09',
    'A': '10',
    'B': '11',
    'C': '12',
    'D': '13',
    'E': '14',
    'F': '15',
    'G': '16',
    'H': '17',
    'I': '18',
    'J': '19',
    'K': '20',
    'L': '21',
    'M': '22',
    'N': '23',
}

if __name__ == '__main__':
    n_line = 4
    lines = []
    while n_line != 0:
        lines.append(input())
        n_line -= 1

    week, hour, second = '', '', ''

    j = 0
    for i, char in enumerate(lines[0]):
        if char == lines[1][i] and char in week_dict.keys():
            j = i
            week = week_dict[char]
            break

    for i in range(j+1, len(lines[0])):
        if lines[0][i] == lines[1][i] and lines[0][i] in hour_dict.keys():
            hour = hour_dict[lines[0][i]]
            break

    for i, char in enumerate(lines[2]):
        if char.isalpha() and char == lines[3][i]:
            second = '%.2d' % i
            break

    print(week + ' ' + hour + ':' + second)
```

