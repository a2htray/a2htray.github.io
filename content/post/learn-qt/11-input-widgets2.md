+++
date = '2026-11-10T09:27:09+08:00'
draft = false
title = 'PySide6 选择输入控件（二）：QDoubleSpinBox、QDateEdit、QTimeEdit'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '选择输入控件', 'QDoubleSpinBox', 'QDateEdit', 'QTimeEdit']
toc = true
+++

![](/imgs/learn-qt/cover11.png)

本文继续讲解数值与时间类选择输入控件：

- **浮点数微调框**：QDoubleSpinBox
- **日期选择框**：QDateEdit
- **时间选择框**：QTimeEdit

这类控件内置**格式化、范围校验**，非常适合表单里的浮点参数、日期、时间录入场景。

## QDoubleSpinBox

`QDoubleSpinBox` 是 `QSpinBox` 的浮点版本，用于输入小数，支持设置小数位数、范围、步长、前后缀。

```python
import sys

from PySide6.QtWidgets import QApplication, QWidget, QVBoxLayout, QDoubleSpinBox, QLabel


class MyWindow(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QDoubleSpinBox 演示")
        self.resize(400, 200)
        layout = QVBoxLayout(self)

        self.spinbox = QDoubleSpinBox()
        self.spinbox.setRange(-10.0, 10.0)
        self.spinbox.setValue(0)
        self.spinbox.setSingleStep(0.1)
        self.spinbox.setPrefix("距离：")
        self.spinbox.setSuffix(" 米")
        self.spinbox.setDecimals(3)
        self.spinbox.valueChanged.connect(self.on_value_changed)
        layout.addWidget(self.spinbox)

        self.label = QLabel(f"当前选中：{self.spinbox.value()}")
        layout.addWidget(self.label)

    def on_value_changed(self, value):
        self.label.setText(f"当前选中：{value}")


app = QApplication(sys.argv)
window = MyWindow()
window.show()
sys.exit(app.exec())
```

`QDoubleSpinBox` 在 `QSpinBox` 的基础上，提供浮点数的输入，新方法包括：

- `setDecimals()`：设置小数位数
- `decimals()`：获取小数位数

演示效果如下：

![](/imgs/learn-qt/img-40.png)

## QDateEdit

`QDateEdit` 用于选择日期，内置日历弹窗，支持自定义日期显示格式、日期范围，底层依赖 `QDate` 对象。

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout,
    QDateEdit, QLabel
)
from PySide6.QtCore import QDate


class DateEditDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QDateEdit 演示")
        self.resize(400, 200)
        layout = QVBoxLayout(self)

        self.label_info = QLabel("选中日期：")
        self.date_edit = QDateEdit()
        self.date_edit.setDate(QDate.currentDate())
        self.date_edit.setDisplayFormat("yyyy-MM-dd")
        self.date_edit.setCalendarPopup(True)
        self.date_edit.setMinimumDate(QDate(2020, 1, 1))
        self.date_edit.setMaximumDate(QDate(2030, 12, 31))

        layout.addWidget(self.label_info)
        layout.addWidget(self.date_edit)

        self.date_edit.dateChanged.connect(self.on_date_changed)

    def on_date_changed(self, qdate: QDate):
        self.label_info.setText(f"选中日期：{qdate.toString('yyyy-MM-dd')}")


app = QApplication(sys.argv)
window = DateEditDemo()
window.show()
sys.exit(app.exec())
```

- `setDate()`：设置日期
- `setDisplayFormat()`：设置日期显示格式
- `setCalendarPopup()`：设置是否弹出日历
- `setMinimumDate()`：设置最小日期
- `setMaximumDate()`：设置最大日期

演示效果如下：

![](/imgs/learn-qt/img-41.png)

弹出日历显示为英文，将其调整为中文：

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout,
    QDateEdit, QLabel
)
from PySide6.QtCore import QDate, QLocale


class DateEditDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QDateEdit 演示")
        self.resize(400, 200)
        layout = QVBoxLayout(self)

        self.label_info = QLabel("选中日期：")
        locale_cn = QLocale(QLocale.Chinese, QLocale.China)
        self.date_edit = QDateEdit()
        self.date_edit.setLocale(locale_cn)
        self.date_edit.setDate(QDate.currentDate())
        self.date_edit.setDisplayFormat("yyyy-MM-dd")
        self.date_edit.setCalendarPopup(True)
        self.date_edit.setMinimumDate(QDate(2020, 1, 1))
        self.date_edit.setMaximumDate(QDate(2030, 12, 31))

        layout.addWidget(self.label_info)
        layout.addWidget(self.date_edit)

        self.date_edit.dateChanged.connect(self.on_date_changed)

    def on_date_changed(self, qdate: QDate):
        self.label_info.setText(f"选中日期：{qdate.toString('yyyy-MM-dd')}")


app = QApplication(sys.argv)
window = DateEditDemo()
window.show()
sys.exit(app.exec())
```

代码中的：

```python
locale_cn = QLocale(QLocale.Chinese, QLocale.China)
self.date_edit = QDateEdit()
self.date_edit.setLocale(locale_cn)
```

修改弹层显示为中文，演示效果如下：

![](/imgs/learn-qt/img-42.png)

## QTimeEdit

`QTimeEdit` 用于选择时分时间，底层依赖 `QTime` 对象，支持自定义时间格式、时间范围。

```python
from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout,
    QTimeEdit, QLabel
)
from PySide6.QtCore import QTime
import sys


class TimeEditDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QTimeEdit 演示")
        self.resize(400, 200)
        layout = QVBoxLayout(self)

        self.label_info = QLabel("选中时间：")
        self.time_edit = QTimeEdit()
        self.time_edit.setTime(QTime.currentTime())
        self.time_edit.setDisplayFormat("HH:mm:ss")

        layout.addWidget(self.label_info)
        layout.addWidget(self.time_edit)

        self.time_edit.timeChanged.connect(self.on_time_changed)

    def on_time_changed(self, qtime: QTime):
        self.label_info.setText(f"选中时间：{qtime.toString('HH:mm:ss')}")


app = QApplication(sys.argv)
window = TimeEditDemo()
window.show()
sys.exit(app.exec())
```

![](/imgs/learn-qt/img-43.png)

## 结合

这类控件常用于表单中，现将本文中的 3 个控件进行组合，实现一个简单的表单。

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout,
    QDoubleSpinBox, QDateEdit, QTimeEdit, QLabel
)
from PySide6.QtCore import QDate, QTime


class DateTimeCombineDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("浮点+日期+时间组合表单")
        self.resize(420, 360)
        layout = QVBoxLayout(self)

        self.spin_val = QDoubleSpinBox()
        self.spin_val.setRange(0, 1000)
        self.spin_val.setDecimals(3)
        self.spin_val.setSuffix(" MPa")

        self.date_edit = QDateEdit(QDate.currentDate())
        self.date_edit.setDisplayFormat("yyyy-MM-dd")
        self.date_edit.setCalendarPopup(True)

        self.time_edit = QTimeEdit(QTime.currentTime())
        self.time_edit.setDisplayFormat("HH:mm:ss")

        self.result_label = QLabel("采集信息：")

        layout.addWidget(QLabel("压力值："))
        layout.addWidget(self.spin_val)
        layout.addWidget(QLabel("采集日期："))
        layout.addWidget(self.date_edit)
        layout.addWidget(QLabel("采集时间："))
        layout.addWidget(self.time_edit)
        layout.addWidget(self.result_label)

        self.spin_val.valueChanged.connect(self.refresh)
        self.date_edit.dateChanged.connect(self.refresh)
        self.time_edit.timeChanged.connect(self.refresh)

    def refresh(self):
        val = self.spin_val.value()
        dt = self.date_edit.date().toString("yyyy-MM-dd")
        tm = self.time_edit.time().toString("HH:mm:ss")
        self.result_label.setText(f"采集信息：{val:.3f} MPa | {dt} {tm}")


app = QApplication(sys.argv)
window = DateTimeCombineDemo()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-44.png)

## 小结

`QDoubleSpinBox` 处理浮点数、`QDateEdit` 负责日期、`QTimeEdit` 管理时间，同属 Spin 体系，API 相近，是 GUI 里参数、时间录入的基础三件套，和之前的下拉框、整数微调框搭配，可搭建绝大多数表单面板。