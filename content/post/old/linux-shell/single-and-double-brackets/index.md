---
title: "Bash：单括号与双括号的区别"
date: 2023-03-23T22:47:16+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Shell
tags:
 - Linux Shell
 - Test Operators
 - 翻译
---

[原文：Differences Between Single and Double Brackets in Bash](https://www.baeldung.com/linux/bash-single-vs-double-brackets)

## 概述

当我们在 Bash 中做变量比较时，通常可以交换地使用单括号 `[]` 和双括号 `[[]]`。比如，我们可以使用表达式 `[ 3 -eq 3 ]` 或 `[[ 3 -eq 3 ]]` 来比较 `3` 是否等于 `3`。两个表达式都会执行成功，那两者的区别是什么呢？

在本文中，我们会讨论单括号和双括号之间的一些区别。

<!--more-->

## 主要区别

在本节中，我们会讨论单括号和双括号之间的主要区别。

### 单括号

在 Unix 和 Linux 中，`[` 是用于执行表达式的内置命令，我们可以使用 `type` 命令进行验证：

```bash
$ type [
[ is a shell builtin
```

`[` 是 `test` 内置命令的替代，两者可以交换使用：

```bash
$ [ 3 -eq 3 ] && echo "Numbers are equal"
Numbers are equal
$ test 3 -eq 3 && echo "Numbers are equal"
Numbers are equal
```

`[` 和 `test` 的唯一区别在于：使用 `[` 时，需要以 `]` 结尾，并且括号前后需要包含空格。

### 双括号

`[[]]` 是由 [Korn Shell](http://www.kornshell.com/) 第一次引入，其增强了 `[]` 在脚本中做比较、测试的功能，我们可以认为它是 `[]` 的增强版。

在 Bash 和 zsh 中，我们可以使用 `[[]]` 的方式，但***脚本可能不向后兼容 POSIX***。

`[[` 是 shell 的一个关键字，让我们再次使用 `type` 命令进行验证：

```bash
$ type [[
[[ is a shell keyword
```

## 其它区别

在本节中，我们会讨论单括号和双括号之间的其它区别。


### 比较操作符

在 `[[]]` 中可以使用比较操作符，如 `>`、`<` 等，如下：

```bash
$ [[ 1 < 2 ]] && echo "1 is less than 2"
1 is less than 2
```

上述命令中，我们使用 `<` 符号来检查 `1` 是否小于 `2`，命令是可以运行。但如果使用 `[]` 时，命令会报错：

```bash
$ [ 1 < 2 ] && echo "1 is less than 2"
bash: 2: No such file or directory
```

在这种情况下，Bash 会认为 `<` 是一个重定向符。因此，我们需要使用转义符 `\` 进行转义：

```bash
$ [ 1 \< 2 ] && echo "1 is less than 2"
1 is less than 2
```

现在，使用 `[]` 也可以执行成功。

类似的，对于 `>` 符号也需要使用转义符 `\` 进行转义。

对整数比较符 `-eq`、`-ne`、`-gt`、`-lt`、`-ge` 和 `-le`，`[]` 和 `[[]]` 都可以正常比较，如下：

```bash
$ [[ 1 -lt 2 ]] && echo "1 is less than 2"
1 is less than 2
$ [ 1 -lt 2 ] && echo "1 is less than 2"
1 is less than 2
```

### 布尔操作符

在 `[[]]` 中，我们可以使用逻辑运算 `&&` 和 `||`，如下：

```bash
$ [[ 3 -eq 3 && 4 -eq 4 ]] && echo "Numbers are equal"
Numbers are equal
```

但是在 `[]` 中，我们必须使用 `-a`、`-o` 分别代替 `&&` 和 `||`，如下：

```bash
$ [ 3 -eq 3 -a 4 -eq 4 ] && echo "Numbers are equal"
Numbers are equal
```

### 聚合表达式

在 `[[]]` 中，我们可以使用括号 `()` 来聚合多个表达式，`()` 的使用可以让脚本更具可读性，如下：

```bash
$ [[ 3 -eq 3 && (2 -eq 2 && 1 -eq 1) ]] && echo "Parentheses can be used"
Parentheses can be used
```

上述命令中，我们使用了 `()` 对表达式 `2 -eq 2 && 1 -eq 1` 做了聚合，然后将其作为第 2 个表达式参与到 `&&` 运算中，最后结果为验证为真。

但如果在 `[]` 中使用 `()`，则会报语法错误，如下：

```bash
$ [ 3 -eq 3 -a (2 -eq 2 -a 1 -eq 1) ] && echo "Parentheses can be used"
bash: syntax error near unexpected token `('
```

上述 `[]` 中，我们使用 `-a` 来代替 `&&` 但还是得到了一个报错。

在这里，我们必须使用转义符对括号进行转义，并且***前后要保留一个空格***，如下：

```bash
$ [ 3 -eq 3 -a \( 2 -eq 2 -a 1 -eq 1 \) ] && echo "Parentheses can be used"
Parentheses can be used
```

通过上面的改造，命令就可以成功执行。

### 模式匹配

在 `[[]]` 中，我们还可以使用通配符进行模式匹配，如下：

```bash
$ name=Alice
$ [[ $name = *c* ]] && echo "Name includes c"
Name includes c
$ echo $?
0
```

* `name=Alice`：完成变量赋值
* 使用 `$name = *c*` 检查变量中是否包含了 `c`
* `$?` 为 0 表示成功执行

如果把 `[[]]` 换成 `[]` 就无法成功执行，如下：

```bash
$ name=Alice
$ [ $name = *c* ] && echo "Name includes c"
bash: [: too many arguments
$ echo $?
2
```

在上述命令中，我们用 `[]` 代替了 `[[]]`，运行就会出错，这是因为 `[` 本身是内置的 Shell 命令，而该命令不支持这么多的参数。

### 正则表达式

正则表达式是另一种字符串模式匹配的方式，在 `[[]]` 中，我们可以使用正则表达式来做模式匹配的工作。

```bash
$ name=Alice
$ [[ $name =~ ^Ali ]] && echo "Regular expressions can be used"
Regular expressions can be used
```

首先，我们将值 "Alice" 赋值给 `name` 变量。然后，我们使用正则表达式来检查变量是否以 `Ali` 开头，其中 `=~` 操作符可以用于正则表达式的匹配、`^` 符号表示开头。最后，终端输出了第 2 个命令的执行结果。

那如果在 `[]` 中使用正则表达式呢？

```bash
$ name=Alice
$ [ $name =~ ^Ali ] && echo "Regular expressions can be used"
bash: [: =~: binary operator expected
```

我们得到了一个错误，所以得到结论：***在 `[]` 中不能使用正则表达式***。

### 单词分割

在 `[[]]` 中，Bash 不会对值中的单词进行分割，比如变量的值是一个包含空格的字符串，Bash 不会将其分割成多个单词。

```bash
$ filename="none existent file"
$ [[ ! -e $filename ]] && echo -n "File doesn't exist;" && echo $filename
File doesn't exist;none existent file
```

上面命令中，我们用一个不存在的文件作为示例，检查文件是否存在，但是值包含字符串。在 `[[]]` 中，我们可以直接使用 `filename` 变量，那在 `[]` 又如何？

```bash
$ filename="none existent file"
$ [ ! -e $filename ] && echo -n "File doesn't exist;" && echo $filename
[: too many arguments
```

上面命令中，我们得到了一个错误。如果希望在 `[]` 也能使用带空格的字符串，需要加上双引号，如下：

```bash
$ filename="none existent file"
$ [ ! -e "$filename" ] && echo -n "File doesn't exist;" && echo $filename
File doesn't exist;none existent file
```

## 结论

在本文中，我们讨论了单方括号和双方括号在 Bash 中的区别。

`[` 是内置的命令，并且其历史要久于 `[[`，作为后继者的 `[[]]` 是 `[]` 的增强版本。如果希望脚本更具兼容性，推荐使用 `[]`，如果希望脚本更具可读性，推荐使用 `[[]]`。

在示例中，我们看到在 `[]` 和 `[[]]` 中都可以使用比较操作符、布尔操作符以及聚合表达式，但也有一些区别。另外，正则表达式和单词分割只有能在 `[[]]` 中使用。

示例代码可在 [Github 地址](https://github.com/a2htray/code-notebook/tree/main/Shell/single_or_double_brackets) 中查看。