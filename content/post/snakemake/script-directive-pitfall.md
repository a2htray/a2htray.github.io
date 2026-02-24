+++
date = '2025-12-26T16:11:19+08:00'
draft = false
title = 'Snakemake 中 script 指令无法使用 input 指令中的值'
categories = ['后端技术', 'Snakemake']
tags = ['Snakemake']
toc = true
+++

如题。

## 示例

写个 demo 示例一下。

我的 `Snakefile` 内容如下：

```yaml
rule read_and_write:
    input:
        "inputs/abc.log"
    output:
        "outputs/abcde.log"
    script:
        "scripts/read_file.py -f {input} {output}"
```

运行上面的脚本，会报：

```text
...
The name 'input' is unknown in this context. Please make sure that you defined that variable. Also note that braces not used for variable access have to be escaped by repeating them, i.e. {{print $1}}
...
```

<span style="color: red;">**显示 `input` 不在上下文中**</span>。

## 问题

看了官方的文档：

[https://snakemake.readthedocs.io/en/stable/snakefiles/rules.html#snakefiles-external-scripts](https://snakemake.readthedocs.io/en/stable/snakefiles/rules.html#snakefiles-external-scripts)

才知道，人家就没打算在字符串传递 input 指令的值，而是希望开发在编写的脚本使用以下方式（python）：

```python
from snakemake.script import snakemake

def do_something(data_path, out_path, threads, myparam):
    # python code

do_something(snakemake.input[0], snakemake.output[0], snakemake.threads, snakemake.config["myparam"])
```

## 处理

所以，老老实实地用 `shell` 指令吧。

我个不人推荐在自定义脚本中过多的集成 snakemake 的代码，因为有时，我们希望不依赖 snakemake 环境也能正常执行。

调整后如下：

```yaml
rule read_and_write:
    input:
        "inputs/abc.log"
    output:
        "outputs/abcde.log"
    shell:
        "python scripts/read_file.py -f {input} {output}"
```

