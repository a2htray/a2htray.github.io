+++
date = '2026-11-12T14:53:17+08:00'
draft = false
title = 'PySide6 弹窗对话框（二）：QFileDialog、QColorDialog'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '弹窗', 'QFileDialog', 'QColorDialog']
toc = true
+++

![](/imgs/learn-qt/cover14.png)

`QFileDialog` 控件是**文件 / 目录**选择对话框，用来和本地文件系统交互，使用场景包括：

1. 编辑器：打开 txt、md、csv、json；导出文件
2. 图片工具：加载图片、保存图片
3. 数据处理程序：选择数据集、选择输出目录
4. 配置工具：选择资源文件夹、模型权重文件
5. 导入导出：导入 Excel、导出报表
6. 批量处理：选择多个文件批量转换

`QColorDialog` 控件是**颜色拾取**对话框，用来选取 RGB/HSV 颜色，支持透明度通道，使用场景包括：

1. 绘图软件：画笔颜色、填充色、线条颜色
2. 富文本编辑器：文字颜色、背景高亮色
3. 界面配置面板：自定义主题、按钮颜色、标签背景色
4. 可视化图表：修改曲线 / 柱状图颜色
5. 标注工具：标注框颜色、分类标签配色
6. 样式配置：自定义 Palette、QSS 预览取色

## QFileDialog

`QFileDialog` 控件用于选文件、选目录、存文件，3 个常用静态方法分别是：

1. `getOpenFileName`：打开文件对话框，返回选中的文件路径
2. `getOpenFileNames`：打开文件对话框，返回选中的文件路径列表
3. `getSaveFileName`：保存文件对话框，返回保存路径

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QPushButton, QFileDialog


class FileDialogDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QFileDialog 三个静态方法")
        layout = QVBoxLayout()

        btn_open_single = QPushButton("getOpenFileName 单选文件")
        btn_open_single.clicked.connect(self.on_open_single)
        layout.addWidget(btn_open_single)

        btn_open_multi = QPushButton("getOpenFileNames 多选文件")
        btn_open_multi.clicked.connect(self.on_open_multi)
        layout.addWidget(btn_open_multi)

        btn_save_single = QPushButton("getSaveFileName 保存文件")
        btn_save_single.clicked.connect(self.on_save_single)
        layout.addWidget(btn_save_single)

        self.setLayout(layout)

    def on_open_single(self):
        filepath, _ = QFileDialog.getOpenFileName(
            self,
            caption="打开单个文件",
            dir="",
            filter="所有文件 (*);;文本文件 (*.txt)"
        )
        print("getOpenFileName ->", filepath)

    def on_open_multi(self):
        filelist, _ = QFileDialog.getOpenFileNames(
            self,
            caption="选择多个文件",
            dir="",
            filter="所有文件 (*);;文本文件 (*.txt)"
        )
        print("getOpenFileNames ->", filelist)

    def on_save_single(self):
        # getSaveFileName：另存为对话框
        filepath, _ = QFileDialog.getSaveFileName(
            self,
            caption="保存文件",
            dir="",
            filter="所有文件 (*);;文本文件 (*.txt)"
        )
        print("getSaveFileName ->", filepath)


app = QApplication(sys.argv)
window = FileDialogDemo()
window.resize(300, 220)
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-59.png)

![](/imgs/learn-qt/img-60.png)

![](/imgs/learn-qt/img-61.png)

过滤器是一个用双分号分隔的字符串，每条规则由“显示名”和括号里的通配符组成，同一类里的多个后缀用空格分隔：

```bash
# 一条规则：显示名 (*.后缀)
"PNG 图片 (*.png)"

# 同一类多个后缀用空格分隔
"图片 (*.png *.jpg *.jpeg *.webp)"

# 多条规则用两个分号分隔，下拉框里各占一项
"图片 (*.png *.jpg);;文本 (*.txt);;所有文件 (*)"
```

静态方法覆盖的是常见形态，多选、默认后缀、列表视图、保存模式这些要逐项设置时，改用实例写法，如：

```python
dlg = QFileDialog(self, "打开文件")
dlg.setAcceptMode(QFileDialog.AcceptMode.AcceptOpen)
dlg.setFileMode(QFileDialog.FileMode.ExistingFile)
dlg.setNameFilter("文本文件 (*.txt);;JSON (*.json);;所有文件 (*)")
dlg.setDefaultSuffix("txt")
dlg.setViewMode(QFileDialog.ViewMode.List)
```

## QColorDialog

`QColorDialog` 控件用于取一个颜色对话框，返回选中的颜色值，核心的静态方法是 `getColor`：

- 返回 `QColor`；**取消返回无效颜色，必须用 `color.isValid()` 判断**

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QPushButton, QLabel, QVBoxLayout, QColorDialog
from PySide6.QtGui import QColor


class Demo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QColorDialog 示例")
        self.btn = QPushButton("选择颜色")
        self.label = QLabel("颜色预览")
        self.label.setFixedSize(200, 80)
        self.label.setStyleSheet("background-color: #888888;")

        lay = QVBoxLayout(self)
        lay.addWidget(self.btn)
        lay.addWidget(self.label)

        self.btn.clicked.connect(self.pick_color)

    def pick_color(self):
        # 初始红色，开启Alpha通道
        c = QColorDialog.getColor(
            initial=QColor(255,0,0),
            parent=self,
            title="选择颜色",
            options=QColorDialog.ColorDialogOption.ShowAlphaChannel
        )
        if c.isValid():
            self.label.setStyleSheet(f"background-color: {c.name(QColor.HexArgb)};")
            print(f"hex: {c.name(QColor.HexArgb)} | rgba: {c.red()},{c.green()},{c.blue()},{c.alpha()}")


app = QApplication(sys.argv)
window = Demo()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-62.png)

## 小结

`QFileDialog` 是 Qt 的文件选择对话框，支持文件、文件夹选取与筛选，返回选中路径，需判断路径有效性；`QColorDialog` 是颜色拾取对话框，可选择 RGB/RGBA 颜色，支持实时预览与自定义色板，二者都提供静态便捷调用方法，也支持实例化绑定信号实现非阻塞交互，开发时留意平台原生对话框差异，取消操作会返回无效结果，务必做合法性校验。
