+++
date = '2026-10-16T18:16:05+08:00'
draft = false
title = 'RocketMQ：使用 5.x 客户端及代码示例'
categories = ['后端技术', 'RocketMQ']
tags = ['RocketMQ']
toc = true
+++

![](/imgs/learn-rocketmq/20261019-2.png)

## vs. 4.x 客户端

Python 4.x SDK 是 Remoting TCP 协议的老富客户端；Python 5.x（`rocketmq-python-client`）是全新 gRPC 协议 SDK，两套代码完全重构，API 不兼容，不是简单小版本升级

**服务端兼容性**

- **Python 4.x SDK**：可连接 RocketMQ 4.x/ 5.x 服务端
- **Python 5.x gRPC SDK**：**只能连接 RocketMQ 5.x 服务端，无法连 4.x broker**

## Docker 部署

编写 docker-compose.yml 文件，内容如下：

```yml
services:
  namesrv:
    image: apache/rocketmq:5.5.0
    platform: linux/amd64
    container_name: rmqnamesrv
    ports:
      - "9876:9876"
    environment:
      - JAVA_OPT_EXT=-Xms256m -Xmx256m -Xmn128m
    command: sh mqnamesrv
    restart: unless-stopped

  broker:
    image: apache/rocketmq:5.5.0
    platform: linux/amd64
    container_name: rmqbroker
    ports:
      - "10909:10909"
      - "10911:10911"
      - "10912:10912"
    environment:
      - NAMESRV_ADDR=namesrv:9876
      - AUTO_CREATE_TOPIC_ENABLE=true
      - JAVA_OPT_EXT=-Xms512m -Xmx512m -Xmn256m
    command: sh mqbroker
    depends_on:
      - namesrv
    restart: unless-stopped

  proxy:
    image: apache/rocketmq:5.5.0
    platform: linux/amd64
    container_name: rmqproxy
    ports:
      - "8081:8081"
    environment:
      - NAMESRV_ADDR=namesrv:9876
      - JAVA_OPT_EXT=-Xms256m -Xmx256m -Xmn128m
    command: sh mqproxy
    depends_on:
      - namesrv
      - broker
    restart: on-failure
```

三个组件之间的关系及请求流程如下图：

![](/imgs/learn-rocketmq/20261019.png)

**安装依赖包**

```bash
pip install rocketmq-python-client
```

## 代码

**连接 Proxy，启动 Producer 实例**

```python
# main.py
from rocketmq import ClientConfiguration, Credentials, Producer

# 配置 Proxy 地址
# 分号分隔多个 Proxy 地址
endpoints = "localhost:8081"

client_config = ClientConfiguration(endpoints, Credentials())

# 创建 Producer 实例
producer = Producer(client_config)

# 启动 Producer
try:
    producer.startup()
    print("producer startup")
except Exception as e:
    print(f"producer startup exception: {e}")
finally:
    producer.shutdown()
    print("producer shutdown")
```

执行：

```bash
$ python main.py
producer startup
producer shutdown
```

**生产者、消费者独立脚本**

```python
# consumer.py
import argparse
import time

from rocketmq import (
    ClientConfiguration, Credentials, PushConsumer,
    MessageListener, ConsumeResult, FilterExpression, Message,
)


class Listener(MessageListener):
    def consume(self, msg: Message) -> ConsumeResult:
        print(
            f"receive | topic={msg.topic} tag={msg.tag} "
            f"body={msg.body.decode('utf-8')} keys={msg.keys} "
            f"message_id={msg.message_id}"
        )
        return ConsumeResult.SUCCESS


def parse_args():
    p = argparse.ArgumentParser(description="RocketMQ 5.x 消费者")
    p.add_argument("--endpoints", required=True,
                   help="gRPC 接入点, 例如 localhost:8081")
    p.add_argument("--topic", required=True,
                   help="订阅的 topic")
    p.add_argument("--group", required=True,
                   help="消费者组 (必填)")
    p.add_argument("--tag", default="",
                   help="tag 过滤表达式, 为空表示匹配该 topic 下全部消息")
    return p.parse_args()


def main():
    args = parse_args()

    client_config = ClientConfiguration(args.endpoints, Credentials())
    subscription = {args.topic: FilterExpression(args.tag)}
    consumer = PushConsumer(
        client_config, args.group, Listener(), subscription
    )
    consumer.startup()
    print(
        f"consumer started | group={args.group} topic={args.topic} "
        f"tag='{args.tag}', Ctrl+C 退出"
    )

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        pass
    finally:
        consumer.shutdown()
        print("consumer shutdown")


if __name__ == "__main__":
    main()
```

```python
# producer.py
import argparse

from rocketmq import (
    ClientConfiguration, Credentials, Producer, Message,
)


def parse_args():
    p = argparse.ArgumentParser(description="RocketMQ 5.x 生产者")
    p.add_argument("--endpoints", required=True,
                   help="gRPC 接入点, 例如 localhost:8081")
    p.add_argument("--topic", required=True,
                   help="目标 topic")
    p.add_argument("--group", default="",
                   help="生产者组 (5.x gRPC 客户端通常不需要, 这里仅作预留参数)")
    p.add_argument("--tag", default="",
                   help="消息 tag, 为空表示不带 tag")
    p.add_argument("--count", type=int, default=10,
                   help="发送消息条数, 默认 10")
    p.add_argument("--body", default="hello rocketmq",
                   help="消息内容模板, 实际发送会在末尾追加序号")
    return p.parse_args()


def main():
    args = parse_args()

    client_config = ClientConfiguration(args.endpoints, Credentials())
    producer = Producer(client_config)
    producer.startup()
    print(f"producer startup, group={args.group or '(none)'}")

    try:
        for i in range(args.count):
            msg = Message()
            msg.topic = args.topic
            if args.tag:
                msg.tag = args.tag
            msg.body = f"{args.body}-{i}".encode("utf-8")
            msg.keys = f"key-{i}"
            receipt = producer.send(msg)
            print(
                f"[{i}] send ok | topic={args.topic} tag={args.tag} "
                f"keys={msg.keys} message_id={receipt.message_id}"
            )
    finally:
        producer.shutdown()
        print("producer shutdown")


if __name__ == "__main__":
    main()
```

开启两个终端，分别执行：

```bash
$ python consumer.py --endpoints localhost:8081 --topic testTopic \
        --group consumerGroup --tag testTag
consumer started | group=consumerGroup topic=testTopic tag='testTag', Ctrl+C 退出
receive | topic=testTopic tag=testTag body=hello-0 keys={'key-0'} message_id=010EC4EE2179D31F490AE76F7D00000000
receive | topic=testTopic tag=testTag body=hello-1 keys={'key-1'} message_id=010EC4EE2179D31F490AE76F7E00000001
receive | topic=testTopic tag=testTag body=hello-2 keys={'key-2'} message_id=010EC4EE2179D31F490AE76F7E00000002
receive | topic=testTopic tag=testTag body=hello-3 keys={'key-3'} message_id=010EC4EE2179D31F490AE76F7E00000003
receive | topic=testTopic tag=testTag body=hello-5 keys={'key-5'} message_id=010EC4EE2179D31F490AE76F7E00000005
receive | topic=testTopic tag=testTag body=hello-4 keys={'key-4'} message_id=010EC4EE2179D31F490AE76F7E00000004
receive | topic=testTopic tag=testTag body=hello-6 keys={'key-6'} message_id=010EC4EE2179D31F490AE76F7E00000006
receive | topic=testTopic tag=testTag body=hello-7 keys={'key-7'} message_id=010EC4EE2179D31F490AE76F7E00000007
receive | topic=testTopic tag=testTag body=hello-8 keys={'key-8'} message_id=010EC4EE2179D31F490AE76F7E00000008
receive | topic=testTopic tag=testTag body=hello-9 keys={'key-9'} message_id=010EC4EE2179D31F490AE76F7E00000009
```

```bash
$ python producer.py --endpoints localhost:8081 --topic testTopic \
        --group producerGroup --tag testTag --count 10 --body "hello"
producer startup, group=producerGroup
[0] send ok | topic=testTopic tag=testTag keys={'key-0'} message_id=010EC4EE2179D31F490AE76F7D00000000
[1] send ok | topic=testTopic tag=testTag keys={'key-1'} message_id=010EC4EE2179D31F490AE76F7E00000001
[2] send ok | topic=testTopic tag=testTag keys={'key-2'} message_id=010EC4EE2179D31F490AE76F7E00000002
[3] send ok | topic=testTopic tag=testTag keys={'key-3'} message_id=010EC4EE2179D31F490AE76F7E00000003
[4] send ok | topic=testTopic tag=testTag keys={'key-4'} message_id=010EC4EE2179D31F490AE76F7E00000004
[5] send ok | topic=testTopic tag=testTag keys={'key-5'} message_id=010EC4EE2179D31F490AE76F7E00000005
[6] send ok | topic=testTopic tag=testTag keys={'key-6'} message_id=010EC4EE2179D31F490AE76F7E00000006
[7] send ok | topic=testTopic tag=testTag keys={'key-7'} message_id=010EC4EE2179D31F490AE76F7E00000007
[8] send ok | topic=testTopic tag=testTag keys={'key-8'} message_id=010EC4EE2179D31F490AE76F7E00000008
[9] send ok | topic=testTopic tag=testTag keys={'key-9'} message_id=010EC4EE2179D31F490AE76F7E00000009
producer shutdown
```

## 与 v4 客户端对比

| 维度 | v4（`rocketmq-client-python`） | v5（`rocketmq-python-client`） |
|------|-------------------------------|-------------------------------|
| 底层 | 封装 C++ 客户端 | gRPC |
| 连接目标 | 直连 NameServer **9876** | 连 Proxy **gRPC 8081** |
| 地址 API | `producer.set_name_server_address("ip:9876")` | `ClientConfiguration("ip:8081", Credentials())` |
| 启动 | `producer.start()` | `producer.startup()` |
| 消费 | `PushConsumer` + `subscribe(topic, callback)` | `PushConsumer(config, group, Listener, {topic: FilterExpression})` |