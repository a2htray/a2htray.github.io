+++
date = '2025-12-22T14:35:39+08:00'
draft = false
title = '图解负载均衡策略'
categories = ['后端技术', '基础原理']
tags = ['负载均衡']
toc = true
+++

看了一篇微信公众号文章，里面关于“负载均衡策略”的图例还挺好，一图胜千言，所以自己也画了下，加深下理解。

## 策略

![](/imgs/illustrated-load-balancing-strategies/load-balancer.png)

## 轮询

![](/imgs/illustrated-load-balancing-strategies/round-robin.png)

## 粘性轮询

![](/imgs/illustrated-load-balancing-strategies/stick-round-robin.png)

## 加权轮询

![](/imgs/illustrated-load-balancing-strategies/weighted-round-robin.png)

## 哈希

![](/imgs/illustrated-load-balancing-strategies/hash.png)

## 最少连接

![](/imgs/illustrated-load-balancing-strategies/least-connections.png)

## 最短响应

![](/imgs/illustrated-load-balancing-strategies/least-time.png)
