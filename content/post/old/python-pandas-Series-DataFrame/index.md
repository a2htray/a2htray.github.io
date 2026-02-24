---
title: "创建 Pandas 的 Series 和 DataFrame "
date: 2022-10-15T19:25:04+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Python
tags:
  - Pandas
  - DataFrame
  - Series
  - Python
---

Series 和 DataFrame 是 Pandas 中两种重要的数据结构，也是我们操作和分析的主要对象。其中 Series 是一种类似于数组、列表或表格中一列的
一维数据对象，DataFrame 则可以表示表格化的数据对象，可由多个 Series 对象组成。

本文主要摘录 Series 和 DataFrame 两种数据结结构的创建方法以及一些注意事项。

<!--more-->

## Series 的创建

Pandas 提供 Series 来表示一系列的元素，与常规一维数组只能以整数作为索引不同，Series 的索引类型可以是任一类型，其索引值与值的映射结构类似于 Python 中的字典。

### Series 的构造函数

Series 的构造函数的声明如下：

```python
def __init__(
    self,
    data=None,
    index=None,
    dtype: Dtype | None = None,
    name=None,
    copy: bool = False,
    fastpath: bool = False,
):
    pass
```

### 从列表中创建 Series

```python
names = ['Jimmy', 'Andy', 'Tony', 'John']
names_series = pd.Series(names)
```

`pd.Series` 的第 1 个参数名为 `data`，即上述调用过程等价于 `names_series = pd.Series(data=names)`。因为列表中没有相应的值可以作为 Series 的 index 值，所以 Pandas 会以 0 起始为每 1 个元素进行索引。如果指定了 index 参数值，可以修改 Series 中的索引值。

```python
names_series = pd.Series(names, index=['ST0001', 'ST0002', 'ST0003', 'ST0004'])
```

此时可以通过 `names_series['ST0002']` 访问 `'Andy'` 元素。

### 从字典中创建 Series

```python
sales = {
    'January': 1000,
    'February': 900,
    'March': 1200,
}
sales_series = pd.Series(sales)
```

Python 中的字典是 key-value 结构，所以字典值作为 `data` 参数值时，其 key 值会作为 Series 的 index 值。此时可以通过 `sales['January']` 访问到 `1000` 元素。

### 从 np.ndarray 中创建 Series

从 np.ndarray 中创建 Series 与从列表中创建 Series 类似：

```python
x_axis_values = np.arange(0, 7)
axis_series = pd.Series(x_axis_values, index=['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday'])
```

### 从标量中创建 Series

```python
scaled_series = pd.Series(7, index=['01', '02', '03'])
```

此时 scaled_series 中一共有 3 个元素，各元素的 index 分别为 '01'、'02' 和 '03'。标量作为 `data` 参数值时，Pandas 会按 `index` 参数
值的长度进行扩展。如果没有指定 `index` 值时，索引值则为 0 且长度为 1。

## DataFrame 的创建

Pandas 提供 DataFrame 结构可用于表示表格化的数据。

### DataFrame 的构造函数

```python
def __init__(
    self,
    data=None,
    index: Axes | None = None,
    columns: Axes | None = None,
    dtype: Dtype | None = None,
    copy: bool | None = None,
):
    pass
```

相较于 Series 的构造函数，DataFrame 的构造函数可以指定列名（`columns`）。因为 DataFrame 可以视为一个或多个 Series 的集合，同样的，
`columns` 值可以视为每一个 Series 的 `name` 参数的传值。

### 从字典（列表为元素）中创建 DataFrame

```python
scores = {
    'Jimmy': [89, 88, 93],
    'Andy': [78, 90, 99],
    'Tony': [83, 85, 87],
    'John': [87, 67, 70],
}

df = pd.DataFrame(data=scores)
```

`scores` 的 key 值会作为 DataFrame 的列名。

### 指定 index、columns 值

`index` 参数用于设置 DataFrame 每行的索引值，`columns` 参数用于设置列名。

```python
df = pd.DataFrame([
    [89, 88, 93],
    [78, 90, 99],
    [83, 85, 87],
    [87, 67, 70],
], index=['Jimmy', 'Andy', 'Tony', 'John'], columns=['Math', 'English', 'Chinese'])
```

上述中，`df` 是一个 4 行 3 列的表格化对象。如果 `data` 参数值是一个字典，则 `columns` 参数还可以起到过滤的作用。
