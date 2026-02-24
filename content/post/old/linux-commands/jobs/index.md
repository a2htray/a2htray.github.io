---
title: "Linux jobs 命令"
date: 2023-03-12T21:38:18+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Shell
tags:
 - Linux Shell
 - jobs
 - fg
 - bg
---

jobs 命令主要用于显示系统中的任务列表及运行状态。在 Linux 中，每一个 job 都有一个唯一 ID，系统管理员通过任务 ID 对任务进行管理，可使其在前后或后台运行。通常任务和进程是等价的，只在于说侧重不同。即任务之于用户，相应地，进程之于系统。

<!--more-->

## 语法及选项

jobs 命令的语法如下：

```bash
jobs [-lnprs] 
```

* -l：在标准信息显示的基础上添加任务 ID 信息
* -n：只显示与上一次显示状态不同的任务
* -p：只显示任务 ID
* -r：显示运行中的任务
* -s：显示暂停中的任务

以 sleep 命令为例：

```bash
$ sleep 100 &
[1] 803
$ jobs
[1]+  Running                 sleep 100 &
$ jobs -l
[1]+   803 Running                 sleep 100 &
$ jobs -lr
[1]+   803 Running                 sleep 100 &
```



## fg 和 bg

fg（foreground）可以将一个任务转置前台运行，相反地，bg（background）则可以将一个任务转置后台运行。两个命令的语法分别如下：

```bash
fg [JOB_SPEC]
bg [JOB_SPEC]
```

JOB_SPEC 可以是以下形式：

* %n：n 是任务 ID
* %abc：以 abc 开头命令启动的任务
* %?abc：以包含 abc 命令启动的任务
* %：特指上一个任务

首先，编写一个 shell 脚本，实现的是每隔离 1s 打印一次 hello world。然后，后台执行该脚本并使用 fg 和 bg 进行管理。最后，使用 ctrl+c 结束任务的执行。

```bash
# file: abc.sh
#!/bin/bash

while :
do
	sleep 1
	echo "hello world"
done
```

步骤一：后台运行 abc.sh

```bash
$ ./abc.sh &
[1] 508
$ jobs
[1]+  Running                 ./abc.sh > abc.log &
```

步骤二：fg 命令使其在前台运行

```bash
$ fg %1
./abc.sh > abc.log
```

步骤三：Ctrl+Z 暂停运行

```bash
^Z
```

步骤四：bg 命令使其在后台运行

```bash
$ bg %1
[1]+ ./abc.sh > abc.log &
```

步骤五：Ctrl+C 结束运行

```bash
$ fg %1
./abc.sh > abc.log
^C
$ jobs
```

