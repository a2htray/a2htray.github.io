---
title: "Redis 的 list 和 stream：异步记录请求信息"
date: 2022-04-23T19:05:01+08:00
draft: false
reward: false
categories:
 - 后端技术
 - Redis
tags:
 - Redis
 - logging
 - Gin
 - Go
images:
 - images/managing_activity_logs.jpg
---

在 Web 开发中，常常需要对请求信息进行记录，形成日志以便于后期评估应用的性能。请求信息通常包含客户端地址、请求的 URL、请求时间及请求执行时间。在程序中，可以以同步或异步的方式完成这一需求。同步方式是指请求信息写入日志文件后才返回数据给客户端，异步方式则是在返回数据之前以新线程或进程完成对请求信息的记录。开源的日志包有：

1. [Zap](https://pkg.go.dev/go.uber.org/zap)：出自 Uber 团队，以高性能著称；
2. [Zerolog](https://github.com/rs/zerolog)：以易用性著称，支持 7 种日志级别；
3. [Logrus](https://github.com/sirupsen/logrus)：兼容标准日志包格式，也是本人常用的日志包；
4. [apex/log](https://github.com/apex/log)：受 Logrus 启发，简化操作后的 Logrus；
5. [ Log15](https://github.com/inconshreveable/log15)：日志可读性强；

5 个日志包的详细介绍可以看[《5 种结构化 Go 日志包对比分析》](/posts/go-five-structured-logging-package/)这篇文章。

<!--more-->

在 Go 开发中，一个非常简单的办法就是启用一个 goroutine 将请求信息发送到目的地，目的地可以是（1）一个日志文件；（2）一个 channel 或其它的应用（如 Redis）。在第 2 种方法中，还需要另一个拉取日志信息的服务，这类方法的优势是可以提高主体应用的性能，缺点是增加了系统的复杂度。本文的重点落在两个方面，分别为：

1. 解析请求，将信息发送到 Redis 服务器；
2. 读取 Redis 服务器中的请求信息，持久化到日志文件；

所以在本次实现中，包含两个组件（app 组件和 micro-dumper 组件）分别完成上述两项功能。

## app

首先需要一个 `Record` 类型来描述请求信息，其结构如下：

```go
// Record Represent a set of a Request passing from the client
type Record struct {
	RemoteAddr    string
	URL           string
	AccessTime    int64
	TimeExecuted  int64
	BodyBytesSent int64
}
```

在 Redis 端，我们也需要两种数据结构，分别为 Stream 和 List。Stream 可以用来保存请求的信息，而 List 则是用保存 Stream 中请求信息的 ID。

> 当前 Redis 不支持以位置索引的方式访问 Stream 中的信息。[Stackoverflow](https://stackoverflow.com/questions/63845306/how-do-i-get-records-in-a-redis-stream-by-position-instead-of-id)

接着定义两种数据类型的键名：

```go
const StreamKey = "api-request-log"
const RecordIDsKey = "api-request-record-ids"
```

然后就是实现 Redis 的连接，用于访问 Redis 服务器。

```go
var client *redis.Client
var once sync.Once

// Client return a redis client
func Client() *redis.Client {
	once.Do(func() {
		// Maybe you should instantiate redis client by reading config file
		client = redis.NewClient(&redis.Options{
			Network: "tcp",
			Addr:    "localhost:6379",
			DB:      0,
		})
	})
	return client
}
```

这里使用单例设计模式，应用只需要维护一个 Redis 客户端即可。有了记录（请求信息）和 Redis 客户端，就应该将记录发送到 Redis 服务器。

```go
// SendRecord send a record to redis server, two things are done as follows:
// 1. add a entry to the stream
// 2. push a entry id to the list
func SendRecord(record Record) {
	// add a entry by call .XAdd method
	xaddCmd := Client().XAdd(context.Background(), &redis.XAddArgs{
		Stream: StreamKey,
		ID:     "*",
		Values: map[string]interface{}{
			"remote_addr":     record.RemoteAddr,
			"url":             record.URL,
			"access_time":     record.AccessTime,
			"time_executed":   record.TimeExecuted,
			"body_bytes_sent": record.BodyBytesSent,
		},
	})
	if xaddCmd.Err() != nil {
		panic(xaddCmd.Err())
	}
	recordID := xaddCmd.Val()
	// push the id to the list by call .LPush method
	lpushCmd := Client().LPush(context.Background(), RecordIDsKey, recordID)
	if lpushCmd.Err() != nil {
		panic(lpushCmd.Err())
	}
}
```

上述代码做了两件事情：

1. 将记录添加到键为 `StreamKey` 的 Stream；
2. 将记录 ID 添加到键为 `RecordIDsKey` 的 List；

最后，app 组件的主程序如下：

```go
import (
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
	"redislog"
	"time"
)

type User struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// dataFromDB mock retrieving data from database
func dataFromDB() []*User {
	return []*User{
		{
			"xiaoming",
			12,
		},
		{
			"xiaohong",
			13,
		},
		{
			"xiaobei",
			14,
		},
	}
}

// RedisLogger a gin.HandlerFunc wrapper
// extract request information, assemble to a record, and send to Redis server via goroutine
func RedisLogger(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		record := redislog.Record{
			RemoteAddr: c.Request.RemoteAddr,
			URL:        c.Request.URL.RequestURI(),
			AccessTime: time.Now().Unix(),
		}
		f(c)
		record.TimeExecuted = time.Now().Unix() - record.AccessTime
		record.BodyBytesSent = int64(c.Writer.Size())
		go redislog.SendRecord(record)
	}
}

func main() {
	engine := gin.Default()
	engine.GET("api/users", RedisLogger(func(c *gin.Context) {
		time.Sleep(time.Second * time.Duration(rand.Int63n(5)))
		c.JSON(http.StatusOK, gin.H{
			"data": dataFromDB(),
		})
	}))
	engine.GET("api/user/:name/age", RedisLogger(func(c *gin.Context) {
		users := dataFromDB()
		name := c.Param("name")
		var user *User
		for _, u := range users {
			if u.Name == name {
				user = u
			}
		}
		var age interface{}
		if user == nil {
			age = "unknown"
		} else {
			age = user.Age
		}
		c.JSON(http.StatusOK, gin.H{
			"age": age,
		})
	}))

	log.Fatalln(engine.Run(":9090"))
}
```

上述代码定义了两个接口：/api/users 和 /api/user/:name/age，两个 gin.HandlerFunc 都被 RedisLogger 进行“包裹”。

## mirco-dumper

micro-dumper 会周期性拉取 Stream 里的记录，如果存在并保存到日志文件，反之则等待下一刻执行。主要实现其实就一个 `ReadRecord` 方法：

```go
// ReadRecord read a record from redis, three things are done as follows:
// 1. retrieve a entry id from the list
// 2. retrieve a entry from the stream via the entry id
// 3. after retrieving the entry, delete the entry from the stream
func ReadRecord() (Record, bool) {
	// retrieve record id from the redis list
	lpopCmd := Client().LPop(context.Background(), RecordIDsKey)
	recordID := lpopCmd.Val()
	if recordID == "" {
		return Record{}, false
	}
	// read the record from the stream
	xreadCmd := Client().XRead(context.Background(), &redis.XReadArgs{
		Streams: []string{StreamKey, recordID},
		Count:   1,
		Block:   0,
	})
	if xreadCmd.Err() != nil {
		panic(xreadCmd.Err())
	}
	// if read successfully, we should remove record from the stream
	xdelCmd := Client().XDel(context.Background(), StreamKey, recordID)
	if xdelCmd.Err() != nil {
		panic(xdelCmd.Err())
	}
	record := xreadCmd.Val()[0].Messages[0].Values

	accessTime, _ := strconv.ParseInt(record["access_time"].(string), 10, 64)
	timeExecuted, _ := strconv.ParseInt(record["time_executed"].(string), 10, 64)
	bodyBytesSent, _ := strconv.ParseInt(record["body_bytes_sent"].(string), 10, 64)

	return Record{
		RemoteAddr:    record["remote_addr"].(string),
		URL:           record["url"].(string),
		AccessTime:    accessTime,
		TimeExecuted:  timeExecuted,
		BodyBytesSent: bodyBytesSent,
	}, true
}
```

读记录的逻辑如下：

1. 从 List 中得到记录的 ID；
2. 通过 ID 得到 Stream 中的记录；
3. 记录获取完毕后，将该记录在 Stream 中删除；

> LPOP 是左端弹出操作，在获取的同时，List 中已经不存在该 ID 了。

最后编写主程序：

```go
func main() {
	fmt.Println("start to retrieve request record ...")
	f, _ := os.Create("./request.log")
	for {
		time.Sleep(1)
		if record, found := redislog.ReadRecord(); found {
			_, err := f.WriteString(fmt.Sprintf("remote addr: %s url: %s access time: %d time executed: %d body bytes sent: %d\n", record.RemoteAddr, record.URL, record.AccessTime, record.TimeExecuted, record.BodyBytesSent))
			if err != nil {
				panic(err)
			} else {
				fmt.Println("write a request record to log file")
			}
		}
	}
}
```

## 测试

运行 app 和 mirco-dumper：

```bash
$ nohup go run app/main.go > app.out &
$ nohup go run micro-dumper/main.go > micro-dumper.out &
```

使用 ab 测试工具发送大量请求：

```bash
$ ab -n 1000 -c 10 http://localhost:9090/api/users
$ ab -n 1000 -c 10 http://localhost:9090/api/user/xiaohong/age
```

查看日志：

```bash
head -10 request.log
remote addr: [::1]:58644 url: /api/users access time: 1650716423 time executed: 0 body bytes sent: 96
remote addr: [::1]:58636 url: /api/users access time: 1650716423 time executed: 1 body bytes sent: 96
remote addr: [::1]:58666 url: /api/users access time: 1650716424 time executed: 0 body bytes sent: 96
remote addr: [::1]:58640 url: /api/users access time: 1650716423 time executed: 2 body bytes sent: 96
remote addr: [::1]:58658 url: /api/users access time: 1650716423 time executed: 2 body bytes sent: 96
remote addr: [::1]:58680 url: /api/users access time: 1650716425 time executed: 0 body bytes sent: 96
remote addr: [::1]:58648 url: /api/users access time: 1650716423 time executed: 3 body bytes sent: 96
remote addr: [::1]:58654 url: /api/users access time: 1650716423 time executed: 3 body bytes sent: 96
remote addr: [::1]:58682 url: /api/users access time: 1650716425 time executed: 1 body bytes sent: 96
remote addr: [::1]:58646 url: /api/users access time: 1650716423 time executed: 4 body bytes sent: 96
```

## 完整代码

[redislog](https://github.com/a2htray/redislog)

## 总结

1. `ResponseWriter` 接口定义了 `Size` 方法可以得到响应体中的字节数；
2. 因为 Redis 保存的都是字符串形式，所以在 Go 代码中总是要做字符串转换；

