---
title: "通过索引加快数据表更新速度"
date: 2022-04-14T11:16:11+08:00
draft: true
reward: false
tags:
 - PostgreSQL
 - 索引
 - 数据更新
 - sqlx
categories:
 - 数据库
 - Go
images:
 - "images/postgresql-index.jpg"
---

当开发过程中，由于业务逻辑的调整，已有的数据库表设计也会相应地改变，所以需要对修改后的表做数据更新。数据更新的速度往往受限于数据的规模，数据量越大更新速度越慢。本文示例来自于实际工作内容，但会以其它数据进行演示。

## 需求

有一个 200W 条记录的 `work_problem` 表 ，业务要求需要新增了一个 `handler` 字段，而 `handler` 字段的值其实是存在于当前记录的其它字段中的（相当于做了表的冗余设计），数据更新的值要通过其它信息而来。原先的 `work_problem` 表的结构如下：

```sql
CREATE TABLE IF NOT EXISTS work_problem
(
    id bigint,
    item jsonb,
    type varchar(10) -- 取 JavaScript|React|Project
);
```

`type` 字段取值不同，`item` 字段存储的信息也不同，规则如下：

* `type = JavaScript` 时，`item` 字段存储的是该记录所属 React 记录的 ID，如 `{"react": react_record_id}`；
* `type = React` 时，`item` 字段存储的是该记录所属 Project 记录的 ID，如 `{"project": project_record_id}`；
* `type = Project` 时，Project 的 ID 用 `id` 字段存储，所以 `item` 存储空对象，如 `{}`。

### Go 程序填充数据



## 问题

## 解决方案

## 总结



