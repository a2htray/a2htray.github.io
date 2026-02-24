---
title: "集群相关"
date: 2022-04-10T11:34:31+08:00
draft: false
reward: false
toc: true
pinned: false
categories:
  - 后端技术
  - Redis
tags:
  - Redis
series:
  - "Redis 命令手册"
type: "docs"
---

Redis 服务器中与服务相关的命令，集群的配置过程可参考[《Redis 集群配置过程》](/posts/redis-cluster-deployment/)。

<!--more-->

## CLUSTER nodes：显示集群主从关系

```bash
127.0.0.1:6380> CLUSTER nodes
ac3f5aaae24bcd142b8303abedc5f57187ebc20e 127.0.0.1:6385@16385 slave 82ac1fe6c0af98d252ea96d4e84e7315eff31f8c 0 1649562125672 2 connected
9242c455757da4ad2f5aa818d75d12cde38231c2 127.0.0.1:6383@16383 slave 9a91b1be7a42e1bfe4b03deaa200c0e72fcf9b8e 0 1649562126000 3 connected
5cc9ed2602986aeffaf9997f3c38c675092a4810 127.0.0.1:6384@16384 slave 92a1fd50a1c8dc003e905ac828b0db64773d6b66 0 1649562126674 1 connected
9a91b1be7a42e1bfe4b03deaa200c0e72fcf9b8e 127.0.0.1:6382@16382 master - 0 1649562124670 3 connected 10923-16383
92a1fd50a1c8dc003e905ac828b0db64773d6b66 127.0.0.1:6380@16380 myself,master - 0 1649562122000 1 connected 0-5460
82ac1fe6c0af98d252ea96d4e84e7315eff31f8c 127.0.0.1:6381@16381 master - 0 1649562125000 2 connected 5461-10922
```

## CLUSTER slots：显示哈希槽分配信息

```bash
127.0.0.1:6380> CLUSTER slots
1) 1) (integer) 0
   2) (integer) 5460
   3) 1) "127.0.0.1"
      2) (integer) 6380
      3) "92a1fd50a1c8dc003e905ac828b0db64773d6b66"
   4) 1) "127.0.0.1"
      2) (integer) 6384
      3) "5cc9ed2602986aeffaf9997f3c38c675092a4810"
2) 1) (integer) 5461
   2) (integer) 10922
   3) 1) "127.0.0.1"
      2) (integer) 6381
      3) "82ac1fe6c0af98d252ea96d4e84e7315eff31f8c"
   4) 1) "127.0.0.1"
      2) (integer) 6385
      3) "ac3f5aaae24bcd142b8303abedc5f57187ebc20e"
3) 1) (integer) 10923
   2) (integer) 16383
   3) 1) "127.0.0.1"
      2) (integer) 6382
      3) "9a91b1be7a42e1bfe4b03deaa200c0e72fcf9b8e"
   4) 1) "127.0.0.1"
      2) (integer) 6383
      3) "9242c455757da4ad2f5aa818d75d12cde38231c2"
```

## CLUSTER keyslot：显示 key 对应的哈希槽

```bash
127.0.0.1:6380> CLUSTER keyslot user:1:name
(integer) 12440
```



