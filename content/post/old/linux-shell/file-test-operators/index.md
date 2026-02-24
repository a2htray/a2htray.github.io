---
title: "Linux Shell：文件测试比较操作符"
date: 2023-03-25T20:13:25+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Shell
tags:
 - Linux Shell
 - Test Operators
---

在 Shell 中，我们经常需要与文件打交道，所以了解文件测试操作符十分有必要。本文罗列了文件测试操作符并逐一给出示例，其中测试表达式的形式以 `[[]]` 给出。

<!--more-->



## 文件测试操作符

文件测试操作符如下表：

| 操作符     | 描述                   | 示例             |
| ------- | -------------------- | -------------- |
| -b file | 测试给定文件是否是块设备         | [[ -b $file ]] |
| -c file | 测试给定文件是否是字符文件        | [[ -c $file ]] |
| -d file | 测试给定文件是否是目录          | [[ -d $file ]] |
| -f file | 测试给定文件是否是文件          | [[ -f $file ]] |
| -g file | 测试给定文件是否设置了 SGID 比特位 | [[ -g $file ]] |
| -k file | 测试给定文件是否设置了黏滞位       | [[ -k $file ]] |
| -p file | 测试给定文件是否是一个命名管道      | [[ -p $file ]] |
| -u file | 测试给定文件是否设置了 SUID 比特位 | [[ -u $file ]] |
| -r file | 测试给定文件是否可读           | [[ -r $file ]] |
| -w file | 测试给定文件是否可写           | [[ -w $file ]] |
| -x file | 测试给定文件是否可执行          | [[ -x $file ]] |
| -s file | 测试给定文件的大小是否大于 0      | [[ -s $file ]] |
| -L file | 测试给定文件是否是一个链接文件      | [[ -L $file ]] |

## 示例

### -b

`测试给定文件是否是块设备`

```bash
$ block_file=my_block_file
$ mknod $block_file b 1 2
$ [[ -b $block_file ]] && echo "$block_file is a block special file"
my_block_file is a block special file
```

### -c

`测试给定文件是否是字符文件`

```bash
$ character_file=my_character_file
$ mknod $character_file c 1 2
$ [[ -c $character_file ]] && echo "$character_file is a character special file"
my_character_file is a character special file
```

### -d

`测试给定文件是否是目录`

```bash
$ [[ -d /etc ]] && echo "/etc is a directory"
/etc is a directory
```

### -f

`测试给定文件是否是文件`

```bash
$ [[ -f /etc/passwd ]] && echo "/etc/passwd is a file"
/etc/passwd is a file
```

### -g

`测试给定文件是否设置了 SGID 比特位`

```bash
$ mkdir SGID_test
$ chmod g+s SGID_test
$ [[ -g ./SGID_test ]] && echo "./SGID_test has been set SGID bit"
./SGID_test has been set SGID bit
```

### -k

`测试给定文件是否设置了黏滞位`

```bash
$ mkdir sticky_file
$ chmod +t sticky_file
$ [[ -k ./sticky_file ]] && echo "./sticky_file has been set sticky bit"
./sticky_file has been set sticky bit
```

### -p

`测试给定文件是否是一个命名管道`

```bash
$ mknod named_pipe p
$ [[ -p ./named_pipe ]] && echo "./named_pipe is a named pipe"
./named_pipe is a named pipe
```

### -u

`测试给定文件是否设置了 SUID 比特位`

```bash
$ touch SUID_test
$ chmod u+s SUID_test
$ [[ -u ./SUID_test ]] && echo "./SUID_test has been set SUID bit"
./SUID_test has been set SUID bit
```

### -r

`测试给定文件是否可读`

```bash
$ touch readable_file
$ chmod u+r readable_file
$ [[ -r ./readable_file ]] && echo "./readable_file is a readable file"
./readable_file is a readable file
```

### -w

`测试给定文件是否可写`

```bash
$ touch writable_file
$ chmod u+w writable_file
$ [[ -w ./writable_file ]] && echo "./writable_file is a writable file"
./writable_file is a writable file
```

### -x

`测试给定文件是否可执行`

```bash
$ touch executable_file
$ chmod u+x executable_file
$ [[ -x ./executable_file ]] && echo "./executable_file ia an executable file"
./executable_file ia an executable file
```

### -s

`测试给定文件的大小是否大于 0`

```bash
$ echo "content" > test_file
$ [[ -s ./test_file ]] && echo "the size of ./test_file is larger than 0"
the size of ./test_file is larger than 0
```

### -L

`测试给定文件是否是一个链接文件`

```bash
$ touch origin_file
$ ln -s origin_file linked_file
$ [[ -L ./linked_file ]] && echo "./linked_file is a linked file"
./linked_file is a linked file
```

## 参考

1.  [Unix / Linux - Shell File Test Operators Example](https://www.tutorialspoint.com/unix/unix-file-operators.htm)
2.  [mknod](https://wangchujiang.com/linux-command/c/mknod.html)
3.  [SUID、SGID 详解](https://www.cnblogs.com/qinshizhi/p/9595375.html)
