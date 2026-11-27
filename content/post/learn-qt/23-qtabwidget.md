+++
date = '2026-11-23T22:58:39+08:00'
draft = false
title = 'PySide6：QTabWidget 标签页控件'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QTabWidget', '标签页控件']
toc = true
+++

![](/imgs/learn-qt/cover23.png)

`QTabWidget` 是标签页容器，每个标签对应一个页面（`QWidget`），用户可以通过点击标签来切换不同的页面。

## 信号与方法

常用信号包括：

- `currentChanged(int)`：当用户切换到不同的标签时触发，参数为当前选中的标签索引。
- `tabCloseRequested(int)`：当用户点击关闭按钮时触发，参数为要关闭的标签索引。
- `tabBarClicked(int)`：当用户点击标签页时触发，参数为点击的标签索引。

常用方法包括：

- `addTab(QWidget, QString)`：添加一个标签页，参数为标签页的 `QWidget` 和标签页的标题。
- `insertTab(int, QWidget, QString)`：在指定位置插入一个标签页，参数为插入位置、标签页的 `QWidget` 和标签页的标题。
- `removeTab(int)`：删除指定标签页，参数为标签页的索引。
- `setCurrentIndex(int)`：设置当前选中的标签索引，参数为标签页的索引。
- `setTabText(int, QString)`：设置指定标签页的标题，参数为标签页的索引和标题。
- `count()`：返回标签页的数量。
- `tabBar()`：返回标签页栏（`QTabBar`），用于自定义标签页栏的外观。

## 示例

**演示信号**

```python
import sys

from PySide6.QtWidgets import QApplication, QLabel, QTabWidget, QVBoxLayout, QWidget


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QTabWidget 信号演示")
        self.resize(400, 260)

        self.tabs = QTabWidget()
        self.tabs.setTabsClosable(True)
        for i in range(1, 4):
            self.tabs.addTab(QLabel(f"第 {i} 页内容"), f"标签 {i}")

        self.log = QLabel("操作提示会显示在这里")
        self.log.setWordWrap(True)

        self.tabs.currentChanged.connect(
            lambda index: self.log.setText(f"currentChanged({index})")
        )
        self.tabs.tabCloseRequested.connect(self.close_tab)
        self.tabs.tabBarClicked.connect(
            lambda index: self.log.setText(f"tabBarClicked({index})")
        )

        root = QVBoxLayout(self)
        root.addWidget(self.tabs)
        root.addWidget(self.log)

    def close_tab(self, index):
        self.tabs.removeTab(index)
        self.log.setText(f"tabCloseRequested({index})")


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-94.png)

![](/imgs/learn-qt/img-95.png)

**演示方法**

```python
import sys

from PySide6.QtWidgets import (
    QApplication,
    QHBoxLayout,
    QLabel,
    QPushButton,
    QTabWidget,
    QVBoxLayout,
    QWidget,
)


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QTabWidget 方法演示")
        self.resize(560, 260)

        self.tabs = QTabWidget()
        self.tabs.addTab(QLabel("内容 A"), "标签 A")
        self.tabs.addTab(QLabel("内容 B"), "标签 B")

        self.status = QLabel()
        buttons = QHBoxLayout()
        for text, slot in (
            ("addTab()", self.add_tab),
            ("insertTab(1)", self.insert_tab),
            ("removeTab(0)", self.remove_tab),
            ("setCurrentIndex(1)", self.set_current),
            ("setTabText(0)", self.set_text),
            ("count()", self.show_count),
            ("tabBar() 样式", self.style_tab_bar),
        ):
            button = QPushButton(text)
            button.clicked.connect(slot)
            buttons.addWidget(button)

        root = QVBoxLayout(self)
        root.addWidget(self.tabs)
        root.addLayout(buttons)
        root.addWidget(self.status)
        self.show_count()

    def add_tab(self):
        self.tabs.addTab(QLabel("新内容"), f"标签 {self.tabs.count() + 1}")
        self.show_count()

    def insert_tab(self):
        self.tabs.insertTab(1, QLabel("插入内容"), "插入页")
        self.show_count()

    def remove_tab(self):
        if self.tabs.count():
            self.tabs.removeTab(0)
        self.show_count()

    def set_current(self):
        self.tabs.setCurrentIndex(1 if self.tabs.count() > 1 else 0)
        self.show_count()

    def set_text(self):
        if self.tabs.count():
            self.tabs.setTabText(0, "已改名")
        self.show_count()

    def show_count(self):
        self.status.setText(f"count() = {self.tabs.count()}")

    def style_tab_bar(self):
        self.tabs.tabBar().setStyleSheet("""
            QTabBar::tab { background: #ddd; padding: 6px 12px; }
            QTabBar::tab:selected { background: white; }
        """)
        self.status.setText(f"已通过 tabBar() 设置样式；count() = {self.tabs.count()}")


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-96.png)

![](/imgs/learn-qt/img-97.png)

![](/imgs/learn-qt/img-98.png)

![](/imgs/learn-qt/img-99.png)

![](/imgs/learn-qt/img-100.png)

## 小结

`QTabWidget` 把“页面容器 + 标签切换”打包成一个控件：

1. **增删改查**：`addTab` / `insertTab` / `removeTab` 管页面的增删，`setCurrentIndex` / `setTabText` 管选中与改标题，`count()` 随时确认页面数。
2. **信号**：`currentChanged` 监听页面切换，是加载数据、刷新内容的主要入口；可关闭标签需 `setTabsClosable(True)` 并在 `tabCloseRequested` 中自行调用 `removeTab`。
3. **外观定制**：通过 `tabBar()` 拿到底层 `QTabBar`，配合 QSS 的 `QTabBar::tab`、`::tab:selected` 子控件即可美化标签栏。
