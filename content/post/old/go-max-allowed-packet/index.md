---
title: "Go&MySQL`max_allowed_packet`"
date: 2022-04-30T20:25:09+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - MySQL
 - Go
---

当发送给数据库的语句过大时，会报如下错误：

```bash
panic: Error 1105: Parameter of prepared statement which is set through mysql_send_long_data() is longer than 'max_allowed_packet' bytes
```

从报错中可知，需要修改 `max_allowed_packet` 选项的值。

<!--more-->

## my.cnf

在 my.cnf 配置文件中修改该选项是最直接了当的，修改成合适的值后重启服务即可。

```bash
[mysqld]
max_allowed_packet=16M
```

## 会话

当打开一个会话时，也可以使用下面的语句设置全局或当前的 `max_allowed_packet` 值。

```bash
SET GLOBAL max_allowed_packet=1073741824;
SET max_allowed_packet=1073741824;
```

第 1 个在服务重启后失效，第 2 个在会话结束会失效。

## go-sql-driver/mysql

更多情况下，还是希望在程序中控制该选项，所以要看下使用的数据库驱动是否支持该选项。很庆幸，`go-sql-driver/mysql` 是支持的，使用的版本如下：

```bash
go-sql-driver/mysql@v1.6.0
```

在使用 `gorm` 中，返回一个 `*gorm.DB` 实例，我们都会使用下面的代码：

```go
db, err := gorm.Open(mysql.Open("username:password@tcp(host:port)/database?queryString"), &gorm.Config{})
```

下面说一下整行代码的执行顺序：

1. mysql.Open 简单返回一个实现 `gorm.Dialector` 接口的实例（只是简单的赋值，并没有作解析）；
2. gorm.Open 会调用 `Dialector.Initialize` 方法，并且会真正地打开数据库连接；

```go
db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)
```

3. 在 sql.Open 中会调用；

```go
connector, err := driverCtx.OpenConnector(dataSourceName)
```

4. 在 driverCtx.OpenConnector 中会调用 `ParseDSN` 方法；
5. ParseDSN 完成了对 DSN 的解析，包括 `queryString` 部分；

即：

```bash
mysql.Open
\
 \ gorm.Open -> Dialector.Initialize -> sql.Open -> driverCtx.OpenConnector
                                                 -> ParseDSN
```

`go-sql-driver/mysql` 支持的 query 参数包括如下：

```bash
allowAllFiles
allowCleartextPasswords
allowNativePasswords
allowOldPasswords
checkConnLiveness
clientFoundRows
collation
columnsWithAlias
compress # 这个是没有实现的
interpolateParams
loc
multiStatements
parseTime
readTimeout
rejectReadOnly
serverPubKey
strict
timeout
tls
writeTimeout
maxAllowedPacket
```

## 参考

1. [max_allowed_packet in mySQLmax_allowed_packet in mySQL](https://dba.stackexchange.com/questions/45087/max-allowed-packet-in-mysql)