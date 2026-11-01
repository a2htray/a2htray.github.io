+++
date = '2026-10-29T17:28:37+08:00'
draft = false
title = 'PySide6 窗口基础：QWidget / QMainWindow / QDialog 三大顶层窗口区别'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI']
toc = true
+++

![](/imgs/learn-qt/cover1.png)

## 前言

在 Qt 的世界里，几乎所有界面元素 —— 按钮、输入框、标签、甚至窗口都是 `QWidget`（控件）的子类。

但当我们写一个真正的程序时，顶层窗口通常有三种选择：

* **QWidget**
* **QMainWindow**
* **QDialog**

本文主要介绍三个类的**职责分工、各自能力、适用场景和坑**，并给出各自最小可运行代码。

## 核心概念：三者到底是什么关系

> **QDialog 和 QMainWindow 都是 QWidget 的子类。**

三者都是**可以独立显示的顶层窗口**，但设计之初的定位及适用场景却不同：

| 类 | 设计定位 | 最擅长的场景 |
|---|---|---|
| `QWidget` | 最通用的窗口/控件基类，无任何预设结构 | 自定义控件、弹层、工具小窗、临时画布 |
| `QMainWindow` | 完整的主窗口骨架，自带菜单栏、工具栏、状态栏 | 应用程序的主窗口 |
| `QDialog` | 对话框，与父窗口协作交互 | 设置弹窗、确认框、输入框 |

### 继承关系

```
QObject
 └── QWidget              ← 一切控件与窗口的基类
      ├── QMainWindow     ← 主窗口：菜单栏/工具栏/状态栏
      └── QDialog         ← 对话框：模态交互
```

## QWidget：最朴素的空白画布

`QWidget` 就是一个"可以放东西的矩形区域"，**不带任何预设的菜单栏、状态栏结构**。你可以在里面放布局和子控件，自己组装一切。

**适用场景**

- 一个自定义控件
- 一个极简的、不需要菜单栏/状态栏的工具小窗
- 作为 `QMainWindow` 的 centralWidget 内容区
- 弹出一个"非交互式"的小窗口（如提示浮层）

**自定义控件**

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QLabel
from PySide6.QtCore import Qt


class CustomWidget(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("自定义 Widget 控件")
        self.resize(300, 200)

        layout = QVBoxLayout(self)
        layout.addWidget(QLabel(
            text="没有菜单栏、状态栏、工具栏",
            alignment=Qt.AlignmentFlag.AlignCenter,
        ))


app = QApplication(sys.argv)
window = CustomWidget()
window.show()

sys.exit(app.exec())
```

![](/imgs/learn-qt/img-04.png)

## QMainWindow：应用主窗口的完整骨架

`QMainWindow` 是**为"应用程序主窗口"量身定制**的。它内置了一套固定的布局骨架，包括：

* 菜单栏
* 工具栏
* 停靠区（左）
* 中央区
* 停靠区（右）
* 状态栏

### 最小可运行示例

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QVBoxLayout,
    QLabel, QToolBar,
    QDockWidget, QTextEdit
)
from PySide6.QtGui import QAction
from PySide6.QtCore import Qt


class DemoMainWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QMainWindow 窗口布局")
        self.resize(800, 600)

        # 1. 菜单栏
        menu_bar = self.menuBar()
        menu = menu_bar.addMenu("菜单栏")
        menu.addAction(QAction("示例", self))

        # 2. 工具栏
        tool_bar = QToolBar("工具栏", self)
        tool_bar.addAction(QAction("工具栏", self))
        self.addToolBar(tool_bar)

        # 3. 停靠区（左）
        left_dock = QDockWidget("左侧停靠区", self)
        left_dock.setWidget(QTextEdit("左侧停靠区"))
        self.addDockWidget(Qt.DockWidgetArea.LeftDockWidgetArea, left_dock)
        left_dock.show()

        # 4. 中央区
        central_widget = QWidget()
        layout = QVBoxLayout(central_widget)
        self.setLayout(layout)
        layout.addWidget(QLabel("中央区", alignment=Qt.AlignmentFlag.AlignCenter))

        # 5. 停靠区（右）
        right_dock = QDockWidget("右侧停靠区", self)
        right_dock.setWidget(QTextEdit("右侧停靠区"))
        self.addDockWidget(Qt.DockWidgetArea.RightDockWidgetArea, right_dock)
        right_dock.show()

        # 6. 状态栏
        status_bar = self.statusBar()
        status_bar.showMessage("就绪 · 状态栏消息", 100000)


app = QApplication(sys.argv)
window = DemoMainWindow()
window.show()
sys.exit(app.exec())
```

![](/imgs/learn-qt/img-05.png)

> Mac 系统菜单栏在界面顶部。


## QDialog：对话框，天生为交互而生

`QDialog` 的设计目标非常明确：**和用户做一次性交互**（确认、选文件、输参数），然后返回结果。它比 `QMainWindow` 轻量，且有 `exec()` / `accept()` / `reject()` 这一整套"模态交互"机制。

**两种运行模式**

| 模式 | 方法 | 行为 |
|---|---|---|
| 模态（阻塞父窗口） | `dialog.exec()` | 程序阻塞在 exec()，用户关掉后才继续 |
| 非模态（不阻塞） | `dialog.show()` | 立即返回，可同时操作主窗口 |

### 最小可运行示例（模态 + 返回值）

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QPushButton,
    QDialog, QLabel, QVBoxLayout, QDialogButtonBox,
)


class MyDialog(QDialog):
    def __init__(self, parent=None):
        super().__init__(parent)
        self.setWindowTitle("设置对话框")

        layout = QVBoxLayout(self)
        layout.addWidget(QLabel("这是一个模态对话框"))

        # 标准按钮组：OK = accept，Cancel = reject
        buttons = QDialogButtonBox(
            QDialogButtonBox.Ok | QDialogButtonBox.Cancel
        )
        buttons.accepted.connect(self.accept)   # 点 OK -> 关闭并返回 True
        buttons.rejected.connect(self.reject)   # 点 Cancel -> 关闭并返回 False
        layout.addWidget(buttons)


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("对话框演示")
        self.resize(600, 400)

        central = QWidget(self)
        layout = QVBoxLayout(central)
        self.btn = QPushButton("打开设置对话框")
        self.btn.clicked.connect(self.open_dialog)
        self.setCentralWidget(central)
        layout.addWidget(self.btn)

    def open_dialog(self):
        dlg = MyDialog(self)
        # 阻塞用户操作，直到用户关闭对话框
        if dlg.exec(): # 用户点了 OK
            print("用户确认")
        else: # 用户点了 Cancel 或关闭
            print("用户取消")


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

![](/imgs/learn-qt/img-06.png)

**代码解析**

- `exec()` 是**模态**核心：阻塞父窗口，返回 `True/False`。
- `accept()` 与 `reject()`：内部会调用 `close()` 并设置返回结果，是"优雅关闭"的标准做法。
- `QDialogButtonBox`：标准对话框按钮组，自动排好 OK/Cancel，省去手动布局。

### 最小可运行示例（非模态）

将 `exec()` 改为 `show()` 即可，不阻塞父窗口操作。

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QMainWindow, QWidget, QPushButton,
    QDialog, QLabel, QVBoxLayout, QDialogButtonBox,
)


class MyDialog(QDialog):
    def __init__(self, parent=None):
        super().__init__(parent)
        self.setWindowTitle("设置对话框")

        layout = QVBoxLayout(self)
        layout.addWidget(QLabel("这是一个模态对话框"))

        # 标准按钮组：OK = accept，Cancel = reject
        buttons = QDialogButtonBox(
            QDialogButtonBox.Ok | QDialogButtonBox.Cancel
        )
        buttons.accepted.connect(self.accept)   # 点 OK -> 关闭并返回 True
        buttons.rejected.connect(self.reject)   # 点 Cancel -> 关闭并返回 False
        layout.addWidget(buttons)


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("对话框演示")
        self.resize(600, 400)

        central = QWidget(self)
        layout = QVBoxLayout(central)
        self.btn = QPushButton("打开设置对话框")
        self.btn.clicked.connect(self.open_dialog)
        self.setCentralWidget(central)
        layout.addWidget(self.btn)

    def open_dialog(self):
        dlg = MyDialog(self)
        dlg.show()


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

## 一张表搞定选择

| 维度 | QWidget | QMainWindow | QDialog |
|---|---|---|---|
| 预设结构 | 无 | 菜单栏/工具栏/状态栏/停靠区 | 无（交互机制） |
| 是否必须有中央控件 | 否 | **是（setCentralWidget）** | 否 |
| 模态运行 | 无 | 无 | `exec()` 模态 / `show()` 非模态 |
| 返回值机制 | 无 | 无 | `accept()/reject()` |
| 典型用途 | 自定义控件/小窗 | 主程序窗口 | 设置/确认/输入弹窗 |
| 最小代码量 | 最少 | 中等 | 中等 |

**选型口诀**

- 程序主窗口 → `QMainWindow`
- 和用户对话的临时框 → `QDialog`
- 一个可复用的"部件" → `QWidget`
