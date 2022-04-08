---
title: "列表模拟栈和队列"
date: 2022-03-27T06:05:41+08:00
draft: false
comment: false
featured: true
reward: false
toc: true
pinned: false
categories:
  - "数据集"
tags:
  - "Redis"
series:
  - "Redis 命令手册"
weight: 1
type: "docs"
---

栈（Stack）和队列（Queue）是编程中常用的两种数据结构，下面通过 Redis 的列表（List）类型来实现栈和队列。

<!--more-->

## 栈

栈是一种受限的线性表，即“**只能在一端进行插入和删除操作**”，其特点是后进先出（Last In First Out，LIFO）。假设列表的右端为栈顶（插入和删除的一端），则需要使用到 `RPUSH` 和 `RPOP` 两个命令。

```bash
127.0.0.1:6379> RPUSH stack 1 2 3 4 5 6
(integer) 6
127.0.0.1:6379> RPOP stack
"6"
127.0.0.1:6379> RPUSH stack 7 8 9
(integer) 8
127.0.0.1:6379> RPOP stack
"9"
127.0.0.1:6379> RPOP stack
"8"
```

上述命令的说明如下：

* `RPUSH stack 1 2 3 4 5 6` ：按插入的先后顺序，此时列表为 [1, 2, 3, 4, 5, 6]；
* `RPOP`：从列表右端弹出 1 个元素，该元素为 6，此时列表为 [1, 2, 3, 4, 5]；
* `RPUSH stack 7 8 9`：列表右端插入 3 个元素，此时列表为 [1, 2, 3, 4, 5, 7, 8, 9]；
* 最后的两次 `RPOP stack`：分别弹出元素 9 和 8，此时列表为 [1, 2, 3, 4, 5, 7]；

## 队列

队列也是一种受限的线性表，即“**一端只能插入，另一端只能删除**”，其特点是先进先出（First In First Out，FIFO）。假设列表的左端是队首（删除的一端），右端是队尾（插入的一端），则需要使用到 `LPUSH` 和 `LPOP` 两个命令。

```bash
127.0.0.1:6379> RPUSH queue 1 2 3 4 5 6
(integer) 6
127.0.0.1:6379> LPOP queue
"1"
127.0.0.1:6379> RPUSH queue 7 8 9
(integer) 8
127.0.0.1:6379> LPOP queue
"2"
127.0.0.1:6379> LPOP queue
"3"
```

上述命令的说明如下：

* `RPUSH queue 1 2 3 4 5 6`：按插入的先后顺序，此时列表为 [1, 2, 3, 4, 5, 6]；
* `LPOP queue`：从列表左端弹出 1 个元素，此时列表为 [2, 3, 4, 5, 6]；
* `RPUSH queue 7 8 9`：列表右端插入 3 个元素，此时列表为 [2, 3, 4, 5, 6, 7, 8, 9]；
* 最后两次 `LPOP queue`：分别从左端弹出元素 2 和 3，此时列表为 [4, 5, 6, 7, 8, 9]；