---
title: "如何使用 Python 复制文件"
date: 2023-06-15T20:44:17+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Python
tags:
 - 翻译
 - I/O
 - Python
---

[原文 How to Copy a File With Python](https://builtin.com/data-science/copy-a-file-with-python)

> Python 的 shutil 模块提供了 4 种复制文件的方法，根据你的实际情况选择合适的方法。或许，本文的内容可能帮助到你。

&emsp;&emsp;在每天的软件开发过程中，通过程序复制文件是一项平常的工作任务。我们将学习 Python `shutil` 模块提供的 4 种方法来完成文件复制，包括：

1. shutil.copy
2. shutil.copyfile
3. shutil.copy2
4. shutil.copyfileobj

&emsp;&emsp;shutil 模块是 Python 标准库的一部分，提供了许多高级别的文件操作方法。该模块提供了多种复制文件的方法，取决于你是否需要复制文件的元信息或权限。

&emsp;&emsp;本文的内容会覆盖上述所有方法，并在最后给出各方法特点的总结。

## 使用 shutil.copy 复制文件

&emsp;&emsp;`shutil.copy()` 方法可以复制一个<font color="red">不带元信息</font>的源到目标文件或目录并返回新创建后的路径，方法声明如下：

```python
shutil.copy(src, dst, *, follow_symlinks=True)
```

其中，`src` 可以是一个类路径（path-like）的对象或一个字符串。

> 该方法有以下特点：
> 1. 保留了文件的权限信息；
> 2. dst 参数可以是一个目录；
> 3. 不会复制元信息；
> 4. 参数不可以是 File 对象；

### shutil.copy 示例

```python
import shutil

# 复制 example.txt 并生成新文件 example_copy.txt

shutil.copy('example.txt', 'example_copy.txt')

# 复制 example.txt 到 test/ 目录下

shutil.copy('example.txt', 'test/')
```

## 使用 shutil.copyfile 复制文件

&emsp;&emsp;`shutil.copyfile()` 方法同样可以复制一个<font color="red">不带元信息</font>的源到目标文件，方法声明如下：

```python
shutil.copyfile(src, dst, *, follow_symlinks=True)
```

其中，`src` 可以是一个类路径（path-like）的对象或一个字符串。

> 该方法有以下特点：
> 1. 不保留文件的权限信息；
> 2. dst 参数不可以是目录；
> 3. 不会复制元信息；
> 4. 参数不可以是 File 对象；

### shutil.copyfile 示例

```python
import shutil

# 复制 source.txt 并生成新文件 destination.txt

shutil.copyfile('source.txt', 'destination.txt')
```

## 使用 shutil.copy2

&emsp;&emsp;`shutil.copy2()` 有别于 `shutil.copy()`，该方法会保留文件的元信息，方法声明如下：

```python
shutil.copy2(src, dst, *, follow_symlinks=True)
```

> 该方法有以下特点：
> 1. 保留文件的权限信息；
> 2. dst 参数可以是目录；
> 3. 可以复制元信息；
> 4. 参数不可以是 File 对象；

### shutil.copy2 示例

```python
import shutil

# 复制 example.txt 并生成新文件 example_copy.txt

shutil.copy2('example.txt', 'example_copy.txt')

# 复制 example.txt 到 test/ 目录下

shutil.copy2('example.txt', 'test/')
```

## 使用 shutil.copyfileobj

&emsp;&emsp;如果你需要使用 File 对象进行传参，那么你可以使用 `shutil.copyfileobj()` 方法。该方法会将源文件的内容复制到目标文件，同时使用 `length` 参数来控制复制内容的大小，方法声明如下：

```python
shutil.copyfileobj(fsrc, fdst[, length])
```

> 该方法有以下特点：
> 1. 不保留文件的权限信息；
> 2. dst 参数不可以是目录；
> 3. 不会复制元信息；
> 4. 参数可以是 File 对象；

### shutil.copyfileobj 示例

```python
import shutil

source_file = open('example.txt', 'rb')

dest_file = open('example_copy.txt', 'wb')

shutil.copyfileobj(source_file, dest_file)
```

## 如何选择合适的复制文件的方法

&emsp;&emsp;本文中我们学习了 4 种由 Python shutil 标准包提供的复制文件的方法，选择合适的方法取决于特定的使用场景，如：

1. 是否需要复制文件权限或元信息；
2. 是否需要复制到目录；
3. 传参是否是 File 对象；

&emsp;&emsp;下表示总结了各方法的特点：

| 方法                 | 保留文件权限 | 目标可以是目录 | 复制元信息 |
| ------------------ | ------ | ------- | ----- |
| shutil.copy        | Y      | Y       | N     |
| shutil.copyfile    | N      | N       | N     |
| shutil.copy2       | Y      | Y       | Y     |
| shutil.copyfileobj | N      | N       | N     |