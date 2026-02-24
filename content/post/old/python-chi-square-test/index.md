---
title: "卡方检验 - 检验特征对是否相关"
date: 2022-04-11T18:57:39+08:00
draft: false
reward: false
categories:
 - 数据分析
 - 机器学习
tags:
 - Chi-Square
 - scipy
 - NumPy
 - Pandas
 - Python
images:
 - "images/chi-square-test.webp"
---

在本文开头，贴一段百科对卡方检验基本原理的介绍：

<p style="color: blue; text-decoration: underline;">卡方检验就是统计样本的实际观测值与理论推断值之间的偏离程度，实际观测值与理论推断值之间的偏离程度就决定卡方值的大小，如果卡方值越大，二者偏差程度越大；反之，二者偏差越小；若两个值完全相等时，卡方值就为 0，表明理论值完全符合。</p>

由此可见，卡方检验刻画的是一种偏离程度。那么在相关性计算中也可以利用卡方检验计算出显著性来判断两个特征是否相关。

## 卡方检验

卡方检验的步骤如下：

1. 定义 H0 和 H1 假设；
2. 根据领域知识定义显著性水平 $\alpha$，一般取 0.05，表示有 5% 的容错；
3. 计算卡方值；
4. 计算显著性水平，小于 $\alpha$ 则拒绝 H0 接受 H1；

## 离散型特征对

离散型特征对是指特征为离散值的两维向量，如[帕尔默企鹅数据集](/posts/python-palmer-archipelago-penguin-testing/)中的特征对（species，island）。下面演示特征列（species，island）是否存在相关性。

提出假设：

1. H0：特征 species 和特征 island 不相关（独立）；
2. H1：特征 species 和特征 island 相关；

### 频次统计

首先，根据出现的特征值对进行频次统计：

```python
import numpy as np
import pandas as pd

# 读取数据集
# 特征 species，island 并不包含缺失值
df = pd.read_csv('./dataset/penguins_size.csv', na_values=['NA', '.'])

# 得到 species 和 island 所有的取值
# speices 值为行名，island 为列名
index = df['species'].unique()
columns = df['island'].unique()

# 初始化一个频次矩阵
count_matrix = np.zeros((index.shape[0], columns.shape[0]))
# 遍历所有的 species 值
for i, species in enumerate(index):
    # 当 species 为某个特定值时，计算此时各 island 值出现的次数
    counts = df[df['species']==species]['island'].value_counts()
    # 在对应的行和列设置次数值
    for j, island in enumerate(columns):
        if island in counts.index:
            count_matrix[i][j] = counts[island]
        else:
            count_matrix[i][j] = 0

df_counts = pd.DataFrame(count_matrix, columns=columns, index=index)
print(df_counts)
```

```bash
           Torgersen  Biscoe  Dream
Adelie          52.0    44.0   56.0
Chinstrap        0.0     0.0   68.0
Gentoo           0.0   124.0    0.0
```

输出表示：

1. 当 species 为 `Adelie` 时，表示 `Torgersen` 出现 52次，`Biscoe` 出现 44 次，`Dream` 出现 56 次；
2. 其余类似；

### 计算估计值

估计值的计算公式如下：
$$
E_{ij} = \frac{R_i \times C_j} {N}
$$
其中 $R_i$ 表示第 $i$ 行的总和；$C_j$ 表示第 $j$ 列的总和；$N$ 表示所有值的总和。代码如下：

```python
# 行和
rows_total = df_counts.sum(axis=1).values
# 列和
cols_total = df_counts.sum(axis=0).values
# 总和
total = df_counts.values.sum()

# 根据公式计算估计值
estimated_count_matrix = np.zeros((index.shape[0], columns.shape[0]))
for i in range(index.shape[0]):
    for j in range(columns.shape[0]):
        estimated_count_matrix[i][j] = rows_total[i]*cols_total[j]/total

df_estimated_counts = pd.DataFrame(estimated_count_matrix, columns=columns, index=index)
print(df_estimated_counts)
```

```bash
           Torgersen     Biscoe      Dream
Adelie     22.976744  74.232558  54.790698
Chinstrap  10.279070  33.209302  24.511628
Gentoo     18.744186  60.558140  44.697674
```

### 计算卡方值

卡方值的计算公式如下：
$$
{\chi}^2 = \sum \frac {(O_{ij} - E_{ij})^2} {E_{ij}}
$$
其中 $O_{ij}$ 为实际的频次值。代码如下：

```python
df_chisq = np.power(df_counts - df_estimated_counts, 2) / estimated_count_matrix
print(df_chisq)
```

```bash
           Torgersen     Biscoe      Dream
Adelie     36.660955  12.312759   0.026691
Chinstrap  10.279070  33.209302  77.156789
Gentoo     18.744186  66.462901  44.697674
```

```go
chi = df_chisq.values.sum()
print(chi)
```

```bash
299.55032743148195
```

### 计算显著性

先计算自由度，公式如下：
$$
degree = (r - 1) \times (c - 1)
$$
其中 $r$ 为行数；$c$ 为列数。计算显著性的代码如下：

```python
from scipy import stats

degree = (len(index) - 1) * (len(columns) - 1)
pvalue = 1 - stats.chi2.cdf(x=chi, df=degree)
print(pvalue)
```

```bash
0.0
```

p 值小于 0.05，所以拒绝 H0，说明特征 species 和特征 island 相关。

## 完整代码

```python
import numpy as np
import pandas as pd
from scipy import stats

df = pd.read_csv('./dataset/penguins_size.csv', na_values=['NA', '.'])

index = df['species'].unique()
columns = df['island'].unique()

count_matrix = np.zeros((index.shape[0], columns.shape[0]))
for i, species in enumerate(index):
    counts = df[df['species']==species]['island'].value_counts()
    for j, island in enumerate(columns):
        if island in counts.index:
            count_matrix[i][j] = counts[island]
        else:
            count_matrix[i][j] = 0

df_counts = pd.DataFrame(count_matrix, columns=columns, index=index)
print(df_counts)

rows_total = df_counts.sum(axis=1).values
cols_total = df_counts.sum(axis=0).values
total = df_counts.values.sum()

estimated_count_matrix = np.zeros((index.shape[0], columns.shape[0]))
for i in range(index.shape[0]):
    for j in range(columns.shape[0]):
        estimated_count_matrix[i][j] = rows_total[i]*cols_total[j]/total

df_estimated_counts = pd.DataFrame(estimated_count_matrix, columns=columns, index=index)
print(df_estimated_counts)

df_chisq = np.power(df_counts - df_estimated_counts, 2) / estimated_count_matrix
print(df_chisq)

chi = df_chisq.values.sum()
print(chi)

degree = (len(index) - 1) * (len(columns) - 1)
pvalue = 1 - stats.chi2.cdf(x=chi, df=degree)
print(pvalue)
```



## 参考

1. [STAT #3: Chi-Squared Test(卡方检验)](https://zhuanlan.zhihu.com/p/394084469)
2. [A beginner’s guide to Chi-square test in python from scratch](https://analyticsindiamag.com/a-beginners-guide-to-chi-square-test-in-python-from-scratch/)

