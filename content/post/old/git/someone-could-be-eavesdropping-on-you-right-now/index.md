---
title: "Someone could be eavesdropping on you right now (man-in-the-middle attack)!"
date: 2023-03-25T14:07:10+08:00
draft: false
reward: false
categories:
 - 生产力工具
 - Git
tags:
 - Git
---

今天在 `git pull` 时报了如下的错：

```bash
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
@    WARNING: REMOTE HOST IDENTIFICATION HAS CHANGED!     @
@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@
IT IS POSSIBLE THAT SOMEONE IS DOING SOMETHING NASTY!
Someone could be eavesdropping on you right now (man-in-the-middle attack)!
It is also possible that a host key has just been changed.
The fingerprint for the RSA key sent by the remote host is
SHA256:xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx.
Please contact your system administrator.
Add correct host key in /Users/a2htray/.ssh/known_hosts to get rid of this message.
Offending RSA key in /Users/a2htray/.ssh/known_hosts:2
Host key for github.com has changed and you have requested strict checking.
Host key verification failed.
fatal: Could not read from remote repository.

Please make sure you have the correct access rights
and the repository exists.
```

导致无法拉取远程仓库的代码。

<!--more-->

通过网上的教程，只要删除 `/Users/a2htray/.ssh/known_hosts` 即可。

```bash
$ rm -f /Users/a2htray/.ssh/known_hosts
```

