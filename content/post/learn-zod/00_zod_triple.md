+++
date = '2026-02-24T23:23:06+08:00'
draft = false
title = 'Zod 的学习与使用：Zod 基础知识（一）'
categories = ['前端技术', 'Zod']
tags = ['Zod', 'TypeScript']
toc = true
+++

- 学习资料：[https://zod.dev/](https://zod.dev/)
- 辅助：豆包
- Zod 版本：4

学习一个新的技术、框架或包，首先要解决的是事物的三个问题：

* 是什么
* 做什么
* 有什么

本文就是围绕 Zod 回答上面的三个问题，并在后续的文章中介绍“怎么用”。

## 是什么

> TypeScript-first schema validation with static type inference.

一言以蔽之，<span style="font-weight: bold; background-color: yellow;">Zod 是具有静态类型推导能力，以 TypeScript 优先的模式校验包</span>。

Zod 的出现，解决了“<span style="font-weight: bold; background-color: yellow;">TS 类型编译后消息、无法校验运行时数据</span>”的问题，由此可见，Zod 是将类型安全（模式校验）从编译时延伸到运行时的工具。

编译时的类型安全由 TS 提供，如果使用的 JS，由 Zod 进行提供，这也就是上面“静态类型推导能力”的展现。而运行时的模式校验，则完全由 Zod 提供。

## 做什么

模式校验，即是对数据结构的校验，针对这一点，Zod 提供以下的能力：

* 校验数据类型与值规则
  * 基础类型校验
  * 复合类型校验
  * 格式校验与逻辑判断
* 自动联动 TypeScript 类型：使用 `z.infer`
* 自定义校验规则
* 容错式校验与错误处理：`parse` -> 严格模式，`safeParse` -> 容错模式
* 数据转换：检验的同时可对数据进行转换

上述能力可适用于以下场景：

* 前端表单校验 
* API 数据校验
* 配置文件校验
* 跨端/跨服务数据校验
* 数据库操作校验

总的来说，Zod 能做的事可归纳于 3 点：

1. 校验：校验任意来源数据的类型、格式、值规则，覆盖基础类型到复杂结构
2. 联动：自动推导 TypeScript 类型，实现 “一份定义，同时满足编译时类型检查和运行时数据校验”
3. 适配：支持自定义规则、容错处理、数据转换，能无缝适配前端表单、API 校验、配置校验等各类实战场景

## 有什么

Zod 的目的是验证数据结构，那么 Zod 就必然包含<span style="font-weight: bold; background-color: yellow;">数据类型定义</span>，有：

* 基础数据类型（Primitive Type）
* 胁迫数据类型（Coercion Type）

> 胁迫，代表一定的强制，即值会做一次强制转换。

基础数据类型包括：

```javascript
import * as z from "zod"

z.string()
z.number()
z.bigint()
z.boolean()
z.symbol()
z.undefined()
z.null()
// ...
```

对应的胁迫数据类型，则 `z.coerce`，如：

```javascript
z.coerce.string()
z.coerce.number()
z.coerce.boolean()
z.coerce.bigint()
// ...
```

针对字符串类型，还提供字符串格式校验的数据类型，如下：

```javascript
z.email()
z.uuid()
z.url()
z.httpUrl()
z.hostname()
z.emoji()
z.base64()
z.base64url()
z.hex()
z.jwt()
z.nanoid()
z.cuid()
z.cuid2()
z.ulid()
z.ipv4()
z.ipv6()
z.mac()
z.cidrv4()
z.cidrv6()
z.hash("sha256")
z.iso.date()
z.iso.time()
z.iso.datetime()
z.iso.duration()
```

同时支持自定义字符串格式。

数据结构是为了 Schema 的定义，那么就会有校验失败的时候，所以 Zod 还有“校验错误”，支持自定义错误。

其它还有以下概念：

* Metadata：描述数据结构的数据，通过 Registry 声明
* Registry：一系列数据结构的集合
* JSON Schema：JSON 格式的模式描述
* 编码与解码：用于数据类型的转换

## 总结

Zod 是用于数据结构校验的 TS/JS 包，可以在不同场景下做各种数据的校验与转换，其包含不同的数据类型（基础、胁迫）、字符串格式、Metadata、Registry、JSON Schema、编码与解码。

