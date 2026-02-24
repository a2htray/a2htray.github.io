+++
date = '2025-12-10T09:19:39+08:00'
draft = false
title = '详解：如何将一个大的结构化（CSV）数据导入到 MySQL'
categories = ['后端技术', 'MySQL']
tags = ['工作', '解决问题', 'Shell']
toc = true
+++

## 实际需求

昨天领了一个任务，把一个大文件数据导入到 MySQL，文件以 CSV 格式给出，内容是收集的 600 多万条的文献信息。

接到时，我觉得很简单，实现思路是：

1. （前置）部署一个 MySQL 的服务
2. 根据 CSV 的表头信息，创建一个表
3. 将大文件拆成多个小文件，小文件包含表头信息
4. 遍历所有小文件，应用 `LOAD DATA INFILE` 语句进行导入

为了让整个过程更好追踪，还加入`小文件导入日志`和`汇总主日志`的功能。因为是一个即时的任务，并不需要考虑太多后期的扩展，所以采用 `Shell` 脚本
的方式实现。

但在实现的过程中，还是遇到了一些不少的问题，大概也花了 3 个小时左右。

## MySQL 环境

MySQL，我用 Docker 启的，因为考虑到后面要做数据查询服务，所以按 `docker-compose` 的方式进行编排，内容如下：

```yaml
version: '3.8'

services:
  demo_db:
    image: mysql:8.0
    container_name: demo_db
    ports:
      - 13306:3306
    volumes:
      - demo_mysql_data:/var/lib/mysql
      - ../data:/var/lib/data
      - ../scripts:/var/lib/scripts
    networks:
      - demo_network
    environment:
      - MYSQL_ROOT_PASSWORD=demo_20251209
      - MYSQL_DATABASE=demo
      - TZ=Asia/Shanghai

volumes:
  demo_mysql_data:

networks:
  demo_network:
    driver: bridge
```

> docker-compose 目前还不包含应用的部分。

## 创建 MySQL 表

CSV 文件的表头很多，将近 90 个，所以不太可能逐一映射到数据库的字段。在这种情况下，大语言模型就派上用场了。

我比较常用的工具是豆包。

首先，我把文件前两行的内容给到豆包，它帮我说明了表头信息及对应的内容解释。

然后，又给豆包提了个要求，生成 Create 语句，内容如下：

```text
请根据我给的信息，帮我生成一个名称叫 demo_table 的 mysql 表的 CREATE 语句，
要求字段要小写，且顺序不能变。
```

果然，豆包给了一个详尽可用的 `CREATE` 语句

最后，复制语句为 `./scripts/create_demo_table.sql` 文件，并在容器中执行：

```bash
$ mysql -rroot -pdemo_20251209 demo < ./scripts/create_demo_table.sql
```

**到这里，就会为后面埋下一个坑。**

至此，MySQL 环境已经具备了。

## 编写 Shell 脚本

Shell 脚本要主要实现两个功能：

1. 拆分大文件
2. 按小文件导入并记录日志

### 拆分大文件

**拆分大文件**后小文件都要带相同的表头。

首先，使用 `head` 命令将表头输入到临时文件。

```bash
head -1 ./BIG_FILE.csv > ./header.tmp
```

* `head -1` 表示只取头一行。

然后，从第 2 行（包括第 2 行）开始将大文件拆成小文件，我按 5000 条记录一个文件的方式拆，命令如下：

```bash
mkdir ./split_dir
tail -n +2 ./BIG_FILE.csv > split -l 5000 - ./split_dir/samll_file
```

* `tail -n +2 ./BIG_FILE.csv` 从第 2 行（包括第 2 行）将文件内容输出终端
* `-l 5000` 以每个文件 5000 条记录来拆
* `split` 中的 `-` 表示从标准输出中读取数据
* `./split_dir/samll_file` 小文件放置的目录及其前缀

最后，给每个小文件加入表头，这里需要用到 `for` 循环：

```bash
for samll_file in ./split_dir/small_file*; do
  cat ./header.tmp ${small_file} > ${small_file}.tmp
  mv ${small_file}.tmp ${small_file}
done
```

* `for` 循环，遍历目录下的小文件
* `cat` 接两个文件，将两个文件的内容输出到临时文件
* `mv` 将临时文件命名为原名

至此，大文件拆成小文件已经完成。

### 按小文件导入并记录日志

这里，我需要一个主日志文件来记录每个小文件的导入情况：

```bash
main_log_file=./main.log
```

可以先记录一个拆分的统计信息，如拆分成小文件的数量：

```bash
small_files=(./split_dir/small_file*)
small_file_num=${#small_files[@]}

echo "大文件拆小文件完成，共拆分成 ${small_file_num} 个小文件" >> ${main_log_file}
```

* `${#数组名[@]}` 是 Shell 中获取数组元素个数的标准语法

再建一个 `for` 循环进行数据的导入，同时写入各自小文件的日志文件：

```bash
for samll_file in ${samll_files}; do
  echo "开始导入 ${samll_file}" >> ${main_log_file}
  
  file_name=(basename ${small_file})
  small_file_log="./{$file_name}.log"
  
  mv ${small_file} /var/lib/mysql-files/${file_name} 2>&! >> ${small_file_log}
  
  echo "执行 LOAD DATA INFILE 导入" >> ${small_file_log}
  sql="
LOAD DATA INFILE '/var/lib/mysql-files/${file_name}'
INTO TABLE demo_table
CHARACTER SET utf8mb4
FIELDS TERMINATED BY ','
OPTIONALLY ENCLOSED BY '\"'
LINES TERMINATED BY '\n'
IGNORE 1 ROWS"
(${table_fields});
"

  if mysql -rroot -pdemo_20251209 demo -e "${sql}" >> ${samll_file_log} 2>&1; then
    echo "SUCCESS: ${samll_file_log} 导入成功！" >> ${main_log_file}
  else
    echo "ERROR: 导入失败！详情见 ${samll_file_log}" >> ${main_log_file}
  fi
done  
```

> 注意：`${table_fields}` 里面是按 CSV 中表头顺序映射到数据库的字段名，使用逗号分隔。

脚本的逻辑还是非常清晰的：

1. 遍历每一个小文件，声明小文件的日志文件
2. 将小文件移到到 `/var/lib/mysql-files` 目录下
3. 执行 `LOAD DATA INFIILE` 语句
4. 将命令的终端输出写入到小文件日志文件
5. 根据命令的执行成功与否，将信息写入到主日志文件

主要功能到目前，都已经完成了，后面都是一些锦上添花的事。

## 小坑

还记得上面说的“留了一个坑”。究其原因，在于豆包根据给定信息生成的 CREATE 语句中，字段类型无法涵盖真实字段值。

比如，CREATE 语句中的字段类型是：

```sql
VARCHAR(64)
```

但是，真实值超过 64 个字符。再比较，有一些日期字段：

```sql
DATE
```

而真实值可能如下：

```text
2020
2020.12
```

类似这种情况，只能导一次，看到哪个字段类型声明不恰当，修改一遍再执行。

## 小结

任务其实是很简单的，整体思路也非常清晰，优化的空间也很有，比如将小文件的导入由串行改成并行。

