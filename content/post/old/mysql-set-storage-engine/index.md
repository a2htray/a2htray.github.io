---
title: "MySQL 设置存储引擎的 3 种方法"
date: 2022-04-21T13:35:02+08:00
draft: false
reward: false
tags:
 - MySQL
 - 存储引擎
categories:
 - 后端技术
 - MySQL
---

在 MySQL 5.5 之前，默认存储引擎为 MyISAM，之后版本的默认存储引擎为 InnoDB。

选择一个合适的存储引擎至关重要。

<!--more-->

## 存储引擎

根据是否支持事务，MySQL 的存储引擎可以分为：

1. 事务型；
2. 非事务型；

表 1 为 MySQL 支持的所有存储引擎以及各存储引擎的基本介绍。

<p style="text-align: center;"><i>表 1 不同的存储引擎</i></p>

| 存储引擎      | 支持事务 | 特点                                       | 适用场景                      |
| --------- | ---- | ---------------------------------------- | ------------------------- |
| InnoDB    | 是    | 1. 支持行锁、灾难恢复、多版本并发控制；<br />2. 支持外键、字段约束； | 1. 适用于绝大部分场景；             |
| MyISAM    | 否    | 1. 读写速度快；<br />2. 支持表锁；<br />3. 支持 B 树索引、聚簇索引、全文搜索索引；<br />4. 支持地理数据及其索引；<br />5. 不支持哈希索引、外键、多版本并发控制；<br />6. 存储限制为 256TB； | 1. 读写频繁的应用；<br />2. 数据仓库； |
| Memory    | 否    | 1. 内存数据库；<br />2. 相较于 MyISAM，读写速度更快；<br />3. 支持表锁；<br />4. 非持久化数据；<br />5. 不支持多版本并发控制； | 1. 快存快取                   |
| CSV       | 否    | 1. 通用格式，易于集成；<br />2. 不支持索引；<br />3. 不支持分区；<br />4. 表的所有字段都要设置 `not null` | /                         |
| Merge     | 否    | 1. 底层使用 MyISAM 存储引擎；<br />               | 1. 数据仓库<br />             |
| Archive   | 否    | 1. 插入数据后，数据会被压缩；                         | 1. 存储历史数据；                |
| Federated | 否    | 1. 集群式管理 MySQL 实例；<br />                 | 1. 分布式环境；                 |
| Blackhole | 否    | 1. 可以向表插入数据，但查询只会返回空结果；<br />2. 支持所有的索引类型；<br />3. 不支持分区； | /                         |
| Example   | 否    | 啥也不是存储引擎                                 | /                         |

## 设置方法

以 InnoDB 为例，可通过以下 3 种方法设置表的存储引擎：

1. my.cnf 配置项；
2. SET STORAGE_ENGINE；
3. 创建表时；

### my.cnf 配置项

在 my.cnf 文件或其它引入的配置中，修改 `[mysqld]` 中的 `default-storage-engine` 的值，如：

```bash
[mysqld]
default-storage-engine = InnoDB
```

### SET STORAGE_ENGINE

在执行脚本文件前，先通过 `SET` 设置使用的存储引擎：

```mysq
SET STORAGE_ENGINE = InnoDB;
```

### 创建表时

在创建数据库表时，指定 `ENGINE`：

```mysql
CREATE TABLE IF NOT EXISTS test_name (
  id int
) ENGINE = InnoDB;
```

## 总结

1. MySQL 支持多种不同的存储引擎，InnoDB 存储引擎适用于绝大多数场景，并且支持事务、多版本并发控制；
2. 可以在 3 个层次上对存储引擎进行修改，即：
   1. 服务器层，`my.cnf` 配置项；
   2. 会话层，在当前会话中使用 `SET` 设置存储引擎；
   3. 脚本层，在创建或修改表时声明存储引擎；

## 参考

1. [MySQL storage engines](https://zetcode.com/mysql/storageengines/)
2. [官方文档 Alternative Storage Engines](https://dev.mysql.com/doc/refman/8.0/en/storage-engines.html)

