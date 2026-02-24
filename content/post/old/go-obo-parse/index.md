---
title: "解析 OBO 文件"
date: 2022-04-29T15:31:20+08:00
draft: true
reward: false
categories:
  - 后端技术
  - Go
tags:
 - OBO
 - 生物信息
images:
 - images/go-obo-file.jpg
---

OBO 文件可以用于描述实体（Ontology）以及语义词汇（Semantic Vocabulary），常用于定义分类信息，其格式非常简单，**参考 [1]** 为 1.4 版本的编写规范。也正因为其简单性，许多程序都在内部实现解析，所以没有太多与之相关的 Go 包。本文主要用于记录开发 `github.com/a2htray/oboparser` 包的过程，也使自己更好地理解 OBO 格式。

<!--more-->

## OBO 格式

`so.obo` 可用于描述序列的信息，取 `so.obo` 的片段，如下：

```bash
format-version: 1.2
data-version: so-xp/releases/2015-11-24/so-xp.owl
date: 30:10:2019 11:24
saved-by: david
auto-generated-by: OBO-Edit 2.3.1
! ...

[Term]
id: SO:0000000
name: Sequence_Ontology
subset: SOFA
is_obsolete: true

! ...

[Typedef]
id: variant_of
name: variant_of
def: "A' is a variant (mutation) of A = definition every instance of A' is either an immediate mutation of some instance of A, or there is a chain of immediate mutation processes linking A' to some instance of A." [SO:immuno_workshop]
comment: Added to SO during the immunology workshop, June 2007.  This relationship was approved by Barry Smith.
```

可见，OBO 文件由两部分组成：

1. 头信息：key-value 形式，key 值存在重名；
2. Stanza 信息：用 `[]` 表示一个 stanza，其中 `[]` 内的值可以为 `Term`、`Typedef`、`Instance`；



## 参考

1. [The OBO Flat File Format Guide, version 1.4](https://owlcollab.github.io/oboformat/doc/GO.format.obo-1_4.html)
2. [软件官网 OBO Editor](http://oboedit.org/)
3. [The Open Biological and Biomedical Ontology (OBO) Foundry](https://obofoundry.org/)
4. [Gene Ontology](http://geneontology.org/docs/download-ontology/)

