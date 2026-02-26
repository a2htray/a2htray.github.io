+++
date = '2026-02-26T17:02:42+08:00'
draft = false
title = 'Zod 的学习与使用：Record 数据结构（三）'
categories = ['前端技术', 'Zod']
tags = ['Zod', 'TypeScript']
toc = true
+++

- 学习资料：[https://zod.dev/api](https://zod.dev/api)
- 辅助：豆包
- Zod 版本：4

Zod 中定义的 Record 对应着 TS 中的 Record，用于构建键值对。

```ts
import * as z from 'zod'

const StudentSchema = z.record(z.string(), z.string())
type Student = z.infer<typeof StudentSchema>

// 等价于

type Student = Record<string, string>
```

## JSON 和对象

TS 中 Record 与 JSON 和对象相似，对应的概念如下：

* Record：编译时的类型约束（类型工具），是对 “键值对对象” 的类型描述
* JSON：文本格式（字符串），是数据传输 / 存储的 “序列化标准”
* 对象：运行时的数据结构（实体），是编程语言层面的 “键值对集合”

Record 和对象的键名规则支持字符串 / 数字 / 布尔 / Symbol，但都会转换成字符串，而 JSON 的键名规则，仅支持双引号包裹的字符串。

而 Zod 的 Record 是 TS Record 的运行时校验实现。

## 示例

Record 的使用示例。

**通过 z.union 约束键名规则**

```ts
// union 约束键名规则为字符串、数字或符号
const KeyTypeSchema = z.union([z.string(), z.number(), z.symbol()])
const ObjectSchema = z.record(KeyTypeSchema, z.any())
```

**通过 z.enum 约束键名取值**

```ts
// enum 约束键名取值为 "a"、"b" 或 "c"
const KeyValueSchema = z.enum(["a", "b", "c"])
const ObjectSchema = z.record(KeyValueSchema, z.any())
```

**支持键名的校验**

```ts
const KeySchema = z.union([
  z.string().min(4).max(10),
  z.number().int().min(0).max(100),
])
const ObjectSchema = z.record(KeySchema, z.any())

ObjectSchema.parse({
    "abcd": "123",
    10: "456",
})

// ObjectSchema.parse({
//     "a": "123",
//     10: "456",
// }) // 报错
```

## z.partialRecord

在使用 z.enum 作为 z.record 的第一个参数时，会穷尽检查所有的枚举值是否存在，而使用 z.partialRecord 则会进行部分选择。

```ts
const KeySchema = z.enum(["id", "name", "email"])
const ObjectSchema = z.partialRecord(KeySchema, z.string())

ObjectSchema.parse({
  id: "123",
  name: "张三",
})
ObjectSchema.parse({
  id: "123",
  name: "张三",
  email: "zhangsan@example.com",
})
```
