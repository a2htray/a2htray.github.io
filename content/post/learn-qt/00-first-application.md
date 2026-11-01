+++
date = '2026-10-29T16:43:03+08:00'
draft = false
title = 'PySide6 的第一个应用'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI']
toc = true
+++

PySide，也叫 Qt for Python，是一套 Python 图形界面框架，依托 Qt 工具包来开发 GUI 应用。PySide 是 Qt 官方提供的 Python 绑定库，目前由 Qt 公司自主维护开发。

## 安装

```bash
$ pip install pyside6
```

验证是否安装成功：

```bash
$ python -c "import PySide6; print(PySide6.__version__)"
6.12.0
```

### 创建第一个应用

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget

app = QApplication(sys.argv)

window = QWidget(windowTitle="Hello Python GUI!")
window.show()

app.exec()
```

* `QApplication` 一个应用程序有且仅有一个 `QApplication` 实例，传递 `sys.argv` 是为了让应用程序能够处理命令行参数，例如 `--debug` 或 `--verbose`
* `QWidget` 是一个窗口基类，所有 GUI 窗口都必须继承自 `QWidget`
* `window.show()` 显示窗口，默认是隐藏的
* `app.exec()` 进入事件循环，等待用户交互，直到应用程序退出

运行效果如下：

![](/imgs/learn-qt/img-00.png)

## 事件循环

应用上的每一次交互操作 —— 无论是按键、鼠标点击还是鼠标移动 —— 都会生成一个事件，并存入事件队列。在事件循环中，每一轮迭代都会检查该队列；如果发现有待处理事件，就会将事件连同程序控制权一并交给对应的事件处理器。事件处理器处理完事件后，再把控制权交回事件循环，继续等待后续事件。每个应用同一时刻只会运行一个事件循环。

示意图如下：

![](/imgs/learn-qt/img-01.png)

* `QApplication` 类负责承载整个应用的事件循环

## 自定义窗口

在 QT 中，任何控件都可以是窗口，如下：

```python
import sys

from PySide6.QtWidgets import QApplication, QPushButton

app = QApplication(sys.argv)

window = QPushButton("Hello Python GUI!")
window.show()

app.exec()
```

运行效果如下：

![](/imgs/learn-qt/img-02.png)

PySide6 提供了 `QMainWindow` 类，用于创建主窗口，自定义窗口可以继承自 `QMainWindow`，并在 `__init__` 方法中初始化窗口。

```python
import sys

from PySide6.QtWidgets import QApplication, QMainWindow, QPushButton

app = QApplication(sys.argv)


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("Hello Python GUI!")
        
        button = QPushButton("Click me!")
        self.setCentralWidget(button)


window = MyWindow()
window.show()

app.exec()
```

运行效果如下：

![](/imgs/learn-qt/img-03.png)
