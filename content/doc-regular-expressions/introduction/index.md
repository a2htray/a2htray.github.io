---
title: "介绍"
date: 2022-04-29T23:45:33+08:00
draft: false
toc: true
reward: false
series:
  - "正则表达式"
type: "docs"
categories:
 - 正则表达式
---

正则表达式教程 - 介绍。

<!--more-->

## 学习如何使用和充分利用正则表达式

这篇教程会教会你所有关于如何创造高效正则表达式的方法。教程从最基础的概念开始，所以即使你不懂正则表达式也没有关系，你也可能跟上教程的节奏。本教程会讲解正则表达式引擎的工作机制，理解正则表达式引擎有助于对一些疑难杂症的问题快速寻找解决方案。同时，本教程也旨在节约你学习正则表达式的时间，带你更好地入门。

## 什么是正则表达式

一般来说，正则表达式是指用于描述特写长度文本的**模式字符串**，正则的名字来源于数学方法，但本教程并不会深入数学。正则一般可以缩写为 regex 和 regexp，本教程使用 regex（PS：翻译过程中使用 regex 表示正则）。在书写的过程中，使用 `regex` 表示一个正则表达式字符串。

第 1 个示例 `regex` 是一个完全合法的正则表达式

This first example is actually a perfectly valid regex. It is the most basic pattern, simply matching the literal text `regex`. A “match” is the piece of text, or sequence of bytes or characters that pattern was found to correspond to by the regex processing software. Matches are highlighted in blue on this site.

`\b[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}\b` is a more complex pattern. It describes a series of letters, digits, dots, underscores, percentage signs and hyphens, followed by an at sign, followed by another series of letters, digits and hyphens, finally followed by a single dot and two or more letters. In other words: this pattern describes an [email address](https://www.regular-expressions.info/email.html). This also shows the syntax highlighting applied to regular expressions on this site. Word boundaries and quantifiers are blue, character classes are orange, and escaped literals are gray. You’ll see additional colors like green for grouping and purple for meta tokens later in the tutorial.

With the above regular expression pattern, you can search through a text file to find email addresses, or verify if a given string looks like an email address. This tutorial uses the term “string” to indicate the text that the regular expression is applied to. This website highlights them in `green`. The term “string” or “character string” is used by programmers to indicate a sequence of characters. In practice, you can use regular expressions with whatever data you can access using the application or programming language you are working with.