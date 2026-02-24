---
title: "PostgreSQL JSON 类型字段"
date: 2022-04-13T10:51:09+08:00
draft: true
reward: false
tags:
 - "PostgreSQL"
 - "JSON 类型"
 - "数据库函数"
categories:
 - 后端技术
 - PostgreSQL
images:
 - "images/postgresql-column-type-json.webp"
---

JSON 类型是 PostgreSQL 9.2 加的数据类型，



环境：

**PostgreSQL**

```sql
SELECT version()
```

```bash
PostgreSQL 9.5.24, compiled by Visual C++ build 1800, 64-bit
```

## JSONB 类型

PostgreSQL 9.4 支持存储 JSON 的二进制格式（Binary JSON，JSONB）

### 索引

PostgreSQL 支持对 JSONB 类型字段创建索引，可以加速数据的查询，通过下方示例加以说明。

1. 首先创建一个包含 JSONB 类型字段的表；

```sql
CREATE TABLE test_jsonb_index
(
    student jsonb
);
```

2. 接着向表 `test_jsonb_index` 插入 10000 条数据。因为 PostgreSQL 9.5 不支持存储过程，所以这里使用函数实现；

```sql
CREATE OR REPLACE FUNCTION func_insert_100000_students()
RETURNS void
LANGUAGE plpgsql
AS $BODY$
    DECLARE name varchar(10);
    DECLARE age integer;
    BEGIN
        FOR i IN 1..100000 LOOP
            name = concat('user', i);
            SELECT floor(random()*(24 - 12+1) + 12) INTO age;
            INSERT INTO test_jsonb_index VALUES (concat('{"name": "', name, '", "age": ', age, '}')::jsonb);
        END LOOP;
    END
$BODY$;

SELECT func_insert_10000_students(); -- 插入数据
```

3. 在没有索引的情况下，完成两次查询操作；

**`age` 进行区间查询（无索引）**

```sql
EXPLAIN SELECT * FROM test_jsonb_index WHERE (student->>'age')::integer IN (15, 20);
```

```bash
Seq Scan on test_jsonb_index  (cost=0.00..2935.00 rows=1000 width=45)
  Filter: (((student ->> 'age'::text))::integer = ANY ('{15,20}'::integer[]))
```

**`age` 进行 GROUP 操作（无索引）**

```sql
EXPLAIN SELECT student->>'age', count(*) FROM test_jsonb_index GROUP BY student->>'age';
```

```plsql
GroupAggregate  (cost=13566.32..15566.32 rows=100000 width=45)
  Group Key: ((student ->> 'age'::text))
  ->  Sort  (cost=13566.32..13816.32 rows=100000 width=45)
        Sort Key: ((student ->> 'age'::text))
        ->  Seq Scan on test_jsonb_index  (cost=0.00..2185.00 rows=100000 width=45)
```

4. 为 `student->>'age'` 字段创建索引；

```sql
CREATE INDEX idx_on_age ON test_jsonb_index ((student->>'age'));
```

5. 创建索引后，再次执行两条 `EXPLAIN` 语句；

 **`age` 进行区间查询（有索引）**

```sql
EXPLAIN SELECT * FROM test_jsonb_index WHERE (student->>'age')::integer IN (15, 20);
```

```bash
Seq Scan on test_jsonb_index  (cost=0.00..2935.00 rows=1000 width=45)
  Filter: (((student ->> 'age'::text))::integer = ANY ('{15,20}'::integer[]))
```

**`age` 进行 GROUP 操作（有索引）**

```sql
EXPLAIN SELECT student->>'age', count(*) FROM test_jsonb_index GROUP BY student->>'age';
```

```bash
GroupAggregate  (cost=0.29..8344.29 rows=100000 width=45)
  Group Key: (student ->> 'age'::text)
  ->  Index Scan using idx_on_age on test_jsonb_index  (cost=0.29..6594.29 rows=100000 width=45)
```

对比前后两组查询，当为 `student->>'age'` 创建索引后，可知：

1. `age` 区间查询并没有利用到索引，还是做了全表扫描。这里因为 `WHERE (student->>'age')::integer` 对 `age` 字段做了类型转换；
2. `age` GROUP 操作有效地利用到了索引；

### 与 JSON 的区别



## 操作符

JSON 和 JSONB 数据类型分别存储了数据的文本和二进制格式，所以需要通过操作符来获取数据中信息。*表 1* 是 PostgreSQL 对 JSON 和 JSONB 数据类型支持的操作符。

*表 1 操作符*

| 操作符 | 右值类型描述 | 说明                                            | 示例                                               | 结果           |
| ------ | ------------ | ----------------------------------------------- | -------------------------------------------------- | -------------- |
| `->`   | `int`        | 获取数组 JSON 数据中元素，支持负数              | `'[{"a":"foo"},{"b":"bar"},{"c":"baz"}]'::json->2` | `{"c":"baz"}`  |
| `->`   | `text`       | 获取指定 key 的 JSON 对象                       | `'{"a": {"b":"foo"}}'::json->'a'`                  | `{"b":"foo"}`  |
| `->>`  | `int`        | 以**文本**的形式获取数组 JSON 数据中的元素      | `'[1,2,3]'::json->>2`                              | `3`            |
| `->>`  | `text`       | 以**文本**的形式获取指定 key 的 JSON 对象       | `'{"a":1,"b":2}'::json->>'b'`                      | `2`            |
| `#>`   | `text[]`     | 基于 key 路径匹配的方式获取指定路径的 JSON 对象 | `'{"a": {"b":{"c": "foo"}}}'::json#>'{a,b}'`       | `{"c": "foo"}` |
| `#>>`  | `text[]`     | 以文本的形式获取指定路径的 JSON 对象            | `'{"a":[1,2,3],"b":[4,5,6]}'::json#>>'{a,2}'`      | `3`            |

示例还包括：

```sql
-- 获取数组中第 1 个元素（0 计数）
SELECT '[
  "apple",
  "banana",
  "peach"
]'::json -> 1;
-- "banana"

-- 获取数组中第 1 个元素（0 计数）
SELECT '[
  1.5,
  2.5,
  3.5
]'::jsonb -> 1;
-- 2.5

-- 获取 JSON 对象中 fruit 对应的值
SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::json -> 'fruit';
-- "apple"

-- 获取 JSON 对象中 price 对应的值
SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::jsonb -> 'price';
-- 1.5

-- 取数组中第 1 个元素（0 计数）
SELECT '[
  "apple",
  "banana",
  "peach"
]'::json ->> 1;
-- banana

-- 取数组中第 1 个元素（0 计数）
SELECT '[
  1.5,
  2.5,
  3.5
]'::jsonb ->> 1;
-- 2.5

-- 获取 JSON 对象中 fruit 对应的值
SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::json ->> 'fruit';
-- apple

-- 获取 JSON 对象中 price 对应的值
SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::jsonb ->> 'price';
-- 1.5

-- 获取订单的水果数组
SELECT '{
  "order": [
    "apple",
    "banana",
    "peach"
  ]
}'::json #> '{order}';
-- ["apple", "banana", "peach"]

-- 获取订单数组中第 0 个商品的水果名称
SELECT '{
  "order": [
    {
      "fruit": "apple",
      "price": 1.5
    },
    {
      "fruit": "banana",
      "price": 2.5
    }
  ]
}'::json #> '{order}' -> 0 -> 'fruit';
-- apple

-- 获取订单数组中第 0 个商品的水果价格
-- 这里需要一个转换
-- 因为 #>> 获取的是文本信息
SELECT ('{
  "order": [
    {
      "fruit": "apple",
      "price": 1.5
    },
    {
      "fruit": "banana",
      "price": 2.5
    }
  ]
}'::jsonb #>> '{order}')::json -> 0 -> 'price';
-- 1.5
```

记住：

1. key：整数对应数组，字符串对应对象；
2. 大于：一个大于返回对象，两个大于返回字符串；
3. #：参数是 key 路径；

*表 2* 是 JSONB 数据类型专有操作符。

*表 2 JSONB 的专有操作符*

| 操作符 | 右值类型描述 | 说明                              | 示例                                          |
| ------ | ------------ | --------------------------------- | --------------------------------------------- |
| `@>`   | `jsonb`      | 判断在顶层是否有包含关系          | `'{"a":1, "b":2}'::jsonb @> '{"b":2}'::jsonb` |
| `<@`   | `jsonb`      | 与 `@>` 相反操作                  | `'{"b":2}'::jsonb <@ '{"a":1, "b":2}'::jsonb` |
| `?`    | `text`       | 在顶层是否包含指定字符串的 key    | `'{"a":1, "b":2}'::jsonb ? 'b'`               |
|`?\|`   |`text[]`      |或逻辑，包含其中一个即可 | `'{"a":1, "b":2, "c":3}'::jsonb ?\| array['b', 'c']` |
|`?&`    |`text[]`      |与逻辑，都要包含 | `'["a", "b"]'::jsonb ?& array['a', 'b']` |
|`\|\|`  |`jsonb`       | 连接两个 JSONB 对象 | `'["a", "b"]'::jsonb \|\| '["c", "d"]'::jsonb` |
| `-`    | `text`       | 删除键值对或字符串元素            | `'{"a": "b"}'::jsonb - 'a'`                   |
| `-`    | `integer`    | 删除指定位置的元素                | `'["a", "b"]'::jsonb - 1`                     |
| `#-`   | `text[]`     | 基于路径的删除                    | `'["a", {"b":1}]'::jsonb #- '{1,b}'`          |

示例还包括：

```sql
SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::jsonb @> '{
  "fruit": "apple"
}'::jsonb;
-- true

SELECT '{
  "shop": {
    "name": "happy house",
    "products": [
      "apple",
      "banana",
      "peach"
    ]
  }
}'::jsonb @> '{
  "shop": {}
}'::jsonb;
-- true

SELECT '{
  "shop": {
    "name": "happy house",
    "products": [
      "apple",
      "banana",
      "peach"
    ]
  }
}'::jsonb @> '{
  "shop": 1
}'::jsonb;
-- false

SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::jsonb ? 'fruit';
-- true

SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::jsonb ? 'color';
-- false

SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::jsonb ?| array ['color', 'fruit'];
-- true

SELECT '{
  "fruit": "apple",
  "price": 1.5
}'::jsonb ?& array ['color', 'fruit'];
-- false

SELECT '{
  "fruit": "apple"
}'::jsonb || '{
  "price": 1.5
}'::jsonb;
-- {"fruit": "apple", "price": 1.5}

SELECT '{
  "fruit": "apple",
  "price": 2.5
}'::jsonb || '{
  "price": 1.5
}'::jsonb;
-- {"fruit": "apple", "price": 1.5}

SELECT '{
  "fruit": "apple",
  "price": 2.5,
  "color": {
    "r": "000",
    "b": "111",
    "g": "222"
  }
}'::jsonb || '{
  "color": {
    "r": 333
  }
}'::jsonb;
-- {"color": {"r": 333}, "fruit": "apple", "price": 2.5}
```



## 相关函数

### to_json 和 to_jsonb

### array_to_json

### row_to_json

### json_build_array 和 jsonb_build_array

### json_build_object 和 jsonb_build_object

### json_object 和 jsonb_object



## vs. MongoDB





## 参考

1. [PostgreSQL JSON 相关的函数 - 官方文档](https://www.postgresql.org/docs/9.5/functions-json.html)
2. [Using JSONB in PostgreSQL: How to Effectively Store & Index JSON Data in PostgreSQL](https://scalegrid.io/blog/using-jsonb-in-postgresql-how-to-effectively-store-index-json-data-in-postgresql/)
3. [给 JSON 类型字段创建索引和约束](https://docs.yugabyte.com/preview/api/ysql/datatypes/type_json/create-indexes-check-constraints/)
4. [数据库函数返回 void - stackoverflow](https://stackoverflow.com/questions/14216716/how-to-create-function-that-returns-nothing)