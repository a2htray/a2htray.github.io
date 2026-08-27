+++
date = '2026-08-27T17:24:00+08:00'
draft = false
title = 'PostgreSQL 时间序列扩展 - TimescaleDB'
categories = ['后端技术', 'PostgreSQL']
tags = ['PostgreSQL', 'TimescaleDB', '数据库']
toc = true
+++

## TimescaleDB 介绍

TimescaleDB 是一个基于 **PostgreSQL 的时间序列数据库扩展**（不是独立数据库），通过安装扩展把 Postgres 变成能高效处理时序数据的引擎，兼容所有 Postgres 生态。TimescaleDB 使用场景下的数据带有 3 个明显的特征：

* 带时间戳：数据中包含时间戳字段，用于记录数据的采集时间
* 写入数据量大：持续高频追加
* 按时间查询：特定时间区间内计算聚合、趋势、同比等指标

## 六大核心概念

**Hypertable（超表）**

对用户呈现的 "一张表"，内部按时间自动拆分成多张物理子表。

* 超表本身不存数据，只存定义和 chunk 元信息
* 每次 INSERT 按`策略`自动路由到对应 chunk
* 创建超表后，应用层依然像操作普通 Postgres 表一样读写，无需关心底层分片
* 支持**二级分区**，将单个时间块再打散到多个 chunk，避免单个 chunk 过大、提升写入并发

**Chunk（数据块）**

Chunk 是超表的基本单位，每个 chunk 包含一定时间范围内的数据。

* **一个 chunk 对应一个物理子表**，名字形如 `_hyper_X_Y`
* 查询命中时间条件时，TimescaleDB 会自动做**分区裁剪（partition pruning）**，只扫描相关 chunk，跳过无关时间范围，这是时序查询快的核心原因
* chunk 的数量、大小可通过 `chunk_time_interval` 调优：写入越密集，应设越小的区间以避免单 chunk 过大；查询跨度越大，适当放大区间减少 chunk 数量
* 可通过 `timescaledb_information.chunks` 视图查看每个 chunk 的时间范围与体积

**Compression（压缩）**

TimescaleDB 的压缩不是简单的列存 ZIP，而是**列式压缩 + 时序专用编码**结合：

* 默认对 chunk 启用压缩后，数据从行存转为列存，并对数值列使用 delta-of-delta、GORILLA、字典编码等算法
* 典型时序场景（高频、列相关性高）可获得 **90%+ 的空间压缩比**
* 压缩是**按 chunk 异步进行的**，默认只压缩"较旧、不再频繁写入"的 chunk，避免冻结正在写入的数据
* 压缩后的 chunk 仍然**可查询**，查询时透明解压，对应用无感

**Continuous Aggregates（连续聚合）**

连续聚合相当于"物化视图 + 自动增量刷新"，用于预计算高频访问的聚合结果：

* 定义一次聚合 SQL（如按天/按股票汇总），TimescaleDB 后台定时把**新增数据**增量合并进结果，而不是每次全量重算
* 业务查询直接读聚合表，避免对原始大表反复 `GROUP BY`，查询延迟从秒级降到毫秒级
* 通过 `WITH NO DATA` / `REFRESH` 控制首次填充，通过 `refresh_lag` / `refresh_interval` 控制刷新滞后与频率
* 非常适合 K 线、同比环比、监控指标等"固定聚合 + 实时查询"场景

**Retention Policy（数据保留策略）**

保留策略用于自动清理过期数据，控制存储成本与查询性能：

* 通过 `add_retention_policy` 为超表设定"保留窗口"（如只保留最近 2 年）
* 到期 chunk 会被后台 job **自动删除**，无需手写 `DELETE` 脚本
* 可配合连续聚合：原始明细只保留短期，聚合结果长期保留，既省空间又保历史统计

**Data Tiering（分层存储）**

分层存储把"冷热数据"放到不同成本的存储介质上：

* 热数据（近期、频繁访问）留在本地高速 SSD
* 冷数据（历史、很少访问）自动下沉到对象存储（如 S3 / 兼容 S3 的存储），按访问透明拉取
* 在云版本（Timescale Cloud / TimescaleDB 2.13+ 的 `tiered storage`）中开箱即用，自托管版本需结合外部对象存储配置
* 对上层查询无感：查询历史区间时，引擎自动从对应层级读取，用户只看到一张统一的超表

## TimescaleDB Docker 安装

创建 `docker-compose.yml` 文件，内容如下：

```yml
services:
  timescaledb:
    image: timescale/timescaledb:latest-pg16
    container_name: timescaledb
    restart: unless-stopped
    ports:
      - "5432:5432"
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: yourpassword
      POSTGRES_DB: postgres
    volumes:
      - timescale-data:/var/lib/postgresql/data
    # 健康检查：确认扩展已启用
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 10s
      timeout: 5s
      retries: 5

volumes:
  timescale-data:
```

检查运行状态：

```bash
$ docker compose ps
NAME          IMAGE                               COMMAND                  SERVICE       CREATED         STATUS                   PORTS
timescaledb   timescale/timescaledb:latest-pg16   "docker-entrypoint.s…"   timescaledb   3 minutes ago   Up 3 minutes (healthy)   0.0.0.0:5432->5432/tcp, [::]:5432->5432/tcp
```

检查是否启用了扩展：

```bash
$ docker compose exec -it timescaledb psql -U postgres
psql (16.15)
Type "help" for help.

postgres=# SELECT extname, extversion FROM pg_extension WHERE extname = 'timescaledb';
   extname   | extversion 
-------------+------------
 timescaledb | 2.29.2
(1 row)
```

## TimescaleDB 使用

### 建表、转为超表

```sql
CREATE EXTENSION IF NOT EXISTS timescaledb;

-- 建表
CREATE TABLE IF NOT EXISTS stock_daily
(
    symbol     TEXT NOT NULL,
    trade_date DATE NOT NULL,
    open       NUMERIC(12, 2),
    high       NUMERIC(12, 2),
    low        NUMERIC(12, 2),
    close      NUMERIC(12, 2),
    volume     BIGINT,
    amount     NUMERIC(18, 2),
    pct_chg    NUMERIC(10, 4),
    turnover   NUMERIC(12, 4),
    PRIMARY KEY (symbol, trade_date)
);

COMMENT ON TABLE stock_daily IS '股票日线行情表';
COMMENT ON COLUMN stock_daily.symbol IS '股票代码, 如 600519.SH / 000001.SZ';
COMMENT ON COLUMN stock_daily.trade_date IS '交易日';
COMMENT ON COLUMN stock_daily.open IS '开盘价(元)';
COMMENT ON COLUMN stock_daily.high IS '最高价(元)';
COMMENT ON COLUMN stock_daily.low IS '最低价(元)';
COMMENT ON COLUMN stock_daily.close IS '收盘价(元)';
COMMENT ON COLUMN stock_daily.volume IS '成交量(股)';
COMMENT ON COLUMN stock_daily.amount IS '成交额(元)';
COMMENT ON COLUMN stock_daily.pct_chg IS '涨跌幅(%)';
COMMENT ON COLUMN stock_daily.turnover IS '换手率(%)';

-- 转为超表
SELECT create_hypertable(
               'stock_daily', -- 要转为超表的表名
               'trade_date', -- 时间分区列，按 trade_date 字段作时间维度分片
               chunk_time_interval => INTERVAL '1 month', -- 时间块跨度
               partitioning_column => 'symbol', -- 二级分区列
               number_partitions => 8, -- 哈希分片数
               if_not_exists => TRUE -- 已存在则跳过
       );
```

### 查询库里已有的超表

```sql
SELECT * FROM timescaledb_information.hypertables;
```

### 插入数据

使用 `akshare` 包获取股票数据（中国中车、招商银行），插入到 `stock_daily` 表中。Python 完整代码如下：

```python
import datetime as dt
from decimal import Decimal

import akshare as ak
import psycopg2
from psycopg2.extras import execute_values


STOCKS = [
    {"name": "中国中车", "code": "601766", "symbol": "601766.SH"},
    {"name": "招商银行", "code": "600036", "symbol": "600036.SH"},
]

END_DATE = dt.date.today()
START_DATE = END_DATE - dt.timedelta(days=365)

INSERT_SQL = """
INSERT INTO stock_daily
    (symbol, trade_date, open, high, low, close, volume, amount, pct_chg, turnover)
VALUES %s
ON CONFLICT (symbol, trade_date) DO NOTHING
"""


def _d(value) -> Decimal | None:
    if value is None:
        return None
    return Decimal(str(value))


def fetch_stock_daily(code: str, start: dt.date, end: dt.date) -> list[dict]:
    df = ak.stock_zh_a_hist(
        symbol=code,
        period="daily",
        start_date=start.strftime("%Y%m%d"),
        end_date=end.strftime("%Y%m%d"),
        adjust="",
    )

    records = []
    for _, row in df.iterrows():
        records.append(
            {
                "trade_date": dt.datetime.strptime(str(row["日期"]), "%Y-%m-%d").date(),
                "open": _d(row["开盘"]),
                "high": _d(row["最高"]),
                "low": _d(row["最低"]),
                "close": _d(row["收盘"]),
                "volume": int(row["成交量"]) if row["成交量"] is not None else None,
                "amount": _d(row["成交额"]),
                "pct_chg": _d(row["涨跌幅"]),
                "turnover": _d(row["换手率"]),
            }
        )
    return records


def main():
    conn = psycopg2.connect(
        host="localhost",
        port=5432,
        user="postgres",
        password="yourpassword",
        dbname="postgres",
    )
    conn.autocommit = False

    try:
        for stock in STOCKS:
            print(f"正在获取 {stock['name']}({stock['code']}) 行情数据 ...")
            records = fetch_stock_daily(stock["code"], START_DATE, END_DATE)

            rows = [
                (
                    stock["symbol"],
                    rec["trade_date"],
                    rec["open"],
                    rec["high"],
                    rec["low"],
                    rec["close"],
                    rec["volume"],
                    rec["amount"],
                    rec["pct_chg"],
                    rec["turnover"],
                )
                for rec in records
            ]

            with conn.cursor() as cur:
                execute_values(cur, INSERT_SQL, rows, page_size=500)
            conn.commit()
            print(f"  -> 插入 {len(rows)} 条记录")

    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()


if __name__ == "__main__":
    main()
```

执行如下：

```bash
$ python insert_stock_daily.py 
正在获取 中国中车(601766) 行情数据 ...
  -> 插入 243 条记录
正在获取 招商银行(600036) 行情数据 ...
  -> 插入 243 条记录
```

### 查询 Chunk 信息

```sql
SELECT hypertable_name, chunk_name, range_start, range_end
FROM timescaledb_information.chunks
WHERE hypertable_name = 'stock_daily';
```

![](/imgs/learn-db/ScreenShot_2026-08-27_151403_242.png)

> 满足定义的按月分区。

### 查询数据

按月聚合数据，查询每个股票的月最高、月最低、月成交量、月成交额、平均涨跌幅。

```sql
SELECT symbol,
       time_bucket('1 month', trade_date) AS month,
       MAX(high)                          AS 月最高,
       MIN(low)                           AS 月最低,
       SUM(volume)                        AS 月成交量,
       SUM(amount)                        AS 月成交额,
       AVG(pct_chg)                       AS 平均涨跌幅
FROM stock_daily
GROUP BY symbol, month
ORDER BY symbol, month DESC;
```

查询中国中车（601766.SH）最近一年的 5 天数据，按涨跌幅降序排序。

```sql
SELECT symbol, trade_date, close, pct_chg
FROM stock_daily
WHERE trade_date >= CURRENT_DATE - INTERVAL '1 year'
  AND symbol = '601766.SH'
ORDER BY pct_chg DESC
LIMIT 5;
```

查询招商银行（600036.SH）近 5 日、10 日、20 日均线。

```sql
SELECT symbol,
       trade_date,
       close,
       ROUND(AVG(close) OVER (
           PARTITION BY symbol ORDER BY trade_date
           ROWS BETWEEN 4 PRECEDING AND CURRENT ROW
           ), 2) AS ma5,
       ROUND(AVG(close) OVER (
           PARTITION BY symbol ORDER BY trade_date
           ROWS BETWEEN 9 PRECEDING AND CURRENT ROW
           ), 2) AS ma10,
       ROUND(AVG(close) OVER (
           PARTITION BY symbol ORDER BY trade_date
           ROWS BETWEEN 19 PRECEDING AND CURRENT ROW
           ), 2) AS ma20
FROM stock_daily
WHERE symbol = '600036.SH'
ORDER BY symbol, trade_date DESC;
```

## 性能与使用建议

* **选好时间列**：`create_hypertable` 的时间列应是最常用的查询过滤维度；本文用 `trade_date` 而非 `timestamp`，契合日线场景
* **chunk 不要太大也不要太小**：单 chunk 建议控制在 25MB ~ 数 GB；写入密集可缩小 `chunk_time_interval`，查询跨度大可适当放大
* **写入用批量**：如文中 `execute_values(page_size=500)`，比单条 `INSERT` 快一个数量级
* **先压缩再保留**：压缩降成本，保留控体积，二者组合是时序存储的标准套路
* **连续聚合下沉热点查询**：把"固定聚合"预计算掉，明细表只做明细分析
* **善用分区裁剪**：查询务必带上时间条件，否则会扫描全部 chunk 失去优势

## 小结

TimescaleDB 不是另起炉灶的数据库，而是把 PostgreSQL 变成一个"为时序而生"的引擎：超表 + chunk 解决分区与扩展，压缩 + 连续聚合解决成本与查询，保留 + 分层解决生命周期。对已经使用 Postgres 的团队而言，无需引入新技术栈，就能平滑获得时序数据库的能力。