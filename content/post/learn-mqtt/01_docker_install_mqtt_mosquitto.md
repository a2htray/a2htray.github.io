+++
date = '2026-03-03T11:05:44+08:00'
draft = false
title = 'MQTT：Docker 安装 MQTT 服务 Mosquitto（二）'
categories = ['后端技术', 'MQTT']
tags = ['MQTT', 'Docker']
toc = true
+++

Mosquitto 是 MQTT 的开源实现之一，GitHub 地址 [https://github.com/eclipse-mosquitto/mosquitto](https://github.com/eclipse-mosquitto/mosquitto)，介绍信息如下：

> Mosquitto is an open source implementation of a server for version 5.0, 3.1.1, and 3.1 of the MQTT protocol. It also includes a C and C++ client library, the mosquitto_pub mosquitto_rr, and mosquitto_sub utilities for publishing and subscribing, and the mosquitto_ctrl, mosquitto_signal, and mosquitto_passwd applications for helping administer the broker.

* 实现了 MQTT 的 5.0、3.1.1 和 3.1 版本
* 提供 C 和 C++ 的客户端调用库
* mosquitto_pub、mosquitto_rr 用于订阅与发布
* mosquitto_ctrl、mosquitto_signal、mosquitto_passwd 用于服务器管理

其它方式的安装见其官网文档 [https://mosquitto.org/download/](https://mosquitto.org/download/)。

## 快速使用

使用默认的配置运行 mosquitto：

```bash
docker run -d \
  --name mosquitto \
  -p 1883:1883 \
  -p 9001:9001 \
  eclipse-mosquitto:2
```

* 端口 1883 处理标准 MQTT 连接
* 端口 9001 是 WebSocket 端口

在默认情况下，mosquitto 不接受外部 IP 的连接请求，所以需要对其进行配置。

## 数据挂载与配置

创建 3 个目录，分别用于存储数据（data）、日志（log）和配置（config）。

```bash
mkdir -p .mosquitto/{data,log,config}
```

在 `.mosquitto/config` 目录下创建 `mosquitto.conf` 文件，内容如下：

```yaml
# 监听的端口
listener 1883
protocol mqtt

# WebSocket 监听的端口
listener 9001
protocol websockets

# 允许匿名连接
allow_anonymous true

# 数据持久化处理
persistence true
# 数据持久化目录
persistence_location /mosquitto/data/

# 日志配置
log_dest file /mosquitto/log/mosquitto.log
log_dest stdout
log_type all
```

关闭已运行的容器并重新执行：

```bash
docker stop mosquitto && docker rm mosquitto
docker run -d \
  --name mosquitto \
  -p 1883:1883 \
  -p 9001:9001 \
  -v `pwd`/.mosquitto/data:/mosquitto/data \
  -v `pwd`/.mosquitto/log:/mosquitto/log \
  -v `pwd`/.mosquitto/config:/mosquitto/config \
  eclipse-mosquitto:2
```

查看运行状态：

```bash
docker ps -a
CONTAINER ID   IMAGE                 COMMAND                  CREATED         STATUS         PORTS                                            NAMES
7262f3e3d40b   eclipse-mosquitto:2   "/docker-entrypoint.…"   7 seconds ago   Up 7 seconds   0.0.0.0:1883->1883/tcp, 0.0.0.0:9001->9001/tcp   mosquitto
```

## 验证

使用 mosquitto 自带的 mosquitto_pub 命令行工具进行验证。

```bash
# 终端一
docker exec mosquitto mosquitto_sub -t "temperature" -v

# 终端二
docker exec mosquitto mosquitto_pub -t "temperature" -m '{"value": 22.5}'
```

则终端一中显示：

```bash
temperature {"value": 22.5}
```

## 密码验证

mosquitto 支持用户密码验证配置，步骤如下：

步骤一：使用 mosquitto_passwd 创建 password 文件

```bash
docker exec mosquitto mosquitto_passwd -c -b /mosquitto/config/passwd appuser appuser123
```

同时创建用户 appuser，密码：appuser123。

步骤二：修改 mosquitto.conf 配置文件部分内容

```yaml
# 关闭匿名连接
allow_anonymous false
# 指定 password 文件
password_file /mosquitto/config/passwd
```

步骤三：重启服务

```bash
docker restart mosquitto
```

步骤四：带密码验证订阅发布

```bash
# 终端一
docker exec mosquitto mosquitto_sub -t "temperature" -u appuser -P "appuser123"

# 终端二
# 正常发送
docker exec mosquitto mosquitto_pub -t "temperature" -u appuser -P "appuser123" -m '{"value": 24.5}'
# 无法发送
docker exec mosquitto mosquitto_pub -t "temperature" -m '{"value": 24.5}'
```

## 访问控制列表

密码验证后，可通过 `acl.conf` 文件控制用户对主题的访问，访问控制列表文件的示例如下：

```yaml
# device1 can publish to its own sensor topics and subscribe to commands
user device1
topic readwrite home/device1/#
topic read commands/device1/#

# device2 has access to its own topics only
user device2
topic readwrite home/device2/#
topic read commands/device2/#

# webapp can read all home topics and write commands
user webapp
topic read home/#
topic write commands/#
```

* `user` 关键字指明用户
* `topic` 关键字声明对主题的权限
* 权限可以是可读可写（readwrite）、只读（read）、只写（write）

访问控制列表文件命名为 `acl.conf`，放置在目录 `/mosquitto/config` 下，配置文件 `mosquitto.conf` 中新增内容：

```yaml
acl_file /mosquitto/config/acl.conf
```

## 其它

### 保留消息 retained message

保留消息：客户端一上线，就能拿到的最新的一条消息，核心特点：

1. 一个主题，最多只有一条消息
2. 新的保留消息会覆盖旧的
3. 客户端刚订阅就能收到，不用等下一次发布
4. 用于展示当前状态，不用于事件流

## 资源

* [How to Run Mosquitto MQTT Broker in Docker](https://oneuptime.com/blog/post/2026-02-08-how-to-run-mosquitto-mqtt-broker-in-docker/view)