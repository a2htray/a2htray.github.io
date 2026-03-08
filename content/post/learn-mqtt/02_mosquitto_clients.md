+++
date = '2026-03-04T17:15:48+08:00'
draft = false
title = 'MQTT：Mosquitto 开源客户端 terdia/mqttui（三）'
categories = ['后端技术', 'MQTT']
tags = ['MQTT', 'Docker', 'mosquitto', 'mqttui']
toc = true
+++

项目地址：[https://github.com/terdia/mqttui](https://github.com/terdia/mqttui)

`terdia/mqttui` 是基于 Flask 的一款开源 Web 应用，用于实时可视化展示 MQTT（消息队列遥测传输）的消息流。它支持用户通过简洁直观的 Web 界面监控 MQTT 主题、发布消息，并查看消息统计信息。

> 提供 Docker 部署化方式。

使用 docker-compose 的方式部署，`docker-compose.yml` 内容如下：

```yaml
services:
  mosquitto:
    image: eclipse-mosquitto:2
    container_name: mosquitto
    ports:
      - 1883:1883
      - 9001:9001
    volumes:
      - ./.mosquitto/data:/mosquitto/data
      - ./.mosquitto/log:/mosquitto/log
      - ./.mosquitto/config:/mosquitto/config
  mqtt_web:
    image: terdia07/mqttui:v1.3.2
    container_name: mqtt_web
    ports:
      - 5000:5000
    environment:
      - MQTT_BROKER=mosquitto
      - MQTT_PORT=1883
    depends_on:
      - mosquitto
```

镜像 `terdia07/mqttui:v1.3.2` 支持的[环境变量](https://github.com/terdia/mqttui/blob/main/docker-compose.yml)，有：

```yaml
- DEBUG=False
- HOST=0.0.0.0
- PORT=5000
- MQTT_BROKER=mosquitto
- MQTT_PORT=1883
- MQTT_USERNAME=
- MQTT_PASSWORD=
- MQTT_KEEPALIVE=60
- MQTT_VERSION=3.1.1
- SECRET_KEY=your-secret-key
- LOG_LEVEL=INFO
- MQTT_TOPICS=#
- DB_ENABLED=True
- DB_PATH=/app/data/mqtt_messages.db
- DB_MAX_MESSAGES=10000
- DB_CLEANUP_DAYS=30
```

```bash
docker-compose -f docker-compose.yml up -d
```

```bash
docker-compose -f docker-compose.yml ps -a
```

```bash
NAME        IMAGE                    COMMAND                  SERVICE     CREATED          STATUS          PORTS
mosquitto   eclipse-mosquitto:2      "/docker-entrypoint.…"   mosquitto   27 minutes ago   Up 27 minutes   0.0.0.0:1883->1883/tcp, 0.0.0.0:9001->9001/tcp
mqtt_web    terdia07/mqttui:v1.3.2   "/entrypoint.sh"         mqtt_web    27 minutes ago   Up 27 minutes   0.0.0.0:5000->5000/tcp
```

访问 [http://localhost:5000/](http://localhost:5000/)，如下图：

![](/imgs/learn-mqtt/2026-03-04_172717_901.png)

## 评价

在界面上操作了一下，一般般，不是很推荐，原因如下：

1. 只有一个页面，没有其它复杂功能操作
2. Message flow 的可视化，没什么用
3. 有 BUG，如超出视角的内容无法滚动
