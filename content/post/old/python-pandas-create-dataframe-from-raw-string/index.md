---
title: "从字符串中创建 DataFrame"
date: 2022-11-17T22:24:05+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Python
tags:
 - Pandas
 - DataFrame
 - Python
---

一般情况下，你们会通过文件（CSV、Excel等） 或 Python 的内置结构（字典）来创建 DataFrame 对象。但有时，数据是字符串的形式，如何将其转换成
DataFrame 对象？

<!--more-->

答案：将字符串强制转换成 io.StringIO，再作为 pd.read_csv 的参数。

```python
import pandas as pd
from io import StringIO

data = '''人口基本情况,1982年,1990年,2000年,2020年,2021年
0-14岁人口,34156,31670,29024,25277,24721
15-64岁人口,62517,76260,88847,96871,96481
65岁以上人口,4981,6403,8872,19064,20059'''

df = pd.read_csv(StringIO(data), index_col=0)
```
