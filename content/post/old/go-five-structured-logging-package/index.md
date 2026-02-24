---
title: "5 种结构化 Go 日志包对比分析"
date: 2022-04-24T12:48:00+08:00
draft: true
reward: false
categories:
  - 后端技术
  - Go
tags:
 - logging
 - 翻译
---

本文是一篇翻译文章，正文为**参考 [1]**。

一般来说，日志就是把记录输出到终端或者写入文件，听起来这是一件非常简单的任务。但是记录日志也有它的最佳实践，比如你必须规定日志的级别、结构化日志信息、多地存储日志以及日志的上下文信息等。将这些需求结合在一起，那么记录日志其实是一件非常难的事情。

<!--more-->

结构化日志的一个隐含思想是为了以后可以更好地处理这些信息，通常会使用 JSON 格式进行存储日志信息。比如你知道一个用户 ID，想知道与这个用户相关的日志信息，就可以通过用户 ID 对日志信息进行检索。当日志记录是结构化的，你可以从日志中得到其它相关的信息。

## 为什么要使用结构化日志包

选择结构化日志包通常有以下几个原因：

1. Go 内置的 `log` 包的输出不易于理解，追踪起来费时费力；
2. 结构化日志包可以增加任意的字段信息，从而更好地检索及 DEBUG；
3. 格式化之后的信息更容易理解及存储；

## 内置 log 包的适用性

Go 内置的日志包不需要任何配置就可以使用，非常方法在本地调试。同时它也支持定制化日志信息，但是缺少了日志级别以及对多种格式的支持。在这篇文章中，我们会介绍和测试 5 种结构化日志包，便于开发者进行选型。

## Zap

[Zap](https://pkg.go.dev/go.uber.org/zap) 是非常流行的 Go 的第三方日志包，由 Uber 团队开发并开源。相比于其它日志包，Zap 可以高性能地记录日志信息，而且其速度要超过 Go 内置的日志包。

Zap 提供了两种日志输出，包括以高性能著称的 `Logger` 和以扩展性著称的 `SugaredLogger`。但其实两种类型的处理速度都比较快。

在下面的例子中，我们使用了一个 `zap.SugaredLogger` 的实例以 JSON 的方式记录日志，格式包括日志级别、时间戳、文件名、行号以及日志消息。

```go
import (
    "log"

    "go.uber.org/zap"
)

func main() {
    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatal(err)
    }

    sugar := logger.Sugar()

    sugar.Info("Hello from zap logger")
}

// Output:
// {"level":"info","ts":1639847245.7665887,"caller":"go-logging/main.go:21","msg":"Hello from zap logger"}
```

通过修改日志配置，你可以修改 JSON 中的字段信息及其值显示。如果将 JSON 中的 `ts` 修改为 `timestamp`，或者改变时间的字符串表示。

```go
func main() {
    loggerConfig := zap.NewProductionConfig()
    loggerConfig.EncoderConfig.TimeKey = "timestamp"
    loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

    logger, err := loggerConfig.Build()
    if err != nil {
        log.Fatal(err)
    }

    sugar := logger.Sugar()

    sugar.Info("Hello from zap logger")
}

// Output:
// {"level":"info","timestamp":"2021-12-18T18:21:34+01:00","caller":"go-logging/main.go:23","msg":"Hello from zap logger"}
```

如果你更关心应用的性能，可以在 `SugaredLogger` 实例上调用 `DeSugar` 方法，这样就可以返回一个标准的、更快的 `Logger`。然而这么做后，你需要指定要打印的具体类型。

```go
l := sugar.Desugar()

l.Info("Hello from zap logger",
  zap.String("tag", "hello_zap"),
  zap.Int("count", 10),
)
```

## Zerolog

[Zerolog](https://github.com/rs/zerolog) 是功能完备的、用于记录 JSON 格式日志信息的日志包，设计理念以性能、易用性为主。你可以使用默认的 logger，它的用法相当简便。如果你要使用这个 logger，需要在代码中引入 log 的子包，就像下面这样：

```go
package main

import (
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    zerolog.SetGlobalLevel(zerolog.InfoLevel)

    log.Error().Msg("Error message")
    log.Warn().Msg("Warning message")
    log.Info().Msg("Info message")
    log.Debug().Msg("Debug message")
    log.Trace().Msg("Trace message")
}

// Output:
// {"level":"error","time":"2021-12-19T17:38:12+01:00","message":"Error message"}
// {"level":"warn","time":"2021-12-19T17:38:12+01:00","message":"Warning message"}
// {"level":"info","time":"2021-12-19T17:38:12+01:00","message":"Info message"}
```

Zerolog 支持 7 种日志级别，从用于追踪的 `trace` 至抛出异常的 `panic`。你可以使用 `SetGlobalLevel()` 方法设置全局 logger 的日志级别。在上述的例子中，日志级别为 `info`，所以只有日志级别大于 `info` 的日志信息才会被写入。

Zerolog 支持上下文日志记录，通过 `zerolog.Event` 类型上的方法，你可以很容易地住日志记录中加入额外的字段。

在一个 logger 上调用任一方法，如 `Error()`，最后以 `Msg()` 或 `Msgf()` 结束，都会创建一个 `zerolog.Event` 实例。在下面的例子中，我们加入了一串上下文：

```go
log.Info().Str("tag", "a tag").Int("count", 123456).Msg("info message")

// Output:
// {"level":"info","tag":"a tag","count":123456,"time":"2021-12-20T09:01:33+01:00","message":"info message"}
```

`Event` 类型上有一个特殊的 `Err()` 方法，可以用来传递 `error` 类型的数据。如果希望修改 `error` 值对应的字段名，可以设置 `zerolog.ErrorFieldName` 的值。

```go
err := fmt.Errorf("An error occurred")

log.Error().Err(err).Int("count", 123456).Msg("error message")

// Output:
// {"level":"error","error":"An error occurred","count":123456,"time":"2021-12-20T09:07:08+01:00","message":"error message"}
```



## Logrus



## apex/log



## log15



## 性能比较





## 参考

1. [5 structured logging packages for Go](https://blog.logrocket.com/five-structured-logging-packages-for-go/)

   