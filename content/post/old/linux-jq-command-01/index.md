---
title: "命令 jq（一）基础过滤器"
date: 2023-02-22T22:13:29+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Shell
tags:
 - Linux Shell
 - jq
---

jq 是 Linux 下处理 JSON 文档字符串的命令行工具，可用于过滤并格式化输出特定的内容。其官网的手册详尽地介绍了 jq 的使用方法，内容相当丰富。本文则注重对过滤器的使用，并通过不同的示例加以说明。

<!--more-->

## Identity `.`

`.` 是最简单的过滤器，该过滤器接收输入并且不做修改地进行输出。如果需要对 JSON 字符串进行格式化输出， 则可以使用该过滤器。

```bash
$ echo '{"title": "大数据分析", "price": 58}' | jq '.'
{
  "title": "大数据分析",
  "price": 58
}
$ cat unformatted.json | jq '.'
{
  "id": 1,
  "name": "myName",
  "position": "5"
}
```

## Object Identifier-Index `.foo`, `.foo.bar`

`.foo` 可以获取输入中键为 `foo` 的值，同理，`.foo.bar` 则是先获取键为 `foo` 的值并在此基础上获取键为 `bar` 的值。

```bash
$ echo '{"name": "app-01", "kind": "Pod"}' | jq '.name'
"app-01"
$ echo '{"fields": {"name": {"type": "string"}}}' | jq '.fields.name'
{
  "type": "string"
}
```

## Optional Object Identifier-Index `.foo?`

`.foo?` 与 `.foo` 作用相同，但当输入不是一个合规的数组或对象时不会报错。

```bash
$ echo '{"name": "a2htray"}' | jq '.name.errorField'
jq: error (at <stdin>:1): Cannot index string with string "errorField"
$ echo '{"name": "a2htray"}' | jq '.name.errorField?'  # 不会报错
```

## Generic Object Index `.[<string>]`

`<string>` 表示一个对象的键名，且 `.[<string>]` 与 `.<string>` 的效果相同。

```bash
$ echo '{"fields": {"name": {"type": "string"}}}' | jq '.["fields"]'
{
  "name": {
    "type": "string"
  }
}
$ echo '{"fields": {"name": {"type": "string"}}}' | jq '.["fields"]["name"]'
{
  "type": "string"
}
$ echo '{"fields": {"name": {"type": "string"}}}' | jq '.fields.name'
{
  "type": "string"
}
$ cat example.json | jq '.properties["name", "email"]'
{
  "type": "string"
}
{
  "type": "string",
  "format": "email"
}
```

## Array Index `.[integer]`

`.[integer]` 用于检索数组对象，并且索引值支持负数。

```bash
$ echo '[1, 2, 3, 4]' | jq '.[0]'
1
$ echo '[1, 2, 3, 4]' | jq '.[-1]'
4
```

## Array/String Index `.[integer1:integer2]`

`.[integer1:integer2]` 对数组或字符串进行切片操作，返回原数据的子集。

```bash
$ echo '[1, 2, 3, 4]' | jq '.[0:2]'
[
  1,
  2
]
$ echo '"1234"' | jq '.[0:2]'
"12"
```

## Array/Object Value Iterator `.[]`

`.[]` 返回数组的所有元素或对象的所有的值。

```bash
$ echo '[1, 2, 3, 4]' | jq '.[]'
1
2
3
4
$ echo '{"numbers": [1, 2, 3, 4]}' | jq '.numbers[]'
1
2
3
4
$ echo '{"name": "a2htray", "like": 2}' | jq '.[]'
"a2htray"
2
```

## Optional Array/Object Value Iterator `.[]?`

`.[]?` 相对于 `.[]`，就好像 `.foo?` 相对于 `.foo`。

```bash
$ echo '{"name": "a2htray", "like": 2}' | jq '.unlike[]'
jq: error (at <stdin>:1): Cannot iterate over null (null)
$ echo '{"name": "a2htray", "like": 2}' | jq '.unlike[]?'  # 不会报错
```

## 小结

jq 的过滤器内容不算多，使用的难点在于：面对不同 schema 的 JSON 结构，如何使用和串联过滤器。针对 JS 中的定义的数据类型，jq 都能提供相应的过滤器来获取特定的信息，包括：

1. object: 按键名取值
2. array: 索引取值、切片
3. string: 切片

本文使用的文件及命令见：[https://github.com/a2htray/code-notebook/tree/main/Shell/jq](https://github.com/a2htray/code-notebook/tree/main/Shell/jq)。