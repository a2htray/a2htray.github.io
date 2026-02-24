---
title: "Python 读写文件"
date: 2022-03-28T13:28:47+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Python
tags:
 - file
---

在开发过程中，开发者常常需要对文件执行读写操作，仅以此文记录读写文件的常规用法。

<!--more-->

## 打开和关闭文件

Python 的内建函数 `open` 可以打开一个文件，可返回一个文件对象 `TextIOWrapper`（也称文件句柄）。打开的文件应当及时关闭，否则过多的文件对象容易造成内存占用，导致程序运行内存不足。按照是否调用文件对象的 `close` 方法，有两种打开和关闭文件的代码书写方式：

1. 显式 close
2. 隐式 close

**显式 close**

```python
def open1():
    f = open('./students.dat')
    try:
        lines = f.readlines()
        print(lines)
    finally:
        f.close()
```

**隐式 close**

```python
def open2():
    with open('./students.dat') as f:
        lines = f.readlines()
        print(lines)
```

> 支持 `with` 语句的对象需要实现 `__enter__` 和 `__exit__` 两个方法，其中 `TextIOWrapper` 类实现了 `__exit__` 方法，`IOBase` 类实现了 `__enter__` 方法。

上述两个函数的作用相同，都是用于打开 `students.dat` 文件并打印所有的行。

### open 函数

内建函数 `open` 的签名如下：

```python
def open(
    file: _OpenFile,
    mode: OpenTextMode = ...,
    buffering: int = ...,
    encoding: str | None = ...,
    errors: str | None = ...,
    newline: str | None = ...,
    closefd: bool = ...,
    opener: _Opener | None = ...,
) -> TextIOWrapper: ...
```

1. file：字符串文件路径或实现了 `os.PathLike` 抽象类的实例；
2. mode：打开模式，默认 `r`，可选；
3. buffering：设置缓冲策略，可选；
4. encoding：编码格式，可选；
5. errors：编码或解码发生错误时的错误信息，可选；
6. newline：断行的方式，可用的参数值有 `None`、`' '`、`'\n'`、`'\r'` 和 `'\r\n'`，可选；
7. closefd：必须为 `True`，否则报错，可选；
8. opener：自定义的打开器，调用的函数，返回一个文件描述符，可选；

#### os.PathLike 抽象类

`os.PathLike` 是一个抽象类，定义了 `__fspath__` 方法，任何实现了 `__fspath__` 方法的类的实例都可以作为 `open` 函数的 file 参数值。

```python
import os

class MyFile(os.PathLike):
    def __init__(self, filename) -> None:
        self.filename = filename

    def __fspath__(self):
        return self.filename
    

def open3():
    with open(MyFile('./students.dat')) as f:
        lines = f.readlines()
        print(lines)
```

#### mode 参数值

mode 参数值可为：

| Mode  |                                                              |
| :---- | :----------------------------------------------------------- |
| `'r'` | 读打开（默认）                                               |
| `'w'` | 读打开，若文件不存在则创建，若文件存在则会清空文件内容       |
| `'x'` | 以独占的方式创建文件，如果文件已存在则报错                   |
| `'a'` | 以追加的形式打开文件，文件不存在会创建，文件存在的话，不会清空文件内容 |
| `'t'` | 文本模式（默认）                                             |
| `'b'` | 二进制打开文件                                               |
| `'+'` | 以更新的方式打开文件                                         |

### f.close 方法

当文件对象调用 `close` 方法后，对象的 `closed` 属性会置为 `True`，也可以通过该属性可以检查文件对象是否关闭。

```python
def view_f_closed():
    f = open('./students.dat')
    f.close()
    print(f.closed)
```

## 读写文件

文件对象中有如下几个方法可用于读取文件内容：

| 方法                                                         | 用法                                                         |
| ------------------------------------------------------------ | ------------------------------------------------------------ |
| [.read(size=-1)](https://docs.python.org/3.7/library/io.html#io.RawIOBase.read) | 按指定的字节数读取文件内容，当 size 为 `-1` 时，表示读取全部 |
| [.readline(size=-1)](https://docs.python.org/3.7/library/io.html#io.IOBase.readline) | 按指定的字符数读取一行的内容，当 size 为 `None` 或 `-1` 时，表示读取整行内容 |
| [.readlines()](https://docs.python.org/3.7/library/io.html#io.IOBase.readlines) | 读取文件中所有的行                                           |

文件对象中有如下几个方法可用于写入文件内容：

| 方法                                                         | 用法                                   |
| ------------------------------------------------------------ | -------------------------------------- |
| [.write(string)](https://docs.python.org/3/library/io.html?highlight=textiowrapper#io.TextIOBase.write) | 向文件写入字符串                       |
| [.writelines(seq)](https://docs.python.org/3/library/io.html?highlight=textiowrapper#io.IOBase.writelines) | 向文件中写入多行，换行符需要开发者指定 |

### 编码问题

若文件中存在中文，需要指定 encoding 参数的值，如：

```python
def read_chinese():
    with open('./students.zh-cn.dat', encoding='utf-8') as f:
        lines = f.readlines()
        print(lines)
```

### 遍历文件所有的行

下面列表读取文件所有的行的几种方式。

**方式 1**

```python
def iterate_lines1():
    with open('./students.dat', mode='r') as f:
        line = f.readline()
        while line != '':
            print(line, end='')
            line = f.readline()
```

**方式 2**

```python
def iterate_lines2():
    with open('./students.dat', mode='r') as f:
        lines = f.readlines()
        for line in lines:
            print(line, end='')
```

**方式 3**

```python
def iterate_lines3():
    with open('./students.dat', mode='r') as f:
        for line in f:
            print(line, end='')
```

## 总结

本文的完整代码如下：

```python
import os

def open1():
    f = open('./students.dat')
    try:
        lines = f.readlines()
        print(lines)
    finally:
        f.close()


def open2():
    with open('./students.dat') as f:
        lines = f.readlines()
        print(lines)


class MyFile(os.PathLike):
    def __init__(self, filename) -> None:
        self.filename = filename

    def __fspath__(self):
        return self.filename


def open3():
    with open(MyFile('./students.dat')) as f:
        lines = f.readlines()
        print(lines)


def view_f_closed():
    f = open('./students.dat')
    f.close()
    print(f.closed)


def read_chinese():
    with open('./students.zh-cn.dat', encoding='utf-8') as f:
        lines = f.readlines()
        print(lines)


def iterate_lines1():
    with open('./students.dat', mode='r') as f:
        line = f.readline()
        while line != '':
            print(line, end='')
            line = f.readline()


def iterate_lines2():
    with open('./students.dat', mode='r') as f:
        lines = f.readlines()
        for line in lines:
            print(line, end='')


def iterate_lines3():
    with open('./students.dat', mode='r') as f:
        for line in f:
            print(line, end='')

if __name__ == '__main__':
    print('书写方式 1:')
    open1()
    print('书写方式 2:')
    open2()
    print('实现了 os.PathLike 抽象类')
    open3()
    print('调用 f.close 后:')
    view_f_closed()
    print('遍历所有的行 1:')
    iterate_lines1()
    print('遍历所有的行 2:')
    iterate_lines2()
    print('遍历所有的行 3:')
    iterate_lines3()
    print("读取中文:")
    read_chinese()
```

运行输出为：

```bash
书写方式 1:
['xiaoming\n', 'xiaohong\n', 'xiaolei\n', 'xiaopang\n']
书写方式 2:
['xiaoming\n', 'xiaohong\n', 'xiaolei\n', 'xiaopang\n']
实现了 os.PathLike 抽象类
['xiaoming\n', 'xiaohong\n', 'xiaolei\n', 'xiaopang\n']
调用 f.close 后:
True
遍历所有的行 1:
xiaoming
xiaohong
xiaolei
xiaopang
遍历所有的行 2:
xiaoming
xiaohong
xiaolei
xiaopang
遍历所有的行 3:
xiaoming
xiaohong
xiaolei
xiaopang
读取中文:
['小明\n', '小红\n', '小磊\n', '小胖\n']
```

本文可总结为以下几点：

1. 打开文件后，应该及时关闭。如果不想显式地调用 `close` 方法，推荐使用 `with` 语句；
2. `open` 函数的 `mode` 参数指定了文件打开的模式，`mode` 参数的几个可选值非常重要；
3. 列举了几个读写文件的常用方法，如读的 `read`、`readline` 和 `readlines`；写的 `write` 和 `writelines`；
4. 3 种遍历文件所有行的方式；

