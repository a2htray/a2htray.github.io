+++
date = '2026-11-24T09:20:44+08:00'
draft = false
title = 'PySide6：QStackedWidget 堆栈窗口'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QTabWidget', '标签页控件']
toc = true
+++

![](/imgs/learn-qt/cover24.png)

**QStackedWidget 是一个堆栈式页面容器**，里面可以放多个子页面（QWidget），**所有页面是叠在一起的**，同一时间只会显示其中一个，其余页面自动隐藏。

`QStackedWidget` **没有自带导航栏、标签按钮**，本身不会提供切换入口，切换页面需要你自己通过代码调用 `setCurrentIndex` / `setCurrentWidget` 来控制。适合的场景有：

- 左侧菜单
- 侧边导航
- 自定义标签面板

## 信号与方法

常用信号包括：

- `currentChanged(int)`：切换页面时触发，传入新索引
- `widgetRemoved(int)`：页面被移除时触发，传入索引

常用方法包括：

- `addWidget(QWidget)`：添加一个页面，参数为页面 `QWidget`
- `insertWidget(int, QWidget)`：在指定位置插入一个页面，参数为插入位置和页面 `QWidget`
- `removeWidget(QWidget)`：移除指定页面，参数为页面 `QWidget`
- `setCurrentIndex(int)`：设置当前页面索引，参数为页面索引
- `setCurrentWidget(QWidget)`：设置当前页面，参数为页面 `QWidget`

## 示例

**演示信号**

```python
import sys

from PySide6.QtWidgets import (
    QApplication,
    QComboBox,
    QHBoxLayout,
    QLabel,
    QPlainTextEdit,
    QPushButton,
    QStackedWidget,
    QVBoxLayout,
    QWidget,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QStackedWidget 信号演示")
        self.resize(420, 280)

        self.stack = QStackedWidget()
        for i in range(1, 4):
            self.stack.addWidget(QLabel(f"这是第 {i} 页"))

        self.combo = QComboBox()
        self.combo.addItems(["第 1 页", "第 2 页", "第 3 页"])
        self.combo.currentIndexChanged.connect(self.stack.setCurrentIndex)

        remove = QPushButton("移除当前页")
        remove.clicked.connect(self.remove_current)

        self.log = QPlainTextEdit()
        self.log.setReadOnly(True)

        self.stack.currentChanged.connect(
            lambda index: self.log.appendPlainText(f"currentChanged({index})")
        )
        self.stack.widgetRemoved.connect(
            lambda index: self.log.appendPlainText(f"widgetRemoved({index})")
        )

        controls = QHBoxLayout()
        controls.addWidget(self.combo)
        controls.addWidget(remove)

        root = QVBoxLayout(self)
        root.addLayout(controls)
        root.addWidget(self.stack)
        root.addWidget(self.log)

    def remove_current(self):
        index = self.stack.currentIndex()
        if index >= 0:
            self.stack.removeWidget(self.stack.widget(index))
            self.combo.removeItem(index)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-101.png)

![](/imgs/learn-qt/img-102.png)

**演示方法**

```python
import sys

from PySide6.QtWidgets import (
    QApplication,
    QHBoxLayout,
    QLabel,
    QPushButton,
    QStackedWidget,
    QVBoxLayout,
    QWidget,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QStackedWidget 方法演示")
        self.resize(560, 220)

        self.stack = QStackedWidget()
        self.stack.addWidget(QLabel("页面 1"))
        self.stack.addWidget(QLabel("页面 2"))

        self.status = QLabel()
        buttons = QHBoxLayout()
        for text, slot in (
            ("addWidget()", self.add_page),
            ("insertWidget(1)", self.insert_page),
            ("removeWidget()", self.remove_page),
            ("setCurrentIndex(1)", self.set_by_index),
            ("setCurrentWidget()", self.set_by_widget),
        ):
            button = QPushButton(text)
            button.clicked.connect(slot)
            buttons.addWidget(button)

        root = QVBoxLayout(self)
        root.addWidget(self.stack)
        root.addLayout(buttons)
        root.addWidget(self.status)
        self.refresh()

    def add_page(self):
        page = QLabel(f"页面 {self.stack.count() + 1}")
        self.stack.addWidget(page)
        self.refresh()

    def insert_page(self):
        self.stack.insertWidget(1, QLabel("插入页"))
        self.refresh()

    def remove_page(self):
        page = self.stack.currentWidget()
        if page:
            self.stack.removeWidget(page)
            page.deleteLater()
        self.refresh()

    def set_by_index(self):
        if self.stack.count() > 1:
            self.stack.setCurrentIndex(1)
        self.refresh()

    def set_by_widget(self):
        if self.stack.count():
            self.stack.setCurrentWidget(self.stack.widget(0))
        self.refresh()

    def refresh(self):
        self.status.setText(
            f"count={self.stack.count()}, currentIndex={self.stack.currentIndex()}"
        )


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-103.png)

## 小结

1. **页面管理**：`addWidget` / `insertWidget` / `removeWidget` 负责页面的增删，删除时记得配合 `page.deleteLater()` 释放页面。
2. **页面切换**：核心就是 `setCurrentIndex` / `setCurrentWidget` 两个方法，索引与控件引用二选一；切换时会触发 `currentChanged` 信号，是页面数据加载的挂载点。
