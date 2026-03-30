+++
date = '2026-03-30T06:56:40+08:00'
draft = false
title = 'JSON Web Token - 介绍'
categories = ['后端技术', 'JWT']
tags = ['JWT', 'JSON Web Token', 'The JWT Handbook']
toc = true
+++

最近在写一个产品，在认证鉴权模块用到了 JWT，为探究其本质，所以打算花一些时间学习一下，并记录成文。

首先，回答两个问答：

1. 什么是 JWT？
2. JWT 解决了什么样的问题？

_什么是 JWT？_

JWT，全称 JSON Web Token，是一个**在空间受限环境下传递鉴权声明**的标准，目前已广泛地应用于大部分主流 Web 框架中。

JWT 提供了多种不同类型的签名与加密方式，其中**签名**为了验证数据是否被篡改，而**加密**则是为了保持数据内容。

> JWTs, by virtue of JWS and JWE, can provide various types of signatures and encryption. Signatures are useful to
> validate the data against tampering. Encryption is useful to protect the data from being read by third parties.

_JWT 解决了什么样的问题？_

JWT 统一标准化解决了**多方安全可信传递用户身份 / 权限声明**的通用难题，替代杂乱私有方案，适配无状态、跨域、分布式场景。

<p style="text-align: center; color: #7c7c7c;">------------------便于阅读，我是分割线---------------------</p>

然后，看一个已签名的的 JWT：

eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ

由 3 部分组成，分别是头部（Header）、载荷（Payload）和签名（Signature），各自的作用如下：

* 头部：描述令牌的类型、签名算法
* 载荷：存储用户数据
* 签名：防篡改

3 部分以 “.” 分隔。

<p style="text-align: center; color: #7c7c7c;">------------------便于阅读，我是分割线---------------------</p>

最后，了解一下 JWT 的使用场景：

1. 客户端无状态会话：客户端存储部分用户信息
2. 联合身份系统：允许用户使用一种身份验证访问不同的系统
