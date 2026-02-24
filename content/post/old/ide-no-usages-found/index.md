---
title: "IDE 代码显示无引用 no usages found"
date: 2022-04-27T20:55:01+08:00
draft: false
reward: false
categories:
 - 生产力工具
 - JetBrains
tags:
 - JetBrains
---

IntelliJ IDEA 代码显示灰色，表示无任何引用，实际上是有引用。出现这种问题，非常不易于 DEBUG。

<!--more--> 

解决方法：File -> Invalidate Caches ...

![](QQ截图20220427210450.png)

勾选：

1. Clear file system cache and Local History
2. Ask before downloading new shared indexes

![](QQ截图20220429133118.png)