+++
date = '2026-11-27T14:28:29+08:00'
draft = false
title = 'PySide6 QMenu、QMenuBar、QContextMenuEvent 菜单栏、右键菜单'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QMenu', 'QMenuBar', 'QContextMenuEvent']
toc = true
+++

![](/imgs/learn-qt/cover27.png)

菜单是桌面应用的标配交互组件：顶部菜单栏收纳全部功能入口，右键菜单提供上下文快捷操作。

PySide6 通过 `QMenuBar`（菜单栏）、`QMenu`（菜单）、`QAction`（动作项）三者配合实现菜单系统，而右键菜单则通过重写 `contextMenuEvent()` 来触发。

## QMenuBar 与 QMenu 介绍

`QMenuBar` 是窗口顶部的菜单栏，`QMenu` 是点击菜单名后弹出的下拉菜单，三者关系：

- `QMenuBar`：容纳多个 `QMenu`，横向排列
- `QMenu`：容纳多个 `QAction`（动作项）、`QSeparator`（分隔线），还可以嵌套 `QMenu`（子菜单）
- `QAction`：菜单项本身，点击后发出 `triggered` 信号，可设图标、快捷键、勾选状态

```python
import sys

from PySide6.QtGui import QAction
from PySide6.QtWidgets import (
    QApplication,
    QMainWindow,
    QPlainTextEdit,
)


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QMenuBar / QMenu 菜单栏演示")
        self.resize(560, 380)

        # 菜单栏：QMainWindow 自带 menuBar()
        menu_bar = self.menuBar()

        # 文件菜单：动作 + 快捷键 + 分隔线
        file_menu = menu_bar.addMenu("文件(&F)")
        open_action = QAction("打开", self)
        open_action.setShortcut("Ctrl+O")
        open_action.triggered.connect(
            lambda: self.log.appendPlainText("打开文件")
        )
        file_menu.addAction(open_action)

        save_action = QAction("保存", self)
        save_action.setShortcut("Ctrl+S")
        save_action.triggered.connect(
            lambda: self.log.appendPlainText("保存文件")
        )
        file_menu.addAction(save_action)

        file_menu.addSeparator()
        exit_action = QAction("退出", self)
        exit_action.triggered.connect(self.close)
        file_menu.addAction(exit_action)

        # 编辑菜单：子菜单 + 勾选项
        edit_menu = menu_bar.addMenu("编辑(&E)")
        format_menu = edit_menu.addMenu("格式")
        wrap_action = QAction("自动换行", self)
        wrap_action.setCheckable(True)
        wrap_action.toggled.connect(
            lambda checked: self.log.appendPlainText(f"自动换行: {checked}")
        )
        format_menu.addAction(wrap_action)

        self.log = QPlainTextEdit()
        self.log.setReadOnly(True)
        self.setCentralWidget(self.log)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

* `self.menuBar()`：`QMainWindow` 自带菜单栏，直接获取使用
* `menu_bar.addMenu("文件(&F)")`：添加菜单，`&F` 是 Alt 助记快捷键
* `QAction("打开", self)`：菜单项，`setShortcut()` 设快捷键，`setCheckable(True)` 支持勾选
* `file_menu.addSeparator()`：添加分隔线
* `edit_menu.addMenu("格式")`：在菜单内嵌套子菜单

> macOS 默认使用系统原生菜单栏（菜单显示在屏幕顶部），如需窗口内显示，调用 `menu_bar.setNativeMenuBar(False)`。

演示效果如下：

![](/imgs/learn-qt/img-113.png)

## 核心 API

QMenuBar 常用方法包括：

- `addMenu(title)`：添加菜单，返回 `QMenu`
- `setNativeMenuBar(bool)`：是否使用系统原生菜单栏

QMenu 常用方法包括：

- `addAction(QAction)`：添加动作项
- `addSeparator()`：添加分隔线
- `addMenu(title)`：添加子菜单
- `exec(pos)`：在指定位置模态弹出，返回被点击的 `QAction`（未选择返回 `None`）
- `popup(pos)`：非模态弹出

QAction 常用信号包括：

- `triggered(bool)`：动作被点击
- `toggled(bool)`：勾选状态变化（需 `setCheckable(True)`）

## QContextMenuEvent 右键菜单

Qt 控件默认没有右键菜单，需要重写 `contextMenuEvent()` 自己弹出，核心就两步：

1. 重写 `contextMenuEvent(self, event)`，右键点击时自动触发
2. 构造 `QMenu`，调用 `menu.exec(event.globalPos())` 在鼠标位置模态弹出

```python
import sys

from PySide6.QtGui import QAction
from PySide6.QtWidgets import (
    QApplication,
    QMainWindow,
    QMenu,
    QPlainTextEdit,
)


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QContextMenuEvent 右键菜单演示")
        self.resize(560, 380)

        self.log = QPlainTextEdit("在窗口任意位置点击右键试试")
        self.setCentralWidget(self.log)

    def contextMenuEvent(self, event):
        menu = QMenu(self)

        copy_action = QAction("复制全部", self)
        copy_action.triggered.connect(self.copy_all)
        menu.addAction(copy_action)

        clear_action = QAction("清空", self)
        clear_action.triggered.connect(self.log.clear)
        menu.addAction(clear_action)

        menu.addSeparator()
        about_action = QAction("关于", self)
        about_action.triggered.connect(
            lambda: self.log.appendPlainText("QMenu 右键菜单示例")
        )
        menu.addAction(about_action)

        # 在鼠标位置模态弹出，返回被点击的 QAction
        menu.exec(event.globalPos())

    def copy_all(self):
        self.log.selectAll()
        self.log.copy()


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

* `event.globalPos()`：鼠标的全局坐标，作为菜单弹出位置
* `menu.exec(event.globalPos())`：模态弹出并阻塞，直到用户选择或关闭菜单
* `exec()` 的返回值是被点击的 `QAction`，也可以不用信号、直接根据返回值分支处理

菜单栏的 `QMenu` 由 `QMenuBar` 长期持有，而右键菜单用完即弃，定义为局部变量即可，无需保存引用。

演示效果如下：

![](/imgs/learn-qt/img-114.png)

## 小结


1. **菜单栏体系**：`QMenuBar` 装菜单，`QMenu` 装动作，`QAction` 是功能入口，`addMenu` / `addAction` 层层挂载，配合快捷键、勾选、子菜单覆盖常规需求。
2. **右键菜单**：重写 `contextMenuEvent()`，用 `menu.exec(event.globalPos())` 在鼠标处模态弹出，局部 `QMenu` 用完即弃。
3. **复用技巧**：`exec()` 会返回被点击的 `QAction`，简单菜单可以直接判断返回值，不必逐个连接信号。
