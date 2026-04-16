+++
date = '2026-04-16T11:17:52+08:00'
draft = false
title = '局部变量优化'
categories = ['后端技术', 'Python']
tags = ['Python', '优化技术']
toc = true
+++

局部变量优化，是指将全局变量 / 函数的地址赋值给局部变量，从而减少内存寻址的开销。

我自己的一个观点，大多数应用的瓶颈在于数据库、磁盘 I/O，所以一般不会用这一优化技术。但我今天在项目代码中看到类似的代码，
又问了豆包，豆包说有场景有必要，所以写一个测试验证一下。

局部变量优化的编码结构如下：

```python
def global_func():
    # ...

def loop_func():
    _global_func = global_func
    for _ in range(BIGGER_NUMBER):
        # call _global_func
```

局部变量优化的适用场景就非常明显，适用于存在大量循环，需要反复调用的函数中。

测试代码如下：

```python
def global_func():
    a, b, c = "局部", "变量", "优化"
    return f"{a}{b}{c}"


def loop_func1():
    for _ in range(100000000):
        global_func()


def loop_func2():
    _global_func = global_func
    for _ in range(100000000):
        _global_func()


if __name__ == "__main__":
    import time
    start_time = time.time()
    loop_func1()
    end_time = time.time()
    print(f"loop_func1 耗时: {end_time - start_time} 秒")

    start_time = time.time()
    loop_func2()
    end_time = time.time()
    print(f"loop_func2 耗时: {end_time - start_time} 秒")
```

输出：

```bash
loop_func1 耗时: 10.524455308914185 秒
loop_func2 耗时: 10.882614850997925 秒
```

看样子，反而还慢了，改成 10000000000，输出：

```bash
loop_func1 耗时: 1051.6418409347534 秒
loop_func2 耗时: 1029.1270589828491 秒
```

从上述例子来看，100 亿次的循环调用，省了 20 秒。

我的观点依旧：先按正常思路实现功能，不要太多考虑“局部变量优化”，因为往往其它因素对执行效率的影响更大。

## 附录

豆包说，当使用 PyPy 而不是 CPython 时，可以完全不用考虑局部变量优化，PyPy 的 JIT 会自动做此类优化。

如果查看解释器实现：

```python
import sys
print(sys.executable)       # 解释器路径
print(sys.version)           # 版本信息
print(sys.implementation)    # 核心判断
```
