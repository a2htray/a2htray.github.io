+++
date = '2026-11-03T10:50:32+08:00'
draft = false
title = 'PySide6 信号与槽函数进阶：自定义信号、传参信号、Lambda 槽函数'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'Qt Signal', 'Qt Slot']
toc = true
+++

![](/imgs/learn-qt/cover4.png)

## 概述

在上一篇文章[《PySide6 信号与槽（Signal & Slot）：Qt 最核心机制，写法详解》](/post/learn-qt/03-signal-slot/)中，学习了如何使用内置信号、槽函数绑定、自定义信号，本文将回顾与学习以下一些进阶内容：

- 自定义信号，发送数据（字符串、数字、对象、多个参数）
- Lambda 槽函数，简洁实现槽函数

## 自定义信号，发送复杂结构数据

示例界面交互说明：

1. 3 个按钮，分别用于触发不同的信号，发送字符串、数字、对象
2. 按钮下方有一个文本框，用于显示收到的信号数据

定义发送的对象及包含信号的类：

```python
@dataclass
class User:
    name: str
    age: int
    email: str

    def __str__(self):
        return f"User(name={self.name!r}, age={self.age}, email={self.email!r})"


class Emitter(QObject):
    # 1) 字符串信号：字符串内容 + 发送时间 + 字符串长度（多个参数）
    string_signal = Signal(str, str, int)
    # 2) 数字信号：数值 + 标签 + 单位（多个参数）
    number_signal = Signal(float, str, str)
    # 3) 对象信号：自定义对象 + 备注说明（多个参数）
    object_signal = Signal(object, str)

    def emit_string(self):
        text = "这是一个字符串"
        self.string_signal.emit(text, datetime.now().strftime("%H:%M:%S"), len(text))

    def emit_number(self):
        # 发送一个数字
        value = 3.14159
        self.number_signal.emit(value, "圆周率", "rad")

    def emit_object(self):
        # 发送一个对象
        user = User("Alice", 30, "alice@example.com")
        self.object_signal.emit(user, "新建用户事件")
```

* `User` 是数据类，包含 `name`、`age`、`email` 三个属性
* `Emitter` 类定义了 3 个信号：
    * `string_signal`：演示字符串、多参数发送
    * `number_signal`：演示数字，多参数发送
    * `object_signal`：演示对象，多参数发送

桌面程序代码如下：

```python
import sys
from dataclasses import dataclass
from datetime import datetime

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QVBoxLayout, QHBoxLayout,
    QPushButton, QTextEdit,
)
from PySide6.QtCore import QObject, Signal, Slot


class MainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("PySide6 自定义信号示例")
        self.resize(500, 380)

        self.emitter = Emitter()

        central = QWidget()
        self.setCentralWidget(central)
        layout = QVBoxLayout(central)

        # 顶部：三个发送按钮
        btn_row = QWidget()
        btn_layout = QHBoxLayout(btn_row)
        btn_str = QPushButton("发送字符串信号")
        btn_num = QPushButton("发送数字信号")
        btn_obj = QPushButton("发送对象信号")
        btn_layout.addWidget(btn_str)
        btn_layout.addWidget(btn_num)
        btn_layout.addWidget(btn_obj)
        layout.addWidget(btn_row)

        # 底部：接收日志
        self.log_edit = QTextEdit()
        self.log_edit.setReadOnly(True)
        layout.addWidget(self.log_edit, stretch=1)

        # 信号 -> 槽 连接
        btn_str.clicked.connect(self.emitter.emit_string)
        btn_num.clicked.connect(self.emitter.emit_number)
        btn_obj.clicked.connect(self.emitter.emit_object)

        self.emitter.string_signal.connect(self.on_string)
        self.emitter.number_signal.connect(self.on_number)
        self.emitter.object_signal.connect(self.on_object)

        self.log("点击上方按钮，观察不同数据结构的自定义信号")

    def log(self, msg):
        ts = datetime.now().strftime("%H:%M:%S")
        self.log_edit.append(f"[{ts}] {msg}")

    @Slot(str, str, int)
    def on_string(self, text, time, length):
        self.log(f"收到字符串信号 -> 内容: {text!r}, 时间: {time}, 长度: {length}")

    @Slot(float, str, str)
    def on_number(self, value, label, unit):
        self.log(f"收到数字信号 -> {label} = {value} {unit}")

    @Slot(object, str)
    def on_object(self, obj, note):
        self.log(f"收到对象信号 -> {note}: {obj}")


app = QApplication(sys.argv)
window = MainWindow()
window.show()
app.exec()
```

依次点击 3 个按钮，运行效果如下：

![](/imgs/learn-qt/img-13.png)

简要流程：

```bash
按钮点击 clicked -> 信号触发 emit -> 槽函数执行
```

## Lambda 槽函数

Lambda 省略了函数定义，直接使用匿名函数，适用于简单场景，通常用于实现临时、单次的槽函数。

```python
import sys
from datetime import datetime

from PySide6.QtWidgets import QApplication, QMainWindow, QWidget, QVBoxLayout, QHBoxLayout, QPushButton, QTextEdit


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Lambda 槽函数演示")
        self.resize(400, 200)

        central_widget = QWidget()
        self.setCentralWidget(central_widget)
        layout = QVBoxLayout(central_widget)

        self.btn = QPushButton("点击调用 Lambda 槽函数")
        self.btn.clicked.connect(lambda: self.log("Lambda 槽函数已执行"))
        layout.addWidget(self.btn)

        btns_widget = QWidget()
        btns_layout = QHBoxLayout(btns_widget)

        for i in range(3):
            btn = QPushButton(f"按钮 {i+1}")
            btn.clicked.connect(lambda _, i=i: self.log(f"按钮 {i+1} 被点击了"))
            btns_layout.addWidget(btn)

        layout.addWidget(btns_widget)

        self.text_edit = QTextEdit()
        self.text_edit.setReadOnly(True)
        layout.addWidget(self.text_edit, stretch=1)

    def log(self, msg):
        self.text_edit.append(f"[{datetime.now().strftime('%H:%M:%S')}] {msg}")
        

app = QApplication(sys.argv)
window = MyWindow()
window.show()
app.exec()
```

运行效果如下：

![](/imgs/learn-qt/img-14.png)

上面的示例中，演示了两个 Lambda 槽函数的使用：

```python
self.btn = QPushButton("点击调用 Lambda 槽函数")
self.btn.clicked.connect(lambda: self.log("Lambda 槽函数已执行"))
```

`connect` 直接绑定一个 Lambda 函数，无需定义单独的函数。

```python
for i in range(3):
    btn = QPushButton(f"按钮 {i+1}")
    btn.clicked.connect(lambda _, i=i: self.log(f"按钮 {i+1} 被点击了"))
    btns_layout.addWidget(btn)
```

获取局部变量 `i`，避免闭包问题，方法如下：

```python
lambda _, i=i: self.log(f"按钮 {i+1} 被点击了")
```

> 局部变量作为函数，传入到 Lambda 函数参数中。

## 小结

* 信号触发，可以发送不同类型的数据，支持多个参数
* connect 方法直接绑定 Lambda 函数，将局部变量传入到 Lambda 函数参数中，避免闭包问题
