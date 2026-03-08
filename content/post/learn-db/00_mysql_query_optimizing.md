+++
date = '2026-03-06T09:18:58+08:00'
draft = false
title = 'MySQL 查询优化技巧'
categories = ['后端技术', 'MySQL']
tags = ['MySQL', '查询优化']
toc = true
+++

基础查询、关联查询、排序与分组、特殊场景、索引设计

## 基础查询优化

**🎯 目标：规避全表扫描、减少网络 / 内存开销，让查询更轻量化**

1. <font color="blue">只查所需字段，拒绝 SELECT *</font>

❌ SELECT * FROM users;

✅ SELECT id, name, phone FROM users;

2. <font color="blue">禁止在索引列做运算 / 函数操作</font>

❌ WHERE YEAR(create_time) = 2026;

✅ WHERE create_time >= '2026-01-01' AND create_time < '2027-01-01';

3. <font color="blue">优化 LIKE 查询，避免前缀模糊</font>

❌ WHERE name LIKE '%张%';

✅ WHERE name LIKE '张%';

4. <font color="blue">优化 OR 查询，拆分为 UNION/UNION ALL</font>

❌ SELECT * FROM users WHERE name = '张三' OR age = 25;

✅ SELECT * FROM users WHERE name = '张三' UNION ALL SELECT * FROM users WHERE age = 25;

> 已验证，在 MySQL 5.5 及以前版本中，该优化点有效，见 [MySQL OR 查询优化为 UNION/UNION AL 验证](/post/learn-db/01_or_operation_testing/)。

5. <font color="blue">优化 IN 查询，控制列表值数量</font>

❌ IN 后跟 2000 个值 

✅ 将 2000 个值拆分为 2 个 IN 查询，各含 1000 个值，或用 JOIN 替代

6. <font color="blue">LIMIT 必搭配 ORDER BY，避免结果无序</font>

❌ SELECT * FROM table LIMIT 10; 

✅ SELECT * FROM table ORDER BY id LIMIT 10;

## 关联查询优化

**🎯 目标：减少嵌套循环次数、利用索引、避免临时表，提升多表关联效率**

1. <font color="blue">小表驱动大表，降低循环次数</font>

A 表 1000 条记录，B 表 100000 条记录

❌ B JOIN A

✅ A JOIN B

2. <font color="blue">JOIN 连接字段必须建立索引</font>

SELECT o.* FROM users u JOIN orders o ON u.id = o.user_id;

❌ orders.user_id 没有索引

✅ orders.user_id 有索引

3. <font color="blue">明确 ON 与 WHERE 区别，合理放置过滤条件</font>

LEFT JOIN 示例：保留左表所有数据则条件写 ON。

❌ SELECT u.*, o.order_no FROM users u LEFT JOIN orders o ON u.id = o.user_id WHERE o.status = 1;

✅ SELECT u.*, o.order_no FROM users u LEFT JOIN orders o ON u.id = o.user_id AND o.status = 1;

4. <font color="blue">用 JOIN 替代子查询，避免创建临时表</font>

❌ SELECT * FROM orders WHERE user_id IN (SELECT id FROM users WHERE age > 20);

✅ SELECT o.* FROM orders o JOIN users u ON o.user_id = u.id WHERE u.age > 20;

5. <font color="blue">限制关联表数量，不超过 3 个</font>

❌ JOIN 超过 3 个

✅ JOIN 不超过 3 个

## 排序与分组优化

**🎯 目标：利用索引避免文件排序、临时表**

1. <font color="blue">优化 GROUP BY，关闭默认无用排序</font>

无需排序时

❌ SELECT user_id, COUNT(*) FROM orders WHERE status=1 GROUP BY user_id;

✅ SELECT user_id, COUNT(*) FROM orders WHERE status=1 GROUP BY user_id ORDER BY NULL;

> TODO 搭建环境，测试下这一条优化点。

2. <font color="blue">给 GROUP BY 字段建索引，避免临时表</font>

SELECT user_id, COUNT(*) FROM orders GROUP BY user_id;

❌ orders.user_id 不建索引

✅ orders.user_id 建索引

3. <font color="blue">索引覆盖 ORDER BY，避免文件排序</font>

SELECT id FROM orders WHERE status=1 ORDER BY create_time;

❌ 不建联合索引 idx_status_create_time(status, create_time)

✅ 建联合索引 idx_status_create_time(status, create_time)

## 特殊场景优化

**🎯 目标：保证索引正常生效、减少额外开销**

1. <font color="blue">避免隐式类型转换，防止索引失效</font>

user_id 为 INT 类型 

❌ WHERE user_id = '100';

✅ WHERE user_id = 100;

2. <font color="blue">优化范围查询，联合索引中范围字段放最后</font>

WHERE status=1 AND create_time > '2026-01-01';

❌ 索引idx_status_create_time(create_time, status)

✅ 索引idx_status_create_time(status, create_time)

3. <font color="blue">优化 NULL 值查询，用默认值替代 NULL</font>

❌ SELECT * FROM orders WHERE remark IS NULL;

✅ SELECT * FROM orders WHERE remark = '';

> 建表的时候，设置字段的 DEFAULT 值。

4. <font color="blue">批量插入，减少事务 / 网络开销</font>

❌ 单条插入 2 次

✅ INSERT INTO orders (user_id, amount) VALUES (1, 100), (2, 200);

5. <font color="blue">批量更新，用 CASE WHEN 替代多次 UPDATE</font>

❌ 2 次单条更新

✅ UPDATE orders SET amount = CASE WHEN order_id=1 THEN 150 WHEN order_id=2 THEN 250 END WHERE order_id IN (1,2);

6. <font color="blue">优先用 GROUP BY 替代 SELECT DISTINCT 去重</font>

❌ SELECT DISTINCT user_id FROM orders WHERE status=1;

✅ SELECT user_id FROM orders WHERE status=1 GROUP BY user_id ORDER BY NULL;

7. <font color="blue">避免 IN 子句包含子查询，改写为 JOIN</font>

❌ SELECT * FROM A WHERE id IN (SELECT id FROM B WHERE status = 1);

✅ SELECT A.* FROM A INNER JOIN B ON A.id = B.id WHERE B.status = 1;

8. <font color="blue">避免 NOT IN 子查询，用 NOT EXISTS/LEFT JOIN 替代</font>

查询「没有下过已支付订单」的所有用户。

❌

```sql
SELECT * 
FROM users 
WHERE id NOT IN (
    SELECT user_id 
    FROM `order` 
    WHERE status = 1
);
```

✅

```sql
SELECT u.* 
FROM users u
LEFT JOIN `order` o 
    ON u.id = o.user_id
    AND o.status = 1
WHERE o.user_id IS NULL;
```

> TODO 搭建环境，测试下这一条优化点。

9. <font color="blue">正确使用 COUNT 函数，优先用 COUNT (*)</font>

❌ COUNT(remark) 

✅ SELECT COUNT(*) FROM orders WHERE status=1;

## 索引设计优化

**🎯 目标：平衡查询效率与写入维护成本**

1. <font color="blue">用联合索引替代多个单列索引</font>

❌ 单独建 status、create_time 索引

✅ 建联合索引idx_status_create_time(status, create_time)

2. <font color="blue">区分度高的列优先建索引</font>

❌ 单独给性别（男 / 女）、状态（0/1）建索引（可纳入联合索引末尾）

✅ 给 id、phone、username 建索引

3. <font color="blue">长字符串字段建**前缀索引**，减小索引体积</font>

给 url 字段建前 20 个字符的前缀索引

❌ CREATE INDEX idx_url_prefix ON table_name(url);

✅ CREATE INDEX idx_url_prefix ON table_name(url(20));

4. <font color="blue">选择精简的数据类型，字段尽量设为 NOT NUL</font>

❌ 用 BIGINT 存性别（仅 0/1）

✅用 TINYINT (1)；

❌ 字段允许 NULL

✅ 设默认值VARCHAR(50) DEFAULT ''、INT DEFAULT 0

## 资源

* [MySQL 查询优化 30 条封神技巧：用好索引，少耗资源，查询快到飞起 ](https://juejin.cn/post/7613361194158440494)