---
title: "命令 jq（二） 特殊字符和构造器"
date: 2023-02-22T22:40:49+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Shell
tags:
 - Linux Shell
 - jq
---

上一篇文章介绍了 jq 中的过滤器，但单独地使用过滤器无法满足复杂的实际需求，所以 jq 中引入了特殊字符（`,` 和 `|`）和构造器来实现定制化的输出。

<!--more-->

## 特殊字符

jq 中有两个特殊字符，`,` 和 `|`，分别用于实现：

- `,` 多过滤器处理相同输入
- `|` 类似于 Linux 管道符操作

### `,`

`,` 用于分隔多个过滤器，多个过滤器可同时对输入起作用并作为连续输出，如：

```bash
$ cat example.json | jq '.type, .required'
"object"
[
  "email"
]
```

### `|`

`|` 的工作模式就是依照 Linux 中的管道符的处理方式，前一个过滤器的输出会作为后一个过滤器的输入，如：

```bash
$ cat example.json | jq '.type | .[0:3]'
"obj"
```

上述示例中，`.type` 作用是获取 JSON 对象的 `type` 值，之后将该值作为 `.[0:3]` 的输入，而 `type` 值为一个字符串，且 `.[0:3]` 作为于字符串时会实现获取字符子串的效果，所以输出值为 `obj`。

## 构造器

构造器可以将输出以某种特定的类型进行展示，包括数组构造器 `[]` 和对象构造器 `{}`。

### `[]`

在第 1 个示例中，其输出并没有包含特定的结构，所以可以使用 `[]` 将其以数组的形式进行展示，如：

```bash
$ cat example.json | jq '[.type, .required]'
[
  "object",
  [
    "email"
  ]
]
```

### `{}`

与 `[]` 类似，只不过输出的形式为对象。

```bash
$ cat example.json | jq '{newKeyOfType: .type, newTypeOfRequired: .required}'
{
  "newKeyOfType": "object",
  "newTypeOfRequired": [
    "email"
  ]
}
```

## 小结

特殊字符和构造器是复杂过滤器的基础，包含特殊字符和构造器的表达式也较为难懂，需要一步步拆分开进行分析。

本文使用的文件及命令见：[https://github.com/a2htray/code-notebook/tree/main/Shell/jq](https://github.com/a2htray/code-notebook/tree/main/Shell/jq)。