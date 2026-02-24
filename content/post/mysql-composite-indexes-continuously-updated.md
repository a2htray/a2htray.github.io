+++
date = '2025-12-17T14:40:10+08:00'
draft = false
title = 'MySQL 联合索引（持续更新）'
categories = ['后端技术', 'MySQL']
tags = ['MySQL', '索引', '持续更新']
toc = true
+++

## Timeline

* [2025-12-17 新增联合索引基础知识](#基础知识2025-12-17)

## 基础知识（2025-12-17）

联合索引，又叫多列索引，是一种可以**包含多个列**的索引技术，使用得当，可以提升过滤和排序执行效率。

最多可包含 16 个字段<sub>[](https://www.mysqltutorial.org/mysql-index/mysql-composite-index/)</sub>。

创建联合索引示例如下：

```sql
-- 建表，同时创建联合索引
CREATE TABLE students (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(64) NOT NULL,
    gender SMALLINT DEFAULT 0,
    grade VARCHAR(16) NOT NULL,
    -- 此处
    INDEX idx_name_grade (name, grade)
);

-- 表建完后，新建索引
CREATE INDEX idx_name_grade ON students (name, grade);
```

> 新建唯一联合索引，则加上 **UNIQUE** 关键字，如：
> * UNIQUE INDEX idx_name_grade (name, grade)
> * CREATE UNIQUE INDEX idx_name_grade ON students (name, grade)

重点：

* 应用<span style="color: blue;">**最左匹配原则**</span>：对 (name)、(name, grade) 有效，但对 (grade, name) 无效
  * **主流数据库优化器会对条件进行重排，所以条件顺序不影响命中索引**
* 联合索引 (name, grade) 和 (grade, name) 不同
* 索引副作用：对 INSERT、UPDATE、DELETE 操作有额外开销
* 应业务要求的过滤逻辑创建联合索引，<span style="color: blue;">**若业务场景下，单个字段的查询更多，则创建单个索引**</span>



