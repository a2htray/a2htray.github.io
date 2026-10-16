+++
date = '2026-10-16T09:58:22+08:00'
draft = false
title = 'RocketMQ 核心概念：Topic、队列、Tag、Key'
categories = ['后端技术', 'RocketMQ']
tags = ['RocketMQ', 'Topic', '队列', '标签', 'Key']
toc = true
+++

![](/imgs/learn-rocketmq/20261016.png)

## Topic、队列、Tag、Key

1. **海量消息如何分类？** 订单消息、日志消息、通知消息混在一起怎么区分？

用 **Topic** 做一级分类。

2. **同一主题下如何让多个消费者并行处理、提升吞吐？**

用 **Message Queue（队列）** 做物理分片，它等于"并行度单位"。

3. **同一主题内如何按维度做轻量过滤？** 消费者只想要"已支付"的消息，不要全部拉过来。

用 **Tag** 做二级标签，消费者根据标签筛选消息。

4. **如何定位 / 追踪 / 去重某一条具体消息？** 出问题时能不能按"订单号"把这条消息捞出来？

用 **Key** 做业务主键标识。

总结来说，**Topic 管"归类"，Queue 管"并行"，Tag 管"细分过滤"，Key 管"单条追踪"**。

![](/imgs/learn-rocketmq/20261016-2.png)

## Topic（主题）

Topic 是消息的一级分类，也是消息的**逻辑**分类单位，生产者将消息发送到某个 Topic，消费者从某个 Topic 拉取消息。

**Topic 相关的命令行使用**

1. 指定集群创建 Topic

```bash
$ ./mqadmin updateTopic -n namesrv:9876 -c DefaultCluster -t OrderEvent -r 8 -w 8
create topic to broker:10911 success.
TopicConfig [topicName=OrderEvent, readQueueNums=8, writeQueueNums=8, perm=RW-, topicFilterType=SINGLE_TAG, topicSysFlag=0, order=false, attributes={}]
```

* `-c DefaultCluster`：指定集群名称为 `DefaultCluster`
* `-t OrderEvent`：指定 Topic 名为 `OrderEvent`
* `-r 8 -w 8`：指定读写队列数为 8

2. 列出集群所有 Topic

```bash
$ ./mqadmin topicList -n namesrv:9876
RMQ_SYS_TRANS_HALF_TOPIC
BenchmarkTest
OFFSET_MOVED_EVENT
TBW102
rmq_sys_REVIVE_LOG_DefaultCluster
SELF_TEST_TOPIC
DefaultCluster
SCHEDULE_TOPIC_XXXX
DefaultCluster_REPLY_TOPIC
OrderEvent
test_topic
rmq_sys_SYNC_BROKER_MEMBER_broker-a
RMQ_SYS_TRANS_OP_HALF_TOPIC
%RETRY%test_consumer
broker-a
```

3. 删除 Topic

```bash
$ ./mqadmin deleteTopic -n namesrv:9876 -c DefaultCluster -t OrderEvent
delete topic [OrderEvent] from cluster [DefaultCluster] success.
delete topic [OrderEvent] from NameServer success.
```

4. 获取 Topic 的路由信息

```bash
$ ./mqadmin topicRoute -n namesrv:9876 -t test_topic
{
        "brokerDatas":[
                {
                        "brokerAddrs":{0:"broker:10911"
                        },
                        "brokerName":"broker-a",
                        "cluster":"DefaultCluster",
                        "enableActingMaster":false
                }
        ],
        "filterServerTable":{},
        "queueDatas":[
                {
                        "brokerName":"broker-a",
                        "perm":6,
                        "readQueueNums":4,
                        "topicSysFlag":0,
                        "writeQueueNums":4
                }
        ]
}
```

5. 获取 Topic 的状态信息

```bash
$ ./mqadmin topicStatus -n namesrv:9876 -t test_topic
#Broker Name                      #QID  #Min Offset           #Max Offset             #Last Updated
broker-a                          0     0                     16                      2026-10-15 09:25:43,873
broker-a                          1     0                     16                      2026-10-15 09:25:44,383
broker-a                          2     0                     13                      2026-10-15 09:25:44,889
broker-a                          3     0                     6                       2026-10-15 09:25:26,018
```

## Queue（队列）

Queue 是 Topic 的**物理分片**。一个 Topic 由若干个 Queue 组成（默认读队列数 = 写队列数，常见 4 或 8，创建 Topic 时指定 `-r` 和 `-w` 参数）。

即 Topic 与 Queue 是一对多的关系。

- 消息**实际存在 Queue 里**
- 同一**消费组内**，一个 Queue 在同一时刻**只被一个消费者实例**消费
- Queue 的 FIFO 性质，想要全局顺序则只能让 Topic 的 Queue 数为 1

## Tag（标签）

Tag 是**同一 Topic 内的二级消息分类**，用于消费者订阅时的轻量过滤。

- 消费者只拉取自己关心的 Tag，减少无效网络传输与消费开销
- 生产者给消息打 Tag；消费者订阅时按 Tag 过滤消息

## Key（消息键）

Key 是挂在单条消息上的**业务维度标识**，可与 Tag 共存。

- **按 Key 查询**：RocketMQ 会为 Key 建立 `IndexFile` 索引，运维控制台 / `mqadmin` 可按 key 检索到这条消息。
- **消费幂等 / 去重**：用 Key（如订单号）做唯一性判定，防止消息重复投递导致业务被重复处理。
- **问题排查 / 消息轨迹**：用同一个 Key 串联一条消息从生产到消费的完整链路，方便定位问题。

## 发一条带 Topic / Tag / Key 的消息，并按 Tag 消费

```python
# main3.py
import multiprocessing
import time
import sys

from rocketmq.client import Producer, PushConsumer, Message, ConsumeStatus

namesrv_addr = "namesrv:9876"
topic = "test_topic"


def producer_process(msg_count: int = 10):
    tags = ["test_tag_1", "test_tag_2"]
    producer = Producer("test_producer")
    producer.set_name_server_address(namesrv_addr)
    producer.start()
    print("Producer 进程已启动，开始发送", flush=True)

    for i in range(msg_count):
        msg = Message(topic)
        # 设置消息标签，交替发送
        msg.set_tags(tags[i % len(tags)])
        msg.set_keys(f"test_key_{i}".encode("utf-8"))
        msg.set_body(f"test_body_{i}".encode("utf-8"))
        producer.send_sync(msg)
        print(f"发送成功 test_body_{i}", flush=True)
        time.sleep(0.5)
    producer.shutdown()
    print("Producer 发送完毕退出", flush=True)


def msg_listener(msg):
    body = msg.body.decode("utf-8")
    tags = msg.tags.decode("utf-8")
    keys = msg.keys.decode("utf-8")
    print(f"【Consumer 收到消息】msgId: {msg.id}, body: {body}, tag: {tags}, key: {keys}", flush=True)
    return ConsumeStatus.CONSUME_SUCCESS


def consumer_process():
    consumer = PushConsumer("test_consumer")
    consumer.set_name_server_address(namesrv_addr)
    # 只消费 test_tag_1 标签消息
    consumer.subscribe(topic, msg_listener, "test_tag_1")
    consumer.start()
    print("Consumer 进程已启动，等待消息", flush=True)
    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\nConsumer 收到中断", flush=True)
    finally:
        consumer.shutdown()
        print("Consumer shutdown", flush=True)


if __name__ == "__main__":
    multiprocessing.set_start_method("spawn", force=True)

    n = int(sys.argv[1]) if len(sys.argv) > 1 else 3

    p_consumer = multiprocessing.Process(target=consumer_process)
    p_consumer.daemon = True
    p_consumer.start()
    time.sleep(3)

    p_producer = multiprocessing.Process(target=producer_process, args=(n,))
    p_producer.start()
    p_producer.join()

    try:
        while True:
            time.sleep(1)
    except KeyboardInterrupt:
        print("\n主进程收到Ctrl+C，准备退出", flush=True)
    finally:
        print("程序退出", flush=True)
```

执行效果：

```bash
$ docker run --rm --platform linux/amd64 --network rocketmq_rocketmq-net -v $(pwd):/app rocketmq:1.1 main3.py 10
Consumer 进程已启动，等待消息
Producer 进程已启动，开始发送
发送成功 test_body_0
【Consumer 收到消息】msgId: AC15000500346F296A464E8EE8BC0000, body: test_body_0, tag: test_tag_1, key: test_key_0
发送成功 test_body_1
发送成功 test_body_2
【Consumer 收到消息】msgId: AC15000500346F296A464E8EED2F0002, body: test_body_2, tag: test_tag_1, key: test_key_2
发送成功 test_body_3
发送成功 test_body_4
【Consumer 收到消息】msgId: AC15000500346F296A464E8EF1290004, body: test_body_4, tag: test_tag_1, key: test_key_4
发送成功 test_body_5
发送成功 test_body_6
【Consumer 收到消息】msgId: AC15000500346F296A464E8EF5180006, body: test_body_6, tag: test_tag_1, key: test_key_6
发送成功 test_body_7
发送成功 test_body_8
【Consumer 收到消息】msgId: AC15000500346F296A464E8EF9070008, body: test_body_8, tag: test_tag_1, key: test_key_8
发送成功 test_body_9
Producer 发送完毕退出
```

输出可行，消费者只消费了 `test_tag_1` 标签的消息。