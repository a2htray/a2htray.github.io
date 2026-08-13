+++
date = '2026-08-13T10:27:52+08:00'
draft = false
title = 'Neo4j 搭建（Docker）与数据备份还原实战'
categories = ['后端技术', '数据库']
tags = ['Neo4j', 'Cypher', '图数据库', '备份还原']
toc = true
+++

> 适用版本：Neo4j 5.x（社区版 / 企业版）
> 环境：已安装 Docker 与 Docker Compose 的 Linux / macOS 主机

## 前言

Neo4j 是目前最流行的原生图数据库，常用于知识图谱、关系网络、推荐系统等场景。使用 Docker 部署可以让环境保持一致、迁移方便，而图数据的备份与还原则是生产环境中不可忽视的一环。

本文覆盖三部分：

1. 用 Docker / Docker Compose 快速搭建 Neo4j
2. 三种数据备份方式：官方 `dump`、企业版在线 `backup`、APOC 导出
3. 对应的数据还原方法，以及定时自动备份脚本

## 环境准备

确认 Docker 与 Docker Compose 已安装：

```bash
$ docker --version
Docker version 29.4.1, build 055a478
$ docker compose version
Docker Compose version v5.1.3
```

创建本地目录用于持久化数据与备份（**强烈建议挂载宿主机目录，否则容器删除后数据丢失**）：

```bash
$ mkdir {data,logs,plugins,import,backups}
# 设置目录权限，确保容器可以读写
$ chmod -R 777 .
```

## 使用 Docker 搭建 Neo4j

### 方式一：docker run 直接启动

```bash
$ docker run -d \
  --name neo4j \
  -p 7474:7474 \
  -p 7687:7687 \
  -e NEO4J_AUTH=neo4j/password \
  -v ./data:/data \
  -v ./logs:/logs \
  -v ./plugins:/plugins \
  -v ./import:/var/lib/neo4j/import \
  -v ./backups:/backups \
  neo4j:5.26
```

* `-p 7474:7474` HTTP 浏览器界面
* `-p 7687:7687` Bolt 协议，用于程序连接

启动后访问 `http://<主机IP>:7474`，用 `neo4j / password` 登录即可，如下图：

![](/imgs/neo4j/01.png)

### 方式二：Docker Compose（推荐）

创建 `docker-compose.yml`：

```yaml
version: '3.8'

services:
  neo4j:
    image: neo4j:5.26
    container_name: neo4j
    ports:
      - "7474:7474"
      - "7687:7687"
    environment:
      # 初始认证：账号/密码
      NEO4J_AUTH: neo4j/password
      # 自动安装 APOC 插件（版本需与 Neo4j 匹配）
      NEO4J_PLUGINS: '["apoc"]'
      # APOC 文件导入导出权限（生产环境按需收紧）
      NEO4J_dbms_security_procedures_unrestricted: apoc.*
      NEO4J_apoc_export_file_enabled: "true"
      NEO4J_apoc_import_file_enabled: "true"
      NEO4J_apoc_import_file_use__neo4j__config: "true"
    volumes:
      - ./data:/data
      - ./logs:/logs
      - ./plugins:/plugins
      - ./import:/var/lib/neo4j/import
      - ./backups:/backups
    restart: unless-stopped
    healthcheck:
      test: ["CMD-SHELL", "cypher-shell -u neo4j -p password 'RETURN 1'"]
      interval: 15s
      timeout: 10s
      retries: 5
```

启动：

```bash
$ docker compose up -d
$ docker compose ps
$ docker compose logs -f neo4j
```

停止与删除（**注意：不加 `-v` 不会删除数据卷**）：

```bash
docker compose stop          # 仅停止，保留数据
docker compose down          # 停止并删除容器，保留卷
docker compose down -v       # ⚠️ 停止并删除容器及所有数据卷（危险）
```

### 安装 APOC 插件（两种方式）

**方式 A：通过环境变量（最简单，见「方式二：Docker Compose」中的 `NEO4J_PLUGINS`）**

容器首次启动时会自动下载匹配版本的 APOC 并安装，重启生效。

**方式 B：手动放 jar 包**

从 https://github.com/neo4j/apoc/releases 下载与 Neo4j 同版本的 `apoc-<version>-core.jar`，放入 `~/neo4j/plugins`，然后：

```bash
docker compose restart neo4j
```

验证插件是否生效（在 Neo4j Browser 或 cypher-shell 执行）：

```bash
CALL apoc.help("path");
```

![](/imgs/neo4j/02.png)

有返回结果即安装成功。

### 准备测试数据（用于备份/还原演练）

为了验证备份与还原是否完整、一致，先灌入一份**可复现**的测试数据。下面创建一个「人物—电影」小图（8 个节点、6 条关系），后续备份前后都用节点/关系数量做对比校验。

在 Neo4j Browser 或 cypher-shell 中执行：

```bash
// 清空当前库（仅测试环境使用，生产环境请勿执行）
MATCH (n) DETACH DELETE n;

// 创建 4 个人物节点 + 4 个电影节点
CREATE
  (tom:Person {name:'Tom Hanks', born:1956}),
  (keanu:Person {name:'Keanu Reeves', born:1964}),
  (mega:Person {name:'Meg Ryan', born:1961}),
  (rob:Person {name:'Rob Reiner', born:1947}),
  (forrest:Movie {title:'Forrest Gump', released:1994, tagline:'Life is like a box of chocolates.'}),
  (sleepless:Movie {title:'Sleepless in Seattle', released:1993, tagline:'What if someone you never met was the one for you?'}),
  (matrix:Movie {title:'The Matrix', released:1999, tagline:'Welcome to the Real World.'}),
  (matrix2:Movie {title:'The Matrix Reloaded', released:2003, tagline:'Free your mind.'});

// 创建关系：演员出演、导演执导
MATCH (tom:Person {name:'Tom Hanks'}), (forrest:Movie {title:'Forrest Gump'})
CREATE (tom)-[:ACTED_IN {roles:['Forrest']}]->(forrest);

MATCH (tom:Person {name:'Tom Hanks'}), (sleepless:Movie {title:'Sleepless in Seattle'})
CREATE (tom)-[:ACTED_IN {roles:['Sam Baldwin']}]->(sleepless);

MATCH (mega:Person {name:'Meg Ryan'}), (sleepless:Movie {title:'Sleepless in Seattle'})
CREATE (mega)-[:ACTED_IN {roles:['Annie Reed']}]->(sleepless);

MATCH (rob:Person {name:'Rob Reiner'}), (forrest:Movie {title:'Forrest Gump'})
CREATE (rob)-[:DIRECTED]->(forrest);

MATCH (keanu:Person {name:'Keanu Reeves'}), (matrix:Movie {title:'The Matrix'})
CREATE (keanu)-[:ACTED_IN {roles:['Neo']}]->(matrix);

MATCH (keanu:Person {name:'Keanu Reeves'}), (matrix2:Movie {title:'The Matrix Reloaded'})
CREATE (keanu)-[:ACTED_IN {roles:['Neo']}]->(matrix2);
```

可视化图如下：

![](/imgs/neo4j/03.png)

写入后立刻校验基线数据量（这个值在备份前、还原后都应一致）：

```bash
MATCH (n) RETURN count(n) AS 节点数;
MATCH ()-[r]->() RETURN count(r) AS 关系数;
```

预期结果：**节点数 = 8，关系数 = 6**。

> 小技巧：把上面的 Cypher 保存为 `./import/test-data.cypher`，之后做还原演练时，可先 `MATCH (n) DETACH DELETE n` 清空，再 `:source import/test-data.cypher` 一键重新灌入，方便反复练习备份/还原。

## 基础配置与验证

修改配置可通过环境变量或在 `~/neo4j/data/../conf` 中编辑（建议用命名卷挂载 `conf` 目录以便持久化）。常用项：

| 配置项 | 说明 |
| --- | --- |
| `NEO4J_AUTH` | 账号密码认证 |
| `NEO4J_dbms_memory_heap_initial__size` | 初始堆内存 |
| `NEO4J_dbms_memory_heap_max__size` | 最大堆内存 |
| `NEO4J_dbms_default__database` | 默认数据库名（默认 `neo4j`） |

验证连接：

```bash
$ docker exec -it neo4j cypher-shell -u neo4j -p password "MATCH (n) RETURN count(n) AS nodeCount;" 
+-----------+
| nodeCount |
+-----------+
| 8         |
+-----------+

1 row
ready to start consuming query after 718 ms, results consumed after another 2 ms
```

查看版本与组件：

```bash
CALL dbms.components();
```

![](/imgs/neo4j/04.png)

## 数据备份

Neo4j 提供三类备份手段，按场景选择：

| 方式 | 适用版本 | 是否需停库 | 产物 |
| --- | --- | --- | --- |
| `neo4j-admin database dump` | 社区版 / 企业版 | **需停库**（离线） | 单个 `.dump` 文件 |
| `neo4j-admin database backup` | 仅企业版 | 无需停库（在线） | 备份目录（可增量） |
| APOC 导出 | 社区版 / 企业版（需装 APOC） | 在线 | Cypher / CSV / JSON 等 |

### 社区版离线备份（dump）—— 最常用

`dump` 要求 **数据库离线**——不能有任何进程持有数据目录的锁文件 `/data/databases/<db>/database_lock`。

> ⚠️ **两点关键坑**
> 1. 不要用「临时容器挂载同一个 `./data`」去对正在运行的 `neo4j` 容器做 dump。两个容器共享同一份数据目录，运行中的容器已经持有了锁，临时容器再来抢同一把锁会报错：`The database is in use` / `Lock file has been locked by another process`。临时容器的意义只是「不必在业务容器里 exec」，**并不能绕过离线要求**。
> 2. `STOP DATABASE` / `START DATABASE` 是**企业版（Enterprise）专属**的数据库管理命令，社区版执行会直接报 `Unsupported administration command: STOP DATABASE neo4j`。**社区版无法只停某一个库**，只能停掉整个服务进程来释放锁。

#### 社区版（Community）正确做法：停掉整个容器再做 dump

```bash
$ docker stop neo4j
neo4j
$ docker run --rm -it -v ./data:/data -v ./backups:/backups neo4j:5.26 \
    neo4j-admin database dump neo4j --to-path=/backups --overwrite-destination=true --verbose
$ docker start neo4j
neo4j
```

执行后会在 `./backups/neo4j.dump` 生成备份文件。

> 为什么这样可行：先 `docker stop neo4j` 让运行中的容器释放数据目录锁，再用一个一次性临时容器挂载同一份 `./data` 执行 dump（此时已无锁冲突），最后 `docker start neo4j` 拉起服务。

### APOC 在线导出（可选）

适合需要把数据导出为可读格式（便于迁移、审计、跨版本）的场景。在 cypher-shell 或 Browser 中执行：

```bash
// 导出为 Cypher 脚本（可直接回灌）
CALL apoc.export.cypher.all('/var/lib/neo4j/import/export.cypher', {});
```

![](/imgs/neo4j/05.png)

```bash
CALL apoc.export.csv.all('/var/lib/neo4j/import/export.zip', {useOptimizations: {type: 'UNWIND_BATCH', batchSize: 20000}});
```

![](/imgs/neo4j/07.png)

```bash
// 导出为 JSON
CALL apoc.export.json.all('/var/lib/neo4j/import/export.json', {});
```

![](/imgs/neo4j/06.png)

> 文件会写入容器内的 `/var/lib/neo4j/import`，因已挂载到宿主机 `~/neo4j/import`，可直接在宿主机拿到。

## 数据还原

还原是备份的逆操作，同样遵循「离线用 dump/load，在线用 backup/restore」。

### 社区版离线还原（load）

`load` 会**覆盖目标库现有数据**，请先确认。与 `dump` 同理，`load` 也要求数据库离线、不能持有锁。注意两点：

- `STOP DATABASE` / `START DATABASE` 是**企业版专属**命令，社区版执行会报 `Unsupported administration command`，因此社区版必须停掉整个容器来释放锁。
- 社区版只有默认的 `neo4j` 一个库，`load` 只能还原到 `neo4j`；还原到**新库名**属于多库能力，仅企业版支持。

#### 社区版（Community）正确做法：停掉整个容器再 load

```bash
$ docker stop neo4j
$ docker run --rm -it -v ./data:/data -v ./backups:/backups neo4j:5.26 \
    neo4j-admin database load neo4j --from-path=/backups --overwrite-destination=true
$ docker start neo4j
```

### APOC 导入

对应上文「APOC 在线导出」的导出格式：

```bash
// 从 Cypher 脚本导入
CALL apoc.cypher.runFile('/var/lib/neo4j/import/export.cypher', {});

// 从 CSV 导入
CALL apoc.import.csv(
  [{fileName:'file:///import/nodes.csv', labels:['Person']}],
  [{fileName:'file:///import/rels.csv'}],
  {}
);

// 从 JSON 导入
CALL apoc.import.json('/var/lib/neo4j/import/export.json', {});
```

> 还原后校验数据完整性：执行 `MATCH (n) RETURN count(n)` 与 `MATCH ()-[r]->() RETURN count(r)`，应与「准备测试数据」节基线一致（节点数 = 8，关系数 = 6）。若数量不符，说明备份或还原过程有问题。

## 定时自动备份

把上文 dump 命令封装成脚本，配合 `cron` 每日执行。

创建脚本 `~/neo4j/backup.sh`：

```bash
#!/usr/bin/env bash
set -euo pipefail

NEO4J_CONTAINER="neo4j"
NEO4J_IMAGE="neo4j:5.26"
DATA_DIR="$HOME/neo4j/data"
BACKUP_DIR="$HOME/neo4j/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p "$BACKUP_DIR"

# 社区版必须停掉整个服务进程才能释放数据目录锁
# （STOP DATABASE 是企业版专属命令，社区版会报 Unsupported administration command）
docker stop "$NEO4J_CONTAINER"

# 无论 dump 成功与否，退出时都拉起容器，避免库长时间停服
trap 'docker start "$NEO4J_CONTAINER" || true' EXIT

docker run --rm -i \
  -v "$DATA_DIR":/data \
  -v "$BACKUP_DIR":/backups \
  "$NEO4J_IMAGE" \
  neo4j-admin database dump neo4j --to-path=/backups --overwrite-destination=true

# 重命名带时间戳，并保留最近 7 份
mv "$BACKUP_DIR/neo4j.dump" "$BACKUP_DIR/neo4j_$DATE.dump"
ls -t "$BACKUP_DIR"/neo4j_*.dump | tail -n +8 | xargs -r rm -f

echo "[$(date)] backup done: neo4j_$DATE.dump"
```

赋权并加入 cron（每日凌晨 3 点）：

```bash
chmod +x ~/neo4j/backup.sh

# 编辑当前用户定时任务
crontab -e
# 添加一行：
# 0 3 * * * /bin/bash $HOME/neo4j/backup.sh >> $HOME/neo4j/backups/backup.log 2>&1
```

## 版本升级与跨环境迁移

### 小版本升级（5.x → 5.y）

1. 全量 `dump` 备份
2. 修改 `docker-compose.yml` 中镜像 tag 后 `docker compose up -d`
3. 首次启动自动执行存储迁移
4. 校验：`CALL dbms.components();`

### 大版本跨越（4.x → 5.x）

必须经官方路径（通常 4.4 LTS → 5.x）：

1. 在 4.x 实例 `dump`
2. 新装 5.x 空实例后 `load`
3. 检查废弃配置项（5.x 的 `neo4j.conf` 改动较大）
4. 验证 APOC / GDS 插件版本匹配

## 常见问题

**Q1：容器启动后无法写入数据 / 权限报错**

A：宿主机目录未授权给 UID 7474。执行 `sudo chown -R 7474:7474 ~/neo4j`。

**Q2：`neo4j-admin database load` 提示数据库已存在**

A：加 `--overwrite-destination=true`，或在旧版本（如 4.x）用 `--force`。

**Q3：`dump` / `load` 是否必须停库？**

A：社区版 `dump`/`load` 需要数据库离线，故推荐用临时容器挂载数据卷执行；企业版 `backup`/`restore` 支持在线。

**Q4：APOC 导出文件找不到**

A：确认已设置 `NEO4J_apoc_export_file_enabled=true`，且路径落在 `/var/lib/neo4j/import`（已挂载到宿主机）。

**Q5：还原后浏览器显示旧数据**

A：确认载入的是正确的数据库名（`NEO4J_dbms_default__database`），或用 `:use <db>` 切换。

**Q6：`dump` / `load` 报 `The database is in use` 或 `Lock file has been locked by another process`**
A：说明数据目录的锁被别的进程（通常是正在运行的 `neo4j` 容器，挂载了同一份 `./data`）占着。社区版 `dump`/`load` 必须数据库离线：先 `docker stop neo4j` 释放锁，再用一次性临时容器执行 dump/load，最后 `docker start neo4j`。**切勿用一个挂载同一份 `./data` 的临时容器去 dump/load 正在运行的库**——两个容器抢同一把锁必然失败。另外，`STOP DATABASE` / `START DATABASE` 是**企业版专属命令**，社区版执行会报 `Unsupported administration command`，所以社区版只能停掉整个容器、而不能只停某个库。企业版则可 `STOP DATABASE neo4j;` → 操作 → `START DATABASE neo4j;`。