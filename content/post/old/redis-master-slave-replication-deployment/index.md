---
title: "Redis 主从复制配置过程"
date: 2022-03-30T10:28:31+08:00
draft: false
reward: false
toc: true
pinned: false
categories:
  - 后端技术
  - Redis
tags:
  - Redis
  - 主从复制
---

Redis 主从复制可以实现数据库的读写分离，即主节点负责接收写请求、从节点负责接收读请求，是高性能 Redis 服务的基础。所以配置 Redis 主从复制应当作为开发者的技能之一，后文内容包括：

1. 单机配置一主二从的主从复制服务
2. 服务验证；

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

主节点使用 `6379` 端口，两个从节点分别使用 `6380` 和 `6381` 端口。

### Redis 配置文件

复制两份 Redis 配置文件分别为两个从节点的配置文件：

```bash
cp -a /etc/redis/redis.conf /etc/redis/redis-server-6380.conf
cp -a /etc/redis/redis.conf /etc/redis/redis-server-6381.conf
```

修改两个配置文件的内容，修改及新增内容如下：

```bash
# redis-6380.conf
# 修改项
port 6380
pidfile /run/redis/redis-server-6380.pid
logfile /var/log/redis/redis-server-6380.log
dbfilename dump-6380.rdb
# 新增项
slaveof 127.0.0.1 6379
```

```bash
# redis-6381.conf
# 修改项
port 6381
pidfile /run/redis/redis-server-6381.pid
logfile /var/log/redis/redis-server-6381.log
dbfilename dump-6381.rdb
# 新增项
slaveof 127.0.0.1 6379
```

### systemd 配置文件

复制两份 systemd 配置文件分别作为两个从节点的服务启动文件：

```bash
cp -a /lib/systemd/system/redis-server.service  /lib/systemd/system/redis-server-6380.service
cp -a /lib/systemd/system/redis-server.service  /lib/systemd/system/redis-server-6381.service
```

修改两个配置文件的内容，均为修改项，如下：

```bash
# redis-server-6380.service
[Service]
ExecStart=/usr/bin/redis-server /etc/redis/redis-server-6380.conf
PIDFile=/run/redis/redis-server-6380.pid
```

```bash
# redis-server-6381.service
[Service]
ExecStart=/usr/bin/redis-server /etc/redis/redis-server-6381.conf
PIDFile=/run/redis/redis-server-6381.pid
```

修改之后，需要执行如下命令进行重新加载：

```bash
systemctl daemon-reload
```



## 启动服务

配置完成后，通过 `systemd` 的管理命令分别在 3 个终端各启动 1 个服务，命令及显示如下：

```bash
# 主节点
$ systemctl start redis-server.service
$ redis-cli -p 6379
127.0.0.1:6379> INFO Replication
# Replication
role:master
connected_slaves:2
slave0:ip=127.0.0.1,port=6380,state=online,offset=308,lag=0
slave1:ip=127.0.0.1,port=6381,state=online,offset=308,lag=1
# 其它 ...
```

```bash
# 从节点 6380
$ systemctl start redis-server-6380.service
$ redis-cli -p 6380
127.0.0.1:6380> INFO Replication
# Replication
role:slave
master_host:127.0.0.1
master_port:6379
# 其它 ...
```

```bash
# 从节点 6381
$ systemctl start redis-server-6381.service
$ redis-cli -p 6381
127.0.0.1:6381> INFO Replication
# Replication
role:slave
master_host:127.0.0.1
master_port:6379
# 其它 ...
```

### 读写 Redis 数据库

是否发生主从复制，可按如下的命令依次执行进行验证。

```bash
# 主节点
127.0.0.1:6379> SET topic master-slave-replication
OK
# 从节点 6380
127.0.0.1:6380> GET topic
"master-slave-replication"
# 从节点 6381
127.0.0.1:6381> GET topic
"master-slave-replication"
```

### 其它

`systemd` 是 Linux 服务器管理服务的其中一种方式，服务的启动与关闭也可以通过 `redis-server` 命令或其它方式进行实现。Redis 使用到的目录及文件信息包括：

* `/run/redis/`：存放 Redis 服务的 pid 文件，由 `pidfile` 配置项决定；
* ` /var/log/redis/`：存放 Redis 服务的日志文件，由 `logfile` 配置项决定；
* `/var/lib/redis`：存放 RDB 文件，由 `dir` 和 `dbfilename` 配置项决定；

## 完整脚本

```bash
#!/bin/bash

SLAVE1_PORT=6380
SLAVE2_PORT=6381

cat /etc/redis/redis.conf | \
sed "s/^port\ 6379/port\ $SLAVE1_PORT/g" | \
sed "s/^pidfile \/run\/redis\/redis-server\.pid/pidfile \/run\/redis\/redis-server-$SLAVE1_PORT\.pid/g" | \
sed "s/^logfile \/var\/log\/redis\/redis-server\.log/logfile \/var\/log\/redis\/redis-server-$SLAVE1_PORT\.log/g" | \
sed "s/^dbfilename dump\.rdb/dbfilename dump-$SLAVE1_PORT\.rdb/g" | \
sed "\$a\\slaveof 127.0.0.1 6379\\" > /etc/redis/redis-server-$SLAVE1_PORT.conf
chown redis.redis /etc/redis/redis-server-$SLAVE1_PORT.conf

cat /etc/redis/redis.conf | \
sed "s/^port\ 6379/port\ $SLAVE2_PORT/g" | \
sed "s/^pidfile \/run\/redis\/redis-server\.pid/pidfile \/run\/redis\/redis-server-$SLAVE2_PORT\.pid/g" | \
sed "s/^logfile \/var\/log\/redis\/redis-server\.log/logfile \/var\/log\/redis\/redis-server-$SLAVE2_PORT\.log/g" | \
sed "s/^dbfilename dump\.rdb/dbfilename dump-$SLAVE2_PORT\.rdb/g" | \
sed "\$a\\slaveof 127.0.0.1 6379\\" > /etc/redis/redis-server-$SLAVE2_PORT.conf
chown redis.redis /etc/redis/redis-server-$SLAVE2_PORT.conf

cat /lib/systemd/system/redis-server.service | \
sed "s/\/etc\/redis\/redis\.conf/\/etc\/redis\/redis-server-$SLAVE1_PORT\.conf/g" | \
sed "s/\/run\/redis\/redis-server\.pid/\/run\/redis\/redis-server-$SLAVE1_PORT\.pid/g" > /lib/systemd/system/redis-server-$SLAVE1_PORT.service
chown root.root /lib/systemd/system/redis-server-$SLAVE1_PORT.service

cat /lib/systemd/system/redis-server.service | \
sed "s/\/etc\/redis\/redis\.conf/\/etc\/redis\/redis-server-$SLAVE2_PORT\.conf/g" | \
sed "s/\/run\/redis\/redis-server\.pid/\/run\/redis\/redis-server-$SLAVE2_PORT\.pid/g" > /lib/systemd/system/redis-server-$SLAVE2_PORT.service
chown root.root /lib/systemd/system/redis-server-$SLAVE2_PORT.service

systemctl daemon-reload

systemctl start redis-server.service
systemctl start redis-server-$SLAVE1_PORT.service
systemctl start redis-server-$SLAVE2_PORT.service
```

## 总结

1. Redis 主从复制的配置过程；
2. Redis 服务相关配置项说明；