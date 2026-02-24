+++
date = '2025-12-31T15:07:14+08:00'
draft = false
title = 'Git Commit 修改作者信息'
categories = ['生产力工具', 'Git']
tags = ['Git', '解决问题']
toc = true
+++

## 起因

现在公司上班，用的是我个人电脑，所以 git 信息是我个人的。

新建项目，推送分支持时，导致线上仓库的 commit 的作者信息是我个人的。

原因就是，项目级别的作者信息 `user.name` 和 `user.email` 没设置，git 自动使用了全局的作者信息。

## 改

毕竟是公司的项目，怎么可以用到个人的信息，改！👊🐼

**第一步：设置项目级的作者信息**

```bash
$ git config user.name "公司用户名"
$ git config user.email "公司邮箱"
```

**第二步：修改 commit 作者信息（但不改 commit 消息）**

```bash
$ git commit --amend --reset-author --no-edit
```

> ❗❗❗上述命令只能修改最近一次 commit 的作者信息。

**第三步：推送至远程仓库**

```bash
$ git push --force-with-lease origin <分支名>
```

## 延伸

--force-with-lease 选项

先检查远程分支内容，若存在本地未有分支，则拒绝推送。相对 `-f`（强制推送），更安全。




