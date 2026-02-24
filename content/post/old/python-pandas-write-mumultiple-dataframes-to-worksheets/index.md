---
title: "Write Mumultiple Dataframes to Worksheets"
date: 2022-09-21T00:01:28+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Python
tags:
 - Pandas
 - xls
---

In normal work, I usually use `Pandas` as my excel read/write utility. 

Here is an example for how to write multiple dataframes to worksheets. We need to use `pd.ExcelWriter` method.

<!--more-->

```python
import pandas as pd

writer = pd.ExcelWriter('./multiple-dataframes-in-a-single-xls.xls')

df1 = pd.DataFrame(columns=[
    'Book',
    'Author',
], data=[
    ('Pride and Prejudice', 'Jane Austen'), 
    ('Healing Her Heart', 'Laura Scott'),
])

df2 = pd.DataFrame(columns=[
    'Website',
    'Scale',
], data=[
    ('Alibaba', 'Larger'),
    ('Google', 'Larger'),
])

df1.to_excel(writer, sheet_name='Book', index=False)
df2.to_excel(writer, sheet_name='Website', index=False)

writer.save()
```

If you get into trouble that the script raises an Error `ModuleNotFoundError: No module named 'xlwt'`, you need to run `pip install xlwt` to install the xlwt package. It is important to note that the xlwt package is no longer maintained, and the xlwt engine will been removed in a furture version of pandas.

At last, don't forget to use `writer.save` method to flush the buffers into the output file.


