---
title: "Linux history 命令"
date: 2023-03-25T20:26:08+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Shell
tags:
 - Linux Shell
 - jobs
---

`history` 命令可用于浏览历史执行过的命令，该命令在 Bourne Shell 中不可用，但 Bash 和 Korn 则支持该特性。在 Bash 和 Korn 中，每一次命令的执行都被视为一次事件，并且配有相应的事件号（Event Number）。在需要的情况下，可以通过事件号将执行过的命令再次执行。

不带选项的 `history` 会在终端输出所有执行过的历史命令。

<!--more-->

## 选项

### -d

如果要删除一个历史执行命令，可在`-d` 选项后面可以接一个命令关联的事件号，如下：

```bash
$ history
  518  echo 222
  519  echo 333
  520  echo 444
  521  echo 555
  522  history 5
$ history -d 520
$ history 5
  519  echo 333
  520  echo 555
  521  history 5
  522  history -d 520
  523  history 5
```

### -c

移除所有历史命令。

```bash
$ history -c
```

### -w

将历史执行命令记录写入到 `.bash_history` 文件。

```bash
$ history -w
```



## 环境变量

与 `history` 相关的环境变量包括：

* `HISTSIZE`：历史执行命令列表的数量
* `HISTFILESIZE`：保存到 `.bash_history` 的历史执行命令数
* `HISTTIMEFORMAT`：显示命令时的时间格式字符串
* `HISTCONTROL`：记录控制变量

```bash
$ echo $HISTSIZE
500
$ echo $HISTFILESIZE
500
```

## 示例

**限制显示的命令数**

```bash
$ history 5
  505  cd -
  506  ls
  507  expr 1 + 2
  508  echo 'hello world'
  509  history 5
```

**通过事件号执行命令**

```bash
$ !508
echo 'hello world'
hello world
```

**通过事件号预览命令 :p**

```bash
$ !507:p
expr 1 + 2
```

**查看并执行最近一条执行过的命令**

```bash
$ !!
expr 1 + 2
3
```

**查看并执行最近第 n 条执行过的命令**

```bash
$ history 5
  523  history 5
  524  echo HISTSIZE
  525  echo $HISTSIZE
  526  echo $HISTFILESIZE
  527  history 5
$ !-3
echo $HISTSIZE
500 
```

**执行指定开头的历史命令**

```bash
$ !exp
expr 1 + 2
3
```

**运行一条命令，但不对其进行记录**

```bash
$ echo "hello world"; history -d $(history 1)
```

**在命令前输出执行的时间**

```bash
$ HISTTIMEFORMAT="%c "
$ history 1
537  Sat Mar 25 22:12:14 2023 history 1
```

**对于重复连续执行的命令只记录一次**

```bash
$ HISTCONTROL="ignoredups"
$ pwd
/Users/a2htray/workspace/a2htray.github.io
$ pwd
/Users/a2htray/workspace/a2htray.github.io
$ pwd
/Users/a2htray/workspace/a2htray.github.io
$ history 5
  537  Sat Mar 25 22:12:14 2023 history 1
  538  Sat Mar 25 22:14:16 2023 echo $HISTCONTROL
  539  Sat Mar 25 22:16:42 2023 HISTCONTROL="ignoredups"
  540  Sat Mar 25 22:16:49 2023 pwd
  541  Sat Mar 25 22:16:54 2023 history 5
```



## 参考

1. [history command in Linux with Examples](https://www.geeksforgeeks.org/history-command-in-linux-with-examples/)
2. [How to Use the Linux history Command](https://phoenixnap.com/kb/linux-history-command)
3. [Bash History Command Examples In Linux](https://www.rootusers.com/17-bash-history-command-examples-in-linux/)

