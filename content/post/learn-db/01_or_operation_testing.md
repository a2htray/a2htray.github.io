+++
date = '2026-03-08T18:41:08+08:00'
draft = false
title = 'MySQL OR 查询优化为 UNION/UNION AL 验证'
categories = ['后端技术', 'MySQL']
tags = ['MySQL', '查询优化', 'OR 查询', 'UNION 操作']
toc = true
+++

验证目标：在特定场景下（OR 连接不同字段的条件），将 OR 查询拆分为 UNION ALL 能显著提升查询效率。

## 通过 Docker 运行 MySQL 服务

```bash
$ docker pull mysql:8.0
$ docker run -d \
    --name mysql-or-testing \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=123456 \
    -e MYSQL_CHARACTER_SET_SERVER=utf8mb4 \
    -e MYSQL_COLLATION_SERVER=utf8mb4_unicode_ci \
    -v mysql_data:/var/lib/mysql \
    --restart=always \
    mysql:8.0
```

查看服务运行状态：

```bash
$ docker ps -a
CONTAINER ID   IMAGE       COMMAND                  CREATED         STATUS         PORTS                               NAMES
3ec8934547ce   mysql:8.0   "docker-entrypoint.s…"   7 seconds ago   Up 7 seconds   0.0.0.0:3306->3306/tcp, 33060/tcp   mysql-or-testing
```

## 创建测试表，插入大量数据

进入容器内的 MySQL 终端：

```bash
$ docker exec -it mysql-or-testing /bin/bash
bash-5.1# mysql -uroot -p123456
```

创建测试数据库

```mysql
CREATE DATABASE IF NOT EXISTS test_optimize; 
USE test_optimize;
```

创建用户行为表，模拟真实业务数据：

```mysql
CREATE TABLE user_behavior (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    user_id BIGINT NOT NULL,
    order_id BIGINT NOT NULL,
    create_time DATETIME NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

创建一个插入测试数据的存储过程：

```mysql
DELIMITER //
CREATE PROCEDURE insert_test_data()
BEGIN
    DECLARE i INT DEFAULT 1;
    WHILE i <= 1000 DO
        INSERT INTO user_behavior (user_id, order_id, create_time)
        SELECT 
            FLOOR(RAND() * 1000000),
            FLOOR(RAND() * 10000000),
            DATE_ADD('2026-01-01', INTERVAL FLOOR(RAND() * 365) DAY)
        FROM information_schema.tables LIMIT 1000;
        SET i = i + 1;
    END WHILE;
END //
DELIMITER ;
```

调用 `insert_test_data` 存储过程，并查看插入的记录数

```mysql
-- 多次调用，使其记录数尽可能多
CALL insert_test_data();
CALL insert_test_data();
CALL insert_test_data();

SELECT COUNT(*) FROM user_behavior;
```

共插入 990000 条记录。给 `user_behavior` 表的 `user_id` 和 `order_id` 创建索引：

```mysql
CREATE INDEX idx_user_id ON user_behavior(user_id);
CREATE INDEX idx_order_id ON user_behavior(order_id);
```

## 验证

逻辑：查询 user_id = 929195 或 order_id = 9804737 的记录。

```mysql
-- 1. 开启执行时间显示
SET profiling = 1;

-- 2. 执行OR查询
SELECT * FROM user_behavior WHERE user_id = 929195 OR order_id = 9804737;

-- 3. 查看执行计划
EXPLAIN
SELECT * FROM user_behavior WHERE user_id = 929195 OR order_id = 9804737\G;

-- 4. 查看查询耗时
SHOW PROFILES;
```

执行计划的输出：

```text
*************************** 1. row ***************************
           id: 1
  select_type: SIMPLE
        table: user_behavior
   partitions: NULL
         type: index_merge
possible_keys: idx_user_id,idx_order_id
          key: idx_user_id,idx_order_id
      key_len: 8,8
          ref: NULL
         rows: 4
     filtered: 100.00
        Extra: Using union(idx_user_id,idx_order_id); Using where
1 row in set, 1 warning (0.00 sec)
```

> 类型为索引合并（index_merge）。

执行时间为：0.00763325。

将查询语句调整为 UNION 方式：

```mysql
-- 1. 执行UNION ALL查询
SELECT * FROM user_behavior WHERE user_id = 929195
UNION ALL
SELECT * FROM user_behavior WHERE order_id = 9804737;

-- 2. 查看执行计划
EXPLAIN
SELECT * FROM user_behavior WHERE user_id = 929195
UNION ALL
SELECT * FROM user_behavior WHERE order_id = 9804737\G;

-- 3. 查看查询耗时（对比OR版本）
SHOW PROFILES;
```

执行计划的输出：

```mysql
*************************** 1. row ***************************
           id: 1
  select_type: PRIMARY
        table: user_behavior
   partitions: NULL
         type: ref
possible_keys: idx_user_id
          key: idx_user_id
      key_len: 8
          ref: const
         rows: 2
     filtered: 100.00
        Extra: NULL
*************************** 2. row ***************************
           id: 2
  select_type: UNION
        table: user_behavior
   partitions: NULL
         type: ref
possible_keys: idx_order_id
          key: idx_order_id
      key_len: 8
          ref: const
         rows: 2
     filtered: 100.00
        Extra: NULL
2 rows in set, 1 warning (0.00 sec)
```

> 执行计划显示有两个查询，类型是非唯一索引。

执行时间为：0.00937475s。

## 总结

**UNION/UNION ALL 操作反而略慢于 OR 查询**。

在 MySQL 8 及两个过滤字段都有索引的情况下，使用 UNION/UNION ALL 操作并不会明显提升查询速度。原因在于：

* MySQL 5.6 以前无法利用多个索引字段，而之后有“索引合并 index_merge”
* **优化器识别到 OR 条件可以拆分为两个独立索引的查询**，自动应用了“索引合并 index_merge”