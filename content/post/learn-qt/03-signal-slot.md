+++
date = '2026-11-02T09:40:45+08:00'
draft = false
title = 'PySide6 信号与槽（Signal & Slot）：Qt 最核心机制，写法详解'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'Qt Signal', 'Qt Slot']
toc = true
+++

![](/imgs/learn-qt/cover3.png)

## 信号与槽

- **信号(Signal)**：对象在状态发生变化时主动发出的通知，本身不包含业务逻辑。
    - 由 Qt 内置控件预先定义（`clicked` / `textChanged` / `valueChanged`），或者自定义信号
    - 信号可以携带参数，用于把数据传递给接收方
- **槽(Slot)**：普通可调用函数（方法），用来响应信号；**信号触发时，槽函数被执行**。
    - 可以是类成员方法、普通函数、lambda
    - 槽的参数数量 ≤ 信号参数数量；多余参数可以设置默认值

> 核心优势：**解耦**。发送方只负责发信号，不需要知道谁接收；接收方只负责处理，不需要持有发送方强依赖。
> 底层：属于**发布订阅模型**。

![](/imgs/learn-qt/img-09.png)

## 基础规则与约束

1. 同一个信号可以绑定多个槽
2. 多个信号可以绑定同一个槽
3. 默认**同步执行**：信号发出后，等待槽函数全部执行完毕才返回
4. 信号参数类型必须和槽函数参数类型兼容

connect 标准签名：

```python
sender.signal.connect(receiver.slot)
```

## 两种完整写法

**写法 1：函数引用写法（内置信号首选，推荐）**

直接引用信号对象、槽函数对象，支持静态检查，IDE 支持跳转。

```python
from PySide6.QtWidgets import QApplication, QPushButton, QWidget, QVBoxLayout


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("信号 & 槽")
        self.resize(300, 200)
        layout = QVBoxLayout()
        self.btn = QPushButton("点击我")
        layout.addWidget(self.btn)
        self.setLayout(layout)

        self.btn.clicked.connect(self.on_click)

    def on_click(self):
        print("slot 被调用")


app = QApplication([])
window = MyWindow()
window.show()
app.exec()
```

运行如下：

![](/imgs/learn-qt/img-10.png)

点击按钮，终端打印：

```bash
slot 被调用
```

**写法 2：自定义 Signal + @Slot 装饰器**

内置信号不足完成需求时，可以自定义信号并绑定到槽函数。

> `@Slot` 是可选但推荐：给 Qt 元对象系统注册槽，优化性能、支持重载区分、跨线程稳定。
> 信号必须在**类级别**定义，不能写在 `__init__` 里面！

```python
from PySide6.QtWidgets import QApplication, QPushButton, QWidget, QVBoxLayout
from PySide6.QtCore import Signal, Slot


class SenderWidget(QWidget):
    # 类级别声明自定义信号：括号内是参数类型
    my_signal = Signal(str, int)

    def __init__(self):
        super().__init__()
        layout = QVBoxLayout()
        self.btn = QPushButton("发送自定义信号")
        self.btn.clicked.connect(self.emit_signal)
        layout.addWidget(self.btn)
        self.setLayout(layout)

    def emit_signal(self):
        # 发射信号，传参
        self.my_signal.emit("Hello", 520)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("信号 & 槽")
        self.resize(300, 200)
        layout = QVBoxLayout()
        self.sender = SenderWidget()
        layout.addWidget(self.sender)
        self.setLayout(layout)
        # 绑定自定义信号到槽函数
        self.sender.my_signal.connect(self.on_my_signal)

    @Slot(str, int)
    def on_my_signal(self, msg, num):
        print(f"自定义信号收到：{msg}, {num}")


app = QApplication([])
window = MyWindow()
window.show()
app.exec()
```

运行效果如下：

![](/imgs/learn-qt/img-11.png)

点击按钮，终端打印：

```bash
自定义信号收到：Hello, 520
```

## 信号重载

一个信息可以 connect 多个槽函数，每个槽函数可以接收不同的参数。

```python
import sys

from PySide6.QtCore import QObject, Signal, Slot
from PySide6.QtWidgets import QApplication


class DemoObject(QObject):
    # 重载信号：传入多个元组，每个元组代表一组参数签名
    data = Signal((int,), (str,))

    @Slot(int)
    def on_data_int(self, val: int):
        print(f"[int版本] 收到数字: {val}")

    @Slot(str)
    def on_data_str(self, val: str):
        print(f"[str版本] 收到字符串: {val}")


app = QApplication(sys.argv)
obj = DemoObject()

# 下标选择指定重载签名
obj.data[int].connect(obj.on_data_int)
obj.data[str].connect(obj.on_data_str)

obj.data[int].emit(520)
obj.data[str].emit("Hello")

sys.exit(app.exec())
```

* data = Signal((int,), (str,))，支持 int 或 str 类型参数，模版化生成两个信号
* 通过下标为不同版本的信息绑定不同的槽函数

运行效果如下：

```bash
[int版本] 收到数字: 520
[str版本] 收到字符串: Hello
```

## disconnect 解除绑定

信号与槽函数通过 `connect` 方法绑定，通过 `disconnect` 方法解除绑定。从而实现动态绑定与解绑。

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QMainWindow, QPushButton, QVBoxLayout, QLabel
)
from PySide6.QtCore import Slot, QObject, Signal, Qt


class MyObject(QObject):
    signal = Signal()

    @Slot()
    def call_slot(self):
        print("call_slot 被调用")


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("connect / disconnect - 动态绑定")
        self.resize(300, 200)

        self.obj = MyObject()

        central = QWidget()
        self.setCentralWidget(central)
        layout = QVBoxLayout(central)

        self.label = QLabel("点击窗口空白处触发信号", alignment=Qt.AlignmentFlag.AlignCenter)
        layout.addWidget(self.label)

        self.btn_connect = QPushButton("信号绑定 connect")
        self.btn_disconnect = QPushButton("解除绑定 disconnect")

        layout.addWidget(self.btn_connect)
        layout.addWidget(self.btn_disconnect)

        self.btn_connect.clicked.connect(self.do_connect)
        self.btn_disconnect.clicked.connect(self.do_disconnect)

    def do_connect(self):
        # 动态绑定
        self.obj.signal.connect(self.obj.call_slot)
        print("已 connect：点击窗口空白处触发的信号现在会调用 call_slot")

    def do_disconnect(self):
        # 动态解除绑定
        self.obj.signal.disconnect(self.obj.call_slot)
        print("已 disconnect：信号不再调用 call_slot")

    def mousePressEvent(self, event):
        # 点击主窗口空白处发射信号，用来验证绑定/解除是否生效
        self.obj.signal.emit()
        super().mousePressEvent(event)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

运行效果如下：

![](/imgs/learn-qt/img-12.png)

点击效果如下：

```bash
已 connect：点击窗口空白处触发的信号现在会调用 call_slot
call_slot 被调用
已 disconnect：信号不再调用 call_slot
```

## 小结

* 信号与槽函数通过 `connect` 方法绑定，通过 `disconnect` 方法解除绑定
* 信号必须是类属性，不能是实例属性
* 支持信号重载，每个信号可以 connect 多个槽函数，每个槽函数可以接收不同的参数
* 支持动态绑定与解绑，实现信号与槽函数的灵活连接与断开
