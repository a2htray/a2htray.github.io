---
title: "Redis 集群配置过程"
date: 2022-04-07T14:07:35+08:00
draft: false
reward: false
toc: true
pinned: false
categories:
  - 后端技术
  - Redis
tags:
  - Redis
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

### 配置文件

3 主使用的端口分别为 6380、6381 和 6382，3 从使用的端口为 6383、6384 和 6385。创建配置文件以及工作目录：

```bash
$ mkdir /etc/redis/cluster
$ chown redis.redis /etc/redis/cluster
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6380.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6381.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6382.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6383.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6384.conf
$ cp -a /etc/redis/redis.conf /etc/redis/cluster/redis-6385.conf
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
# /etc/redis/cluster/redis-6383.conf
port 6383
cluster-enabled yes
pidfile /run/redis/redis-server-6383.pid
logfile /var/log/redis/redis-server-6383.log
dbfilename dump-6383.rdb
cluster-config-file nodes-6383.conf
```

```bash
# /etc/redis/cluster/redis-6384.conf
port 6384
cluster-enabled yes
pidfile /run/redis/redis-server-6384.pid
logfile /var/log/redis/redis-server-6384.log
dbfilename dump-16381.rdb
cluster-config-file nodes-6384.conf
```

### 启动脚本

编写 `redis-cluster.sh` 脚本，记得 `chmod a+x redis-cluster.sh` 一下：

```bash
#!/bin/bash

if [ "$1" == "start" ]; then
  redis-server /etc/redis/cluster/redis-6380.conf
  redis-server /etc/redis/cluster/redis-6381.conf
  redis-server /etc/redis/cluster/redis-6382.conf
  redis-server /etc/redis/cluster/redis-6383.conf
  redis-server /etc/redis/cluster/redis-6384.conf
  redis-server /etc/redis/cluster/redis-6385.conf
elif [ "$1" == "stop" ]; then
  redis-cli -p 6380 shutdown
  redis-cli -p 6381 shutdown
  redis-cli -p 6382 shutdown
  redis-cli -p 6383 shutdown
  redis-cli -p 6384 shutdown
  redis-cli -p 6385 shutdown
else
  echo "please, pass an action [start|stop]"
fi
```

启动所有 Redis 节点：

```bash
$ ./redis-cluster.sh start
```

验证集群节点是否启动：

```bash
 ps -ef | grep redis-server
root     4025397       1  0 11:10 ?        00:00:00 redis-server *:6380 [cluster]
root     4025405       1  0 11:10 ?        00:00:00 redis-server *:6381 [cluster]
root     4025411       1  0 11:10 ?        00:00:00 redis-server *:6382 [cluster]
root     4025417       1  0 11:10 ?        00:00:00 redis-server *:6383 [cluster]
root     4025423       1  0 11:10 ?        00:00:00 redis-server *:6384 [cluster]
root     4025429       1  0 11:10 ?        00:00:00 redis-server *:6385 [cluster]
```

### 配置集群

使用 `redis-cli` 命令配置集群：

```bash
$ redis-cli --cluster create 127.0.0.1:6380 127.0.0.1:6381 127.0.0.1:6382 127.0.0.1:6383  127.0.0.1:6384 127.0.0.1:6385 --cluster-replicas 1
>>> Performing hash slots allocation on 6 nodes...
Master[0] -> Slots 0 - 5460
Master[1] -> Slots 5461 - 10922
Master[2] -> Slots 10923 - 16383
Adding replica 127.0.0.1:6384 to 127.0.0.1:6380
Adding replica 127.0.0.1:6385 to 127.0.0.1:6381
Adding replica 127.0.0.1:6383 to 127.0.0.1:6382
>>> Trying to optimize slaves allocation for anti-affinity
[WARNING] Some slaves are in the same host as their master
M: 92a1fd50a1c8dc003e905ac828b0db64773d6b66 127.0.0.1:6380
   slots:[0-5460] (5461 slots) master
M: 82ac1fe6c0af98d252ea96d4e84e7315eff31f8c 127.0.0.1:6381
   slots:[5461-10922] (5462 slots) master
M: 9a91b1be7a42e1bfe4b03deaa200c0e72fcf9b8e 127.0.0.1:6382
   slots:[10923-16383] (5461 slots) master
S: 9242c455757da4ad2f5aa818d75d12cde38231c2 127.0.0.1:6383
   replicates 9a91b1be7a42e1bfe4b03deaa200c0e72fcf9b8e
S: 5cc9ed2602986aeffaf9997f3c38c675092a4810 127.0.0.1:6384
   replicates 92a1fd50a1c8dc003e905ac828b0db64773d6b66
S: ac3f5aaae24bcd142b8303abedc5f57187ebc20e 127.0.0.1:6385
   replicates 82ac1fe6c0af98d252ea96d4e84e7315eff31f8c
Can I set the above configuration? (type 'yes' to accept): yes
>>> Nodes configuration updated
>>> Assign a different config epoch to each node
>>> Sending CLUSTER MEET messages to join the cluster
Waiting for the cluster to join
.
>>> Performing Cluster Check (using node 127.0.0.1:6380)
M: 92a1fd50a1c8dc003e905ac828b0db64773d6b66 127.0.0.1:6380
   slots:[0-5460] (5461 slots) master
   1 additional replica(s)
S: ac3f5aaae24bcd142b8303abedc5f57187ebc20e 127.0.0.1:6385
   slots: (0 slots) slave
   replicates 82ac1fe6c0af98d252ea96d4e84e7315eff31f8c
S: 9242c455757da4ad2f5aa818d75d12cde38231c2 127.0.0.1:6383
   slots: (0 slots) slave
   replicates 9a91b1be7a42e1bfe4b03deaa200c0e72fcf9b8e
S: 5cc9ed2602986aeffaf9997f3c38c675092a4810 127.0.0.1:6384
   slots: (0 slots) slave
   replicates 92a1fd50a1c8dc003e905ac828b0db64773d6b66
M: 9a91b1be7a42e1bfe4b03deaa200c0e72fcf9b8e 127.0.0.1:6382
   slots:[10923-16383] (5461 slots) master
   1 additional replica(s)
M: 82ac1fe6c0af98d252ea96d4e84e7315eff31f8c 127.0.0.1:6381
   slots:[5461-10922] (5462 slots) master
   1 additional replica(s)
[OK] All nodes agree about slots configuration.
>>> Check for open slots...
>>> Check slots coverage...
[OK] All 16384 slots covered.
```

输出信息为集群的哈希槽的分配信息及主从关系，可知：

1. 共有 16384（以 0 开始计数），将 0-5460 的哈希槽分配给主节点 1，将 5461-10922 的哈希槽分配给主节点 2，将 10923-16383 的哈希槽分配给主节点 3；
2. 127.0.0.1:6384 是 127.0.0.1:6380 的备份，127.0.0.1:6385 是 127.0.0.1:6381 的备份，127.0.0.1:6383 是 127.0.0.1:6382 的备份；

### 管理工具

[RedisInsight](https://github.com/redisinsight/redisinsight) 是 Redis 官方提供的 Redis 服务器图形化管理工具，操作性很强，包含：

1. 数据维护；
2. 性能测试；
3. 可视化；
4. 支持不同类型的 Redis 服务；

## 注意事项

1. 创建集群时，如果一个 Redis 实例中含有键值对，集群会创建失败；

```bash
[ERR] Node 127.0.0.1:6380 is not empty. Either the node already knows other nodes (check with CLUSTER NODES) or contains some key in database 0.
```

2. 集群模式中，各 Redis 节点要使用两个端口，分别为指定的 port 和 `port+10000`；



