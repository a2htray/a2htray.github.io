---
title: "MySQL 修改列的注释信息"
date: 2022-04-25T14:04:54+08:00
draft: false
toc: true
reward: false
categories:
 - 后端技术
 - MySQL
tags:
 - stackoverflow
 - MySQL
---

小小的修改列的注释信息也能引发一些思考。

<!--more-->

正确的修改方式要带上原来列的定义，如：

```sql
 ALTER TABLE `user` CHANGE `id` `id` INT( 11 ) COMMENT 'id of user';
```

高票答案的几个回复：

> Note that altering a comment will cause a full resconstruction of the table. So you may choose to live without it on very big table.
>
> **修改列注释会重构表** **（对于大表慎用）**
>
> That is not (or no longer) true, as long as the column definition matches the existing definition exactly. Comments can be added without causing table reconstruction.
>
> **如果和之前的列定义一致，修改列注释不会重构表**



[**Alter MySQL table to add comments on columns**](https://stackoverflow.com/questions/2162420/alter-mysql-table-to-add-comments-on-columns)