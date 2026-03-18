+++
date = '2026-03-17T18:57:49+08:00'
draft = true
title = 'MySQL：存储过程详解'
categories = ['后端技术', 'MySQL']
tags = ['MySQL', '存储过程']
toc = true
+++

存储过程在平时的开发中，其实用的不多。

在很久之前，也是我刚入行时候，有见过一个 CMS，它是典型的面向存储过程编程，是把所有的业务逻辑都包装到了一个又一个存储过程中。因此，后台代码只要根据
业务名调用存储过程，从而完成相应的业务逻辑操作。这就要求：

1. 业务可细化，至少能用一个存储过程表示
2. 开发人员对存储过程要足够掌握
3. 存储过程间要包含业务逻辑链，如下单 -> 出仓 -> 发货 -> 收货 -> 评价

当然，这种业务实现方式已然过时，但作为开发者，还是要对存储过程有一定的概念。通常，我都是遇到就翻下文档、查下语法，这一次想系统性地学习并记录存储过程。

> 期望：以后遇到不明白了，翻看这篇笔记，能找到答案。

## 环境配置

依旧是以 docker 的方式启动 MySQL 服务，但这次我想以此为契机，将以后要实验的 MySQL 操作都在一个环境、一个数据集上，所以采用了数据卷的方式持久化
数据库数据。

**创建数据卷**

```bash
$ docker volume create mysql8_testing_data
```

**启动 MySQL 服务**

```bash
$ docker run -d \
  --name mysql8 \                                       # 容器名称（自定义）
  --restart=always \                                    # 开机自启
  -p 3306:3306 \                                        # 端口映射：宿主机3306 → 容器3306（外部连接用）
  -v mysql8_testing_data:/var/lib/mysql \               # 数据卷映射：数据持久化
  -e MYSQL_ROOT_PASSWORD=root123456 \                   # 设置root密码（必填）
  -e MYSQL_ROOT_HOST=% \                                # 允许所有主机连接root（核心配置）
  mysql:8.0 \
  --character-set-server=utf8mb4 \                      # 设置默认字符集
  --collation-server=utf8mb4_unicode_ci \
  --default-authentication-plugin=mysql_native_password # 兼容旧客户端认证方式（避坑！）
```


## 数据导入

## 语句

## 语法



