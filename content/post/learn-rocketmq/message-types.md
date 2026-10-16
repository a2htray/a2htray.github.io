+++
date = '2026-10-16T14:42:37+08:00'
draft = false
title = 'RocketMQ 消息类型：同步、异步、顺序、延时、事务'
categories = ['后端技术', 'RocketMQ']
tags = ['RocketMQ', '消息类型']
toc = true
+++

![](/imgs/learn-rocketmq/20261016-3.png)

RocketMQ 按**发送方式、消费顺序、投递时机、事务语义**划分出最常用 5 类消息：**同步消息、异步消息、顺序消息、延时消息、事务消息**。

1. **同步消息**：发送方发消息后**阻塞等待 Broker 返回结果**，成功 / 失败立即感知，适合强可靠性、不能丢消息的业务；
2. **异步消息**：发送方发消息**不阻塞，通过回调函数接收结果**，吞吐量更高，适合高吞吐、允许异步感知发送结果场景；
3. **顺序消息**：保证**同一分组内消息按发送顺序被消费**，分全局顺序 / 分区顺序；牺牲部分并发，适合有先后依赖的业务；
4. **延时消息**：消息发送后**不会立刻投递，等待指定时长才对消费者可见**，不支持任意时间戳，只能选预设级别；
5. **事务消息**：实现**分布式事务（半消息 + 确认 / 回查机制）**，保证本地事务和消息发送原子性，要么本地事务成功 + 消息可见，要么都不生效。

## 同步消息

```python
# main4.py
from rocketmq.client import Producer, Message

namesrv_addr = "namesrv:9876"
topic = "test_topic"
tag = "test_tag"

producer = Producer("test_producer")
producer.set_name_server_address(namesrv_addr)
producer.start()
print("Producer 进程已启动，开始发送", flush=True)

message = Message(topic)
message.set_tags(tag)
message.set_body("test_body".encode("utf-8"))
send_result = producer.send_sync(message)

print(f"发送成功 msgId: {send_result.msg_id}, status: {send_result.status}, offset: {send_result.offset}", flush=True)
```

`producer.send_sync(message)` 用于发送同步消息，返回 `SendResult` 对象，包含消息 ID、状态、偏移量等信息。

执行：

```bash
$ docker run --rm --platform linux/amd64 --network rocketmq_rocketmq-net -v $(pwd):/app rocketmq:1.1 main4.py   
Producer 进程已启动，开始发送
发送成功 msgId: AC15000500012B6E51454EF778FF0000, status: 0, offset: 19
```

## 异步消息

由于新版本去除了异步消息的 API，所以异步代码需要自己实现，原理：使用线程池异步发送消息，回调函数接收结果。

```python
# main5.py
from concurrent.futures import ThreadPoolExecutor

from rocketmq.client import Producer, Message

executor = ThreadPoolExecutor(max_workers=10)

namesrv_addr = "namesrv:9876"
topic = "test_topic"
tag = "test_tag"

producer = Producer("test_producer")
producer.set_name_server_address(namesrv_addr)
producer.start()

message = Message(topic)
message.set_tags(tag)
message.set_body("test_body".encode("utf-8"))


def callback(future):
    print(f"[回调]发送成功 msgId: {future.result().msg_id}, status: {future.result().status}, offset: {future.result().offset}", flush=True)

future = executor.submit(producer.send_sync, message)
future.add_done_callback(callback)
```

使用 `ThreadPoolExecutor` 异步发送消息，回调函数接收结果，但在 RocketMQ 端仍为同步消息。

执行：

```bash
$ docker run --rm --platform linux/amd64 --network rocketmq_rocketmq-net -v $(pwd):/app rocketmq:1.1 main5.py
[回调]发送成功 msgId: AC1500070001373207584F01631F0000, status: 0, offset: 21
```

## 顺序消息

4.x 版本没有设置 message group 的 API，只能通过设置 key 来实现顺序消息。

```python
# main6.py
from rocketmq.client import Producer, Message

namesrv_addr = "namesrv:9876"
topic = "test_topic"
tag = "test_tag"

producer = Producer("test_producer")
producer.set_name_server_address(namesrv_addr)
producer.start()

for i in range(3):
    message = Message(topic)
    message.set_tags(tag)
    message.set_body(f"test_body_{i+1}".encode("utf-8"))
    message.set_keys("test_key")
    send_result = producer.send_sync(message)
    print(f"发送成功 msgId: {send_result.msg_id}, status: {send_result.status}, offset: {send_result.offset}", flush=True)
```

相同 key 的消息会路由到同一队列。

执行：

```bash
$ docker run --rm --platform linux/amd64 --network rocketmq_rocketmq-net -v $(pwd):/app rocketmq:1.1 main6.py
发送成功 msgId: AC150008000122FB192D4F0849F80000, status: 0, offset: 22
发送成功 msgId: AC150008000122FB192D4F084A2C0001, status: 0, offset: 19
发送成功 msgId: AC150008000122FB192D4F084A310002, status: 0, offset: 15
```

## 延时消息

4.x 版本**原生只支持 18 个预设延时等级，不能随便指定任意秒数**，不是 “实时自定义延时”。

```python
# main7.py
from rocketmq.client import Producer, Message

namesrv_addr = "namesrv:9876"
topic = "test_topic"
tag = "test_tag"

producer = Producer("test_producer")
producer.set_name_server_address(namesrv_addr)
producer.start()

for i in range(3):
    message = Message(topic)
    message.set_tags(tag)
    message.set_delay_time_level(2) # 1 到 18
    message.set_body(f"test_body_{i+1}".encode("utf-8"))
    
    send_result = producer.send_sync(message)
    print(f"发送成功 msgId: {send_result.msg_id}, status: {send_result.status}, offset: {send_result.offset}", flush=True)
```

执行：

```bash
$ docker run --rm --platform linux/amd64 --network rocketmq_rocketmq-net -v $(pwd):/app rocketmq:1.1 main7.py
发送成功 msgId: AC150009000118BA44444F17FA290000, status: 0, offset: 0
发送成功 msgId: AC150009000118BA44444F17FA6B0001, status: 0, offset: 1
发送成功 msgId: AC150009000118BA44444F17FA6E0002, status: 0, offset: 2
```

## 事务消息

4.x 版本不支持事务消息。
