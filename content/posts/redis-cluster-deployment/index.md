---
title: "Redis 集群配置过程"
date: 2022-04-07T14:07:35+08:00
draft: true
comment: false
featured: true
reward: false
toc: true
pinned: false
categories:
  - "数据集"
  - "运维"
tags:
  - "Redis"
weight: 1
---

Redis 集群是基于“主从复制”特性之上的分布式 Redis 版本，可提供高并发、高性能、高可用的数据库服务。Redis 集群突破了单台服务器的内存局限，集群中的每一个节点都可以存储数据，同时维护着 "key-node" 的映射表。本文记录了 3 主 3 从的 Redis 集群的配置过程，主要内容包括：

1. Redis 集群的配置过程；
2. 集群相关命令；
3. Go 存取集群数据；

<!--more-->

**环境信息**

```bash
$ lsb_release -a
No LSB modules are available.
Distributor ID: Ubuntu
Description:    Ubuntu 20.04.2 LTS
Release:        20.04
Codename:       focal
```

```bash
$ redis-server -v
Redis server v=6.2.6 sha=00000000:0 malloc=jemalloc-5.1.0 bits=64 build=9c9e426e2f96cc51
```

## 配置过程

3 主使用的端口分别为 6380、6381 和 6382，对应的 3 从使用的端口为 16380、16381 和 16382。创建配置文件以及工作目录：

```bash
$ mkdir /etc/redis/cluster
$ chown redis.redis /etc/redis/cluster
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6380.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6381.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6382.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-16380.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-16381.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-16382.conf
```

各配置文件的修改内容如下：

```bash
# /etc/redis/cluster/redis-6380.conf
port 6380
cluster-enabled yes
pidfile /run/redis/redis-server-6380.pid
logfile /var/log/redis/redis-server-6380.log
dbfilename dump-6380.rdb
cluster-config-file nodes-6380.conf
```

```bash
# /etc/redis/cluster/redis-6381.conf
port 6381
cluster-enabled yes
pidfile /run/redis/redis-server-6381.pid
logfile /var/log/redis/redis-server-6381.log
dbfilename dump-6381.rdb
cluster-config-file nodes-6381.conf
```

```bash
# /etc/redis/cluster/redis-6382.conf
port 6382
cluster-enabled yes
pidfile /run/redis/redis-server-6382.pid
logfile /var/log/redis/redis-server-6382.log
dbfilename dump-6382.rdb
cluster-config-file nodes-6382.conf
```

```bash
# /etc/redis/cluster/redis-16380.conf
port 16380
cluster-enabled yes
pidfile /run/redis/redis-server-16380.pid
logfile /var/log/redis/redis-server-16380.log
dbfilename dump-16380.rdb
cluster-config-file nodes-16380.conf
```

```bash
# /etc/redis/cluster/redis-16381.conf
port 16381
cluster-enabled yes
pidfile /run/redis/redis-server-16381.pid
logfile /var/log/redis/redis-server-16381.log
dbfilename dump-16381.rdb
cluster-config-file nodes-16381.conf
```

## 注意事项

1. 创建集群时，如果一个 Redis 实例中含有键值对，集群会创建失败；

```bash
[ERR] Node 127.0.0.1:6380 is not empty. Either the node already knows other nodes (check with CLUSTER NODES) or contains some key in database 0.
```



