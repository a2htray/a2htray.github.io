+++
date = '2026-10-10T10:39:47+08:00'
draft = false
title = 'PostgreSQL：查询结果输出到文件'
categories = ['后端技术', 'PostgreSQL']
tags = ['PostgreSQL', '查询持久化', '数据库']
toc = true
+++

![](/imgs/learn-db/ScreenShot_2026-10-10_110556_605.png)

> 起因是今天在开发过程中需要将 PG 的查询结果输出到文件，方便后续分析。
> 问题已经解决，使用的 `\copy` 命令，然后跟 Buddy 也沟通了下，它还列举了其他一些方法，所以在这篇文章中做个小结。
> 本文内容大部分由 Buddy 生成，本人则对脚本、结论进行验证，已确认有效。


![](/imgs/learn-db/ScreenShot_2026-10-09_164450_760.png)

## 认识一个区别

`\copy` 和 `COPY`，它们不一样。

- **客户端写文件**（`\copy`、psql 重定向）：文件落在**连库的那台机器**上，普通账号就能用。
- **服务端写文件**（`COPY ... TO`）：文件落在**数据库服务器**的磁盘上，路径是服务器的路径，需要服务端写权限（`superuser` 或 `pg_write_server_files` 角色）。

区别：文件的落盘位置不同。

## 方式一：`\copy`

psql 的**元命令**，文件写在**客户端**机器上。

```sql
-- 导出为 CSV，带表头
\copy (SELECT id, name, created_at FROM users WHERE created_at > '2026-01-01')
  TO '/tmp/users.csv' WITH (FORMAT CSV, HEADER);

-- 导出纯文本（默认 tab 分隔）
\copy (SELECT * FROM orders LIMIT 100) TO '/tmp/orders.txt';
```

- `\copy` 是 **psql 命令，不是 SQL**，不能在应用程序的驱动里调用，只能在 psql 交互或 `-f` 脚本里用。
- 路径 `/tmp/users.csv` 是**本地机器**的路径，不需要任何服务端权限，普通账号即可。

## 方式二：`COPY ... TO`

真正的 SQL 命令，文件写在**数据库服务器**。

```sql
-- 服务端 CSV 导出（需 superuser 或 pg_write_server_files）
COPY (
  SELECT id, name, created_at
  FROM users WHERE created_at > '2026-01-01'
) TO '/var/lib/postgresql/exports/users.csv' WITH (FORMAT CSV, HEADER);

-- 直接导出整张表
COPY users TO '/var/lib/postgresql/exports/users.csv' WITH (FORMAT CSV, HEADER);
```

- 这是**服务端落盘**，路径必须是数据库进程能访问的服务器目录。普通账号会报 `permission denied`。
- 速度通常比 `\copy` 快，因为少了客户端 <-> 服务器之间的网络传输。
- 安全提醒：**不要在应用代码里用 `COPY TO` 写到不可控路径**，存在文件越权与覆盖风险。

## 方式三：psql 命令行重定向

不用进交互，一条命令搞定，文件落在**执行命令的机器**（客户端）。

```bash
# -c 执行查询，-o 指定输出文件（默认带 psql 表格边框，不建议直接当数据用）
psql -h localhost -U app -d mydb -c "SELECT id, name FROM users" -o /tmp/out.txt

# 想要干净的 CSV，配合 psql 的元组模式
psql -h localhost -U app -d mydb -At -F',' \
  -c "SELECT id, name FROM users" > /tmp/out.csv
```

- `-At`：`-A` 关掉表格边框，`-t` 去掉表头行，输出纯数据，适合机器消费。
- `-F','`：指定字段分隔符（Field separator）。

## 方式四：`\o` 与 `\g`

在 psql 交互里，把后续所有结果（或某一次结果）重定向到文件。

```sql
-- \o 打开输出重定向，之后每条查询的结果都写进文件
\o /tmp/result.txt
SELECT * FROM orders LIMIT 10;
SELECT count(*) FROM users;
\o            -- 再敲一次 \o 关闭重定向，恢复打印到屏幕

-- \g 只把当前这条查询的结果写到文件，不影响后续
SELECT * FROM orders LIMIT 10 \g /tmp/one_query.txt
```

- `\o` 是一直生效的开关，**用完记得再敲一次 `\o` 关掉**，否则你后面的查询全写进文件、屏幕啥也看不到。
- `\g filename` 只影响当前这一条，更安全，适合"就导出这一条"。
- 适合交互排查时顺手存一份，不适合放进自动化脚本。

## 方式五：`COPY ... PROGRAM`

服务端 `COPY` 还能把查询结果通过管道交给 shell 命令，比如导出时直接压缩、直接上传。

```sql
-- 导出时顺便 gzip 压缩（文件在服务器）
COPY (SELECT * FROM big_table)
TO PROGRAM 'gzip > /var/lib/postgresql/exports/big_table.csv.gz'
WITH (FORMAT CSV, HEADER);
```

要点：

- `PROGRAM` 在**服务端**执行，权限要求更高，且有**命令注入风险**，务必用白名单参数，绝不可拼接用户输入。
- 适合超大表导出时减少磁盘占用。

## 对比表格

| 方式                       | 文件落在 | 所需权限                                | 适合场景          |
| ------------------------ | ---- | ----------------------------------- | ------------- |
| `\copy`                  | 客户端  | 普通账号                                | 日常手动导出、开发自取   |
| `COPY ... TO`            | 服务端  | superuser / `pg_write_server_files` | 大批量高性能导出      |
| `psql -c / -o`、shell `>` | 客户端  | 普通账号                                | 脚本、crontab、CI |
| `\o` / `\g`              | 客户端  | 普通账号                                | psql 交互顺手存    |
| `COPY ... PROGRAM`       | 服务端  | superuser + 谨慎                      | 导出即压缩 / 加工    |

## 其它格式

* JSON
* 二进制

```sql
-- 1) JSON（PG 9.3+ 的 row_to_json / json_agg）
COPY (
  SELECT json_agg(row_to_json(t))
  FROM (SELECT id, name FROM users LIMIT 100) t
) TO '/tmp/users.json';

-- 2) 二进制（最高效，但只有 PG 自己能读回来，常用于跨库迁移）
COPY users TO '/tmp/users.bin' WITH (FORMAT BINARY);

-- 3) 客户端 \copy 同样支持这些格式
\copy (
  SELECT json_agg(row_to_json(t))
  FROM (SELECT id, name FROM users LIMIT 100) t
) TO '/tmp/users.json';
```
