+++
date = '2026-04-09T14:28:13+08:00'
draft = false
title = 'SQLAlchemy 查询按字段降序或升序'
categories = ['后端技术', 'SQLAlchemy']
tags = ['Python', 'ORM', 'SQLAlchemy', 'PostgreSQL']
toc = true
+++

在数据库中，查询按某一字段降序和升序的关键字，分别是 **DESC** 和 **ASC**，不指定则默认按升序，也就是 **ASC**，如下面三个查询：

*不指定*

```postgresql
SELECT created_at, id FROM llm_models ORDER BY created_at;
```

![](/imgs/learn-backend/sql-query-01.png)

*指定 ASC*

```postgresql
SELECT created_at, id FROM llm_models ORDER BY created_at ASC;
```

![](/imgs/learn-backend/sql-query-02.png)

*指定 DESC*

```postgresql
SELECT created_at, id FROM llm_models ORDER BY created_at DESC;
```

![](/imgs/learn-backend/sql-query-03.png)


## 在 SQLAlchemy 中

SQLAlchemy 是 Python 中的数据库 ORM 框架，在进行降序和升序时，有两种方式：

1. 使用 desc() 或 asc() 函数
2. 字段调用 desc() 或 asc() 函数

因为数据库的默认行为是按升序排列，所以示例以降序给出。

**第一种**

```python
from sqlalchemy import desc

query.order_by(desc(Model.created_at))
```

**第二种**

```python
query.order_by(Model.created_at.desc())
```

如果存在多个排序字段，则依次作为参数传入。