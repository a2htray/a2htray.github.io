+++
date = '2026-11-27T13:33:56+08:00'
draft = false
title = 'PySide6 QSplitter 可拖动分割面板，左右/上下分栏'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QSplitter']
toc = true
+++

![](/imgs/learn-qt/cover25.png)

在桌面 GUI 开发中，可自由拖动的分割面板是非常高频的需求，如：

- 代码编辑器左右分栏
- 资源管理器侧边栏
- 上下预览面板
- 三栏结构

Qt 自带的 `QSplitter` 就是专门用来实现该效果的控件，**不需要自己监听鼠标拖拽、绘制分割线，开箱即用**，同时支持水平、垂直以及多层嵌套布局。

## QSplitter 介绍

`QSplitter` 是 Qt 提供的分割容器控件，核心特点：

1. 两种分割方向：**垂直分割（上下分栏）**、**水平分割（左右分栏）**
2. 内置拖拽逻辑，鼠标拖动分割手柄即可动态调整面板大小
3. 支持多层嵌套，轻松实现复杂多栏布局

**垂直分割（上下分栏）**

```python
import sys

from PySide6.QtCore import Qt
from PySide6.QtWidgets import QApplication, QLabel, QSplitter, QTextEdit, QVBoxLayout, QWidget


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSplitter 上下分栏")
        self.resize(500, 360)

        top = QLabel("上半部分\n拖动中间分隔条调整高度")
        top.setAlignment(Qt.AlignmentFlag.AlignCenter)
        top.setStyleSheet("background: #e8f0fe;")

        bottom = QTextEdit("下半部分")
        bottom.setStyleSheet("background: #fff3e0;")

        splitter = QSplitter(Qt.Orientation.Vertical)
        splitter.addWidget(top)
        splitter.addWidget(bottom)
        splitter.setStretchFactor(0, 1)
        splitter.setStretchFactor(1, 2)

        root = QVBoxLayout(self)
        root.addWidget(splitter)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-104.png)


**水平分割（左右分栏）**

```python
import sys

from PySide6.QtCore import Qt
from PySide6.QtWidgets import QApplication, QLabel, QSplitter, QTextEdit, QVBoxLayout, QWidget


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSplitter 左右分栏")
        self.resize(500, 360)

        left = QLabel("左侧部分\n拖动中间分隔条调整宽度")
        left.setAlignment(Qt.AlignmentFlag.AlignCenter)
        left.setStyleSheet("background: #e8f0fe;")

        right = QTextEdit("右侧部分")
        right.setStyleSheet("background: #fff3e0;")

        splitter = QSplitter(Qt.Orientation.Horizontal)
        splitter.addWidget(left)
        splitter.addWidget(right)
        splitter.setStretchFactor(0, 1)
        splitter.setStretchFactor(1, 2)

        root = QVBoxLayout(self)
        root.addWidget(splitter)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-105.png)


## 核心 API

常用方法包括：

- `QSplitter(Qt.Orientation)`：构造分割容器，`Vertical` 上下分栏、`Horizontal` 左右分栏
- `setOrientation(Qt.Orientation)`：动态调整分割方向
- `addWidget(QWidget)`：在末尾追加一个子面板
- `insertWidget(int, QWidget)`：在指定索引位置插入一个子面板
- `setStretchFactor(int index, int stretch)`：设置面板拉伸因子，窗口缩放时按比例分配空间
- `setSizes(list[int])`：显式指定各面板的初始像素大小
- `sizes()`：返回当前各面板的像素大小，常用于保存/恢复布局
- `setChildrenCollapsible(bool)`：设置子面板能否被拖拽折叠为 0（默认允许折叠）
- `setHandleWidth(int)`：设置分割手柄的宽度（像素）
- `setOpaqueResize(bool)`：拖动时是否实时调整面板大小，`False` 时先显示预览线、松手后重排

常用信号：

- `splitterMoved(int pos, int index)`：分割手柄被拖动后触发，参数为手柄新位置和面板索引，适合用来保存布局状态

```python
import sys

from PySide6.QtCore import Qt
from PySide6.QtWidgets import (
    QApplication,
    QCheckBox,
    QComboBox,
    QFormLayout,
    QHBoxLayout,
    QLabel,
    QPlainTextEdit,
    QPushButton,
    QSpinBox,
    QSplitter,
    QVBoxLayout,
    QWidget,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSplitter 方法与信号演示")
        self.resize(640, 480)

        self.splitter = QSplitter(Qt.Orientation.Vertical)
        self.splitter.addWidget(self.make_panel("面板 1"))
        self.splitter.addWidget(self.make_panel("面板 2"))
        self.splitter.splitterMoved.connect(self.on_splitter_moved)

        orientation = QComboBox()
        orientation.addItems(["Vertical", "Horizontal"])
        orientation.currentTextChanged.connect(self.change_orientation)

        handle_width = QSpinBox()
        handle_width.setRange(1, 20)
        handle_width.setValue(4)
        handle_width.valueChanged.connect(self.splitter.setHandleWidth)

        collapsible = QCheckBox("setChildrenCollapsible(True)")
        collapsible.setChecked(True)
        collapsible.toggled.connect(self.splitter.setChildrenCollapsible)

        opaque = QCheckBox("setOpaqueResize(True)")
        opaque.setChecked(True)
        opaque.toggled.connect(self.splitter.setOpaqueResize)

        settings = QFormLayout()
        settings.addRow("setOrientation()", orientation)
        settings.addRow("setHandleWidth()", handle_width)
        settings.addRow("", collapsible)
        settings.addRow("", opaque)

        buttons = QHBoxLayout()
        for text, slot in (
            ("addWidget()", self.add_panel),
            ("insertWidget(1)", self.insert_panel),
            ("setStretchFactor(0, 2)", self.set_stretch),
            ("setSizes([100, 200])", self.set_sizes),
            ("sizes()", self.show_sizes),
        ):
            button = QPushButton(text)
            button.clicked.connect(slot)
            buttons.addWidget(button)

        self.log = QPlainTextEdit()
        self.log.setReadOnly(True)
        self.log.setMaximumHeight(100)

        root = QVBoxLayout(self)
        root.addLayout(settings)
        root.addLayout(buttons)
        root.addWidget(self.splitter, 1)
        root.addWidget(self.log)

    def make_panel(self, text):
        label = QLabel(text)
        label.setAlignment(Qt.AlignmentFlag.AlignCenter)
        label.setStyleSheet("background: #e8f0fe;")
        return label

    def change_orientation(self, text):
        value = Qt.Orientation.Vertical if text == "Vertical" else Qt.Orientation.Horizontal
        self.splitter.setOrientation(value)

    def add_panel(self):
        self.splitter.addWidget(self.make_panel(f"面板 {self.splitter.count() + 1}"))

    def insert_panel(self):
        self.splitter.insertWidget(1, self.make_panel("插入面板"))

    def set_stretch(self):
        self.splitter.setStretchFactor(0, 2)

    def set_sizes(self):
        self.splitter.setSizes([100, 200])

    def show_sizes(self):
        self.log.appendPlainText(f"sizes() = {self.splitter.sizes()}")

    def on_splitter_moved(self, pos, index):
        self.log.appendPlainText(f"splitterMoved({pos}, {index})")


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-106.png)


## 小结

1. **方向与面板管理**：构造时用 `Vertical` / `Horizontal` 指定分栏方向，`setOrientation()` 可动态切换；`addWidget` / `insertWidget` 管面板的增删。
2. **大小分配**：`setStretchFactor` 负责**窗口缩放时的比例**，`setSizes` 管**初始大小**，`sizes()` 读回当前布局，三者配合可实现布局的保存与恢复。
3. **手柄行为与信号**：`setHandleWidth`、`setChildrenCollapsible`、`setOpaqueResize` 细调拖拽；`splitterMoved` 在拖动结束后触发，是持久化布局的挂载点。
