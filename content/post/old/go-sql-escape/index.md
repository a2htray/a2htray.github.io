---
title: "SQL 转义问题"
date: 2022-05-11T13:49:39+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - sqlx
 - MySQL
images:
 - images/word-image-16.png
---

SQL 转义问题是指执行的 SQL 语句中包含了某些特定的字符，如单引号 `'`、反斜杠 `\` 等，导致 SQL 语句不能正常执行。所以，我们应该在拼接 SQL 语句的过程中对特别的传入参数进行转义。

环境信息：

1. MySQL 8.0.28；
2. Go 1.16.9 windows/amd64

<!--more-->

## 问题发生

创建一个示例表 `test`，然后执行多条 SQL 语句：

```mysql
CREATE TABLE IF NOT EXISTS test (
    content text
) ENGINE = InnoDB;
```

```go
package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

func main() {
	dsn := "root:password$@tcp(127.0.0.1:3306)/database"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	// 元素的第 1 个反斜杠负责 Go 中字符串的转义
	contentSet := []string{
		"content1'",
		"content2\"",
		"content3\\",
	}
	for _, content := range contentSet {
		insertStmt := "INSERT INTO test VALUES ('" + content + "');"
		fmt.Println("insert statement:", insertStmt)
		_, err := db.Exec(insertStmt)
		if err != nil {
			log.Println(err)
		}
	}
}
```

```bash
insert statement: INSERT INTO test VALUES ('content1'');
2022/05/11 20:27:36 Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ''content1'')' at line 1
insert statement: INSERT INTO test VALUES ('content2"');
insert statement: INSERT INTO test VALUES ('content3\');
2022/05/11 20:27:37 Error 1064: You have an error in your SQL syntax; check the manual that corresponds to your MySQL server version for the right syntax to use near ''content3\')' at line 1
```

生成的 3 条 SQL 语句分别为：

```bash
INSERT INTO test VALUES ('content1'');
insert statement: INSERT INTO test VALUES ('content2"');
INSERT INTO test VALUES ('content3\');
```

显然只有第 2 条语句是合法的

从输出中得到，只有第 2 条语句是执行成功的，因为 `'content"'` 

## database/sql 的讨论

在 [issue 18478](https://github.com/golang/go/issues/18478) 中，开发者对 database/sql 包中是否加入转义进行了讨论。kardianos 认为转义功能并不是 database/sql 包最原始的需求，因为根据不同的数据库驱动应该有不同的转义实现。所以他建议开发者应该自行实现转义功能。时隔一年后，kardianos 做出了部分妥协，表示会在部分数据库驱动增加转义实现。

```bash
You can get the driver from sql.DB.Driver() so while it isn't a direct method from it in this shim, you can use in in a similar way right now.

Here's what I'd like to see:

1. evolve and verify the API that is in sqlexp,
2. get a few drivers to implement it (they usually appreciate PRs),
3. then if it makes sense (which I think it would at that point), bring it into the std lib.
Perhaps it would be good to combine the escaper functions with the database name function. We'd probably want to make it easy to expand it in the future as well, which would mean changing it into a struct with methods.
```

## 基本 args 传入

现修改上述代码，内容如下：

```go
package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

func main() {
	dsn := "root:password$@tcp(127.0.0.1:3306)/database"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	contentSet := []string{
		"content1'",
		"content2\"",
		"content3\\",
	}
	for _, content := range contentSet {
		_, err := db.Exec("INSERT INTO test VALUES (?);", content)
		if err != nil {
			log.Println(err)
		}
	}
}
```

数据可以正常导入。

```bash
mysql> SELECT * FROM test;
+-----------+
| content   |
+-----------+
| content2" |
| content1' |
| content2" |
| content3\ |
+-----------+
4 rows in set (0.13 sec)
```

## 自定义转义功能

如果实在需要自己去拼接 SQL 语句，则需要在拼接前对其进行转义

抄一个 issue 里的转义函数作为备用。

```go
func Escape(source string) string {
	var j int = 0
	if len(source) == 0 {
		return ""
	}
	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)
	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte
		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}
		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}
	return string(desc[0:j])
}
```

再次对程序进行修改，修改后内容如下：

```bash
func main() {
	dsn := "root:bgird123$@tcp(127.0.0.1:3306)/minority"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	contentSet := []string{
		"content1'",
		"content2\"",
		"content3\\",
	}
	for _, content := range contentSet {
		insertStmt := "INSERT INTO test VALUES ('" + Escape(content) + "');"
		fmt.Println("insert statement:", insertStmt)
		_, err := db.Exec(insertStmt)
		if err != nil {
			log.Println(err)
		}
	}
}
```

```bash
insert statement: INSERT INTO test VALUES ('content1\'');
insert statement: INSERT INTO test VALUES ('content2\"');
insert statement: INSERT INTO test VALUES ('content3\\');
```

数据可以正常插入。

## 总结

1. 手动拼接 SQL 语句时一定要进行转义；
2. database/sql 内置了 MySQL 驱动的转义功能，其它的没试；

