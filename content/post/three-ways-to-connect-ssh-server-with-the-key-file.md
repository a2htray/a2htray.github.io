+++
date = '2025-12-11T15:06:39+08:00'
draft = false
title = '3 种使用密钥文件连接 SSH 服务的场景'
categories = ['生产力工具', 'SSH']
tags = ['工作', '解决问题']
+++

SSH 服务除了可以通过用户名和密码登录之外，还可以使用密钥文件进行登录。

本文列举 3 种我使用密钥文件的场景。假设我的密钥的放在：

```bash
./.ssh/xxx.pem
```

## 终端命令 ssh

在终端的场景下，直接使用 `-i` 参数，指定密钥文件即可。

```bash
$ ssh -i ~/.ssh/xxx.pem username@host
```

## VS Code 远程连接

VS Code 的 SSH 连接要使用到 `Remote-SSH` 插件，如果要使用密钥文件，有两个入口可以进行配置。

第一种：通过 VS Code。

![](/imgs/three-ways-to-connect-ssh-server-with-the-key-file/vscode-step-1.png)

![](/imgs/three-ways-to-connect-ssh-server-with-the-key-file/vscode-step-2.png)

![](/imgs/three-ways-to-connect-ssh-server-with-the-key-file/vscode-step-3.png)

第二种：直接编辑 `~/.ssh/config` 文件。

在 `~/.ssh/config` 文件中加入以下内容：

```bash
Host <connection name>          # 自己定义一个连接名称
    HostName <ip or domain>     # 用 IP 或者域名
    User <username>             # 用户名
    Port <port>                 # 端口，默认 22
    IdentityFile ~/.ssh/xxx.pem # 密钥文件地址 
```

修改完后，重启下 VS Code，即可以选择新增的连接。

## FileZilla 连接

一张图解释。

![](/imgs/three-ways-to-connect-ssh-server-with-the-key-file/filezilla-steps.png)

## 小结

列举了 3 种使用 SSH 密钥的场景：

1. 命令行
2. VS Code
3. FileZilla

留个痕迹。
