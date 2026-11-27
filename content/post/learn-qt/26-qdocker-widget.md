+++
date = '2026-11-27T14:05:50+08:00'
draft = false
title = 'PySide6 QDockWidget 停靠窗口：类似 IDE 侧边栏，可拖拽悬浮'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', 'QDockWidget']
toc = true
+++

![](/imgs/learn-qt/cover26.png)

开发桌面 GUI 程序时，可停靠面板是十分常见的交互组件：资源树、属性表单、信息日志面板，既可以固定依附在窗口四周，也能拖拽脱离成为独立悬浮窗口，多个面板还支持堆叠为标签页切换。

PySide6 原生提供 `QDockWidget` 实现该能力，无需手写鼠标拖拽、吸附判定等底层逻辑，开箱即用。

## QDockWidget 介绍

`QDockWidget` 是 Qt 提供的停靠容器控件，核心特点：

1. 四个停靠区域：窗口的上下左右均可停靠，拖拽时自动显示吸附预览
2. 支持浮动：拖出主窗口即变成独立悬浮窗，也可再次拖回停靠
3. 支持堆叠：多个面板拖到同一区域，自动组合成标签页切换
4. 标题栏可定制：通过 `setFeatures()` 控制关闭、浮动、可移动等能力

注意 `QDockWidget` 只能放在 `QMainWindow` 中使用：中央区域由 `setCentralWidget()` 占据，停靠面板通过 `addDockWidget()` 挂到四个停靠区域。

```python
import sys

from PySide6.QtCore import Qt
from PySide6.QtWidgets import QApplication, QDockWidget, QLabel, QMainWindow, QTextEdit


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QDockWidget 演示")
        self.resize(700, 500)

        self.setCentralWidget(QTextEdit("中央区域"))

        features = (
            QDockWidget.DockWidgetFeature.DockWidgetClosable
            | QDockWidget.DockWidgetFeature.DockWidgetMovable
            | QDockWidget.DockWidgetFeature.DockWidgetFloatable
        )

        def add_dock(title, area):
            dock = QDockWidget(title, self)
            dock.setWidget(QLabel(f"{title} 内容"))
            dock.setFeatures(features)
            self.addDockWidget(area, dock)
            return dock

        add_dock("左侧面板", Qt.DockWidgetArea.LeftDockWidgetArea)
        add_dock("右侧面板", Qt.DockWidgetArea.RightDockWidgetArea)
        add_dock("上方面板", Qt.DockWidgetArea.TopDockWidgetArea)

        bottom1 = add_dock("下方面板 1", Qt.DockWidgetArea.BottomDockWidgetArea)
        bottom2 = add_dock("下方面板 2", Qt.DockWidgetArea.BottomDockWidgetArea)
        self.tabifyDockWidget(bottom1, bottom2)


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

* `QDockWidget("标题")`：构造面板，参数为标题栏文字
* `setWidget(QWidget)`：设置面板内容
* `addDockWidget(area, dock)`：将面板添加到指定停靠区域
* 运行后拖拽面板标题栏可改变停靠位置，拖出主窗口则自动变浮动窗口

演示效果如下：

![](/imgs/learn-qt/img-107.png)

![](/imgs/learn-qt/img-108.png)

![](/imgs/learn-qt/img-109.png)

## 核心 API

QDockWidget 常用方法包括：

- `setWidget(QWidget)`：设置面板内容
- `widget()`：获取面板内容控件
- `setFeatures(...)`：设置标题栏能力，`DockWidgetClosable` 可关闭、`DockWidgetMovable` 可移动、`DockWidgetFloatable` 可浮动、`NoDockWidgetFeatures` 全部禁用
- `setAllowedAreas(Qt.DockWidgetArea)`：限制可停靠区域，拖到不允许的区域不会吸附
- `setFloating(bool)` / `isFloating()`：设置或查询浮动状态，`True` 时面板脱离主窗口成为独立悬浮窗
- `setTitleBarWidget(QWidget)`：用自定义控件替换标题栏

常用信号包括：

- `visibilityChanged(bool)`：面板显示/隐藏时触发（含被关闭）
- `topLevelChanged(bool)`：停靠 ↔ 浮动切换时触发，参数为是否处于浮动状态
- `dockLocationChanged(Qt.DockWidgetArea)`：停靠区域变化时触发

```python
import sys

from PySide6.QtCore import Qt
from PySide6.QtWidgets import (
    QApplication,
    QDockWidget,
    QLabel,
    QMainWindow,
    QPlainTextEdit,
    QPushButton,
    QVBoxLayout,
    QWidget,
)


class MyWindow(QMainWindow):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QDockWidget 核心 API 演示")
        self.resize(660, 420)

        # 中央区域：主工作区 + 浮动切换按钮
        btn = QPushButton("切换 setFloating()")
        btn.clicked.connect(self.toggle_float)
        box = QVBoxLayout()
        box.addWidget(QPlainTextEdit("中央编辑区"))
        box.addWidget(btn)
        wrapper = QWidget()
        wrapper.setLayout(box)
        self.setCentralWidget(wrapper)

        # 左侧面板：限制停靠区域 + 定制标题栏能力
        self.dock = QDockWidget("资源树")
        self.dock.setWidget(QLabel("项目文件列表"))
        self.dock.setAllowedAreas(
            Qt.DockWidgetArea.LeftDockWidgetArea | Qt.DockWidgetArea.RightDockWidgetArea
        )
        self.dock.setFeatures(
            QDockWidget.DockWidgetFeature.DockWidgetClosable
            | QDockWidget.DockWidgetFeature.DockWidgetMovable
            | QDockWidget.DockWidgetFeature.DockWidgetFloatable
        )
        self.addDockWidget(Qt.DockWidgetArea.LeftDockWidgetArea, self.dock)

        # 底部日志面板：记录布局变化事件
        self.log = QPlainTextEdit()
        self.log.setReadOnly(True)
        log_dock = QDockWidget("事件日志")
        log_dock.setWidget(self.log)
        self.addDockWidget(Qt.DockWidgetArea.BottomDockWidgetArea, log_dock)

        # 三个布局相关信号
        self.dock.visibilityChanged.connect(
            lambda visible: self.log.appendPlainText(f"visibilityChanged({visible})")
        )
        self.dock.topLevelChanged.connect(
            lambda floating: self.log.appendPlainText(f"topLevelChanged({floating})")
        )
        self.dock.dockLocationChanged.connect(
            lambda area: self.log.appendPlainText(f"dockLocationChanged({area})")
        )

    def toggle_float(self):
        self.dock.setFloating(not self.dock.isFloating())


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

* `setFeatures()`：限制标题栏能力后，对应按钮（关闭、浮动）不再显示
* `setAllowedAreas()`：拖拽到被禁止的区域时不会出现吸附预览
* `setFloating(True)`：无需拖拽，直接把面板变成浮动窗口
* 布局持久化：配合 `QMainWindow` 的 `saveState()` / `restoreState()`，把所有面板的位置、大小、堆叠关系保存为 `QByteArray`，结合 `QSettings` 可在下次启动时完整恢复布局

演示效果如下：

![](/imgs/learn-qt/img-110.png)

![](/imgs/learn-qt/img-111.png)

![](/imgs/learn-qt/img-112.png)

## 小结

1. **挂载与内容**：面板必须放进 `QMainWindow`，用 `addDockWidget(area, dock)` 挂到停靠区，`setWidget()` 设内容，中央工作区交给 `setCentralWidget()`。
2. **行为约束**：`setFeatures()` 控制标题栏能力（关闭/移动/浮动），`setAllowedAreas()` 限制可吸附区域，代码层 `setFloating()` 可直接切换浮动态。
3. **状态与布局**：`visibilityChanged` / `topLevelChanged` / `dockLocationChanged` 三个信号覆盖面板的全部布局变化；配合 `saveState()` / `restoreState()` + `QSettings` 即可实现布局持久化。
