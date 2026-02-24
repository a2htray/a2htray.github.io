+++
date = '2025-12-16T13:55:38+08:00'
draft = false
title = '在 Python 镜像的基础上换 apt 源（阿里云）'
categories = ['后端技术', 'Docker']
tags = ['工作', '解决问题', 'Docker', 'Dockerfile', 'Python', 'apt 源']
+++

介绍如何在 Dockerfile 中对基础镜像进行换源，加快软件的安装。

## Python 的基础系统

官方 Python 的基础系统分成两类：

1. Debian：标准标签，如 3.12、3.12-slim
2. Alpine Linux：Alpine 标签，如 3.12-alpine

本文方法适用于第 1 类。

## Docker 指令

将以下指令信息复制到 Dockerfile 中，放置到 `apt-get update` 之前：

```bash
RUN rm -rf /etc/apt/sources.list* \
    && ls -l /etc/apt/ | grep -v sources || true

RUN echo "deb http://mirrors.aliyun.com/debian/ bookworm main contrib non-free non-free-firmware" > /etc/apt/sources.list \
    && echo "deb http://mirrors.aliyun.com/debian/ bookworm-updates main contrib non-free non-free-firmware" >> /etc/apt/sources.list \
    && echo "deb http://mirrors.aliyun.com/debian-security/ bookworm-security main contrib non-free non-free-firmware" >> /etc/apt/sources.list \
    && echo "deb http://mirrors.aliyun.com/debian/ bookworm-backports main contrib non-free non-free-firmware" >> /etc/apt/sources.list
```

* 第 1 个 RUN 指令：删除默认的源信息
* 第 2 个 RUN 指令：将阿里云的源写入到 `sources.list` 文件