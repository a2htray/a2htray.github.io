+++
date = '2026-11-11T09:17:17+08:00'
draft = false
title = 'PySide6 滑块与进度：QSlider、QProgressBar、QDial'
categories = ['后端技术', 'PySide6']
tags = ['Python', 'PySide6', 'GUI', '滑块', 'QSlider', 'QProgressBar', 'QDial']
toc = true
+++

![](/imgs/learn-qt/cover12.png)

本文介绍的 3 个控件处理的是连续量，用拖拽的方式在一段区间里取一个值，或者把一件事情的完成程度可视化表示出来，包括：

- `QSlider`：滑块
- `QProgressBar`：进度条
- `QDial`：仪表盘

程序界面中的音量、亮度、缩放、导出进度等，都是用这 3 个控件实现。

## QSlider

`QSlider` 是最常用的**连续取值**控件，横向用于音量、缩放，纵向用于刻度尺式输入。

> 构造函数中传入 orientation 参数，指定滑块的方向（`Qt.Orientation.Horizontal` 或 `Qt.Orientation.Vertical`）。

`QSlider` 的常用方法包括：

- `setMinimum`：设置最小值
- `setMaximum`：设置最大值
- `setRange`：一次性设置最小、最大值
- `setValue`：设置滑块当前值
- `value`：获取当前值
- `singleStep`：设置方向键单步增量
- `pageStep`：设置 PageUp/PageDown 大步长
- `setTickPosition`：设置刻度位置
- `setTickInterval`：设置刻度间隔
- `invertedAppearance`：是否反向显示（水平条从右到左）
- `invertedControls`：是否反转键盘控制方向

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QHBoxLayout,
    QSlider, QLabel, QPushButton, QSpinBox
)
from PySide6.QtCore import Qt


class SliderDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QSlider 常用方法演示")
        self.resize(600, 400)

        layout = QVBoxLayout()

        # 滑块本体
        self.slider = QSlider(Qt.Horizontal)
        # setRange(min, max) 一次性设置范围
        self.slider.setRange(0, 100)
        # setValue 设置初始值
        self.slider.setValue(50)
        # setSingleStep 方向键单步
        self.slider.setSingleStep(2)
        # setPageStep PageUp/PageDown步长
        self.slider.setPageStep(10)
        # 刻度
        self.slider.setTickPosition(QSlider.TicksBelow)
        self.slider.setTickInterval(10)

        # 显示当前value
        self.label_val = QLabel(f"当前值：{self.slider.value()}")

        layout.addWidget(self.label_val)
        layout.addWidget(self.slider)

        # 按钮组
        h_layout = QHBoxLayout()

        btn_min = QPushButton("setMinimum(10)")
        btn_min.clicked.connect(lambda: self.slider.setMinimum(10))

        btn_max = QPushButton("setMaximum(90)")
        btn_max.clicked.connect(lambda: self.slider.setMaximum(90))

        btn_reset_range = QPushButton("setRange(0,100)")
        btn_reset_range.clicked.connect(lambda: self.slider.setRange(0, 100))

        btn_set_val = QPushButton("setValue(50)")
        btn_set_val.clicked.connect(lambda: self.slider.setValue(50))

        h_layout.addWidget(btn_min)
        h_layout.addWidget(btn_max)
        h_layout.addWidget(btn_reset_range)
        h_layout.addWidget(btn_set_val)
        layout.addLayout(h_layout)

        h_layout2 = QHBoxLayout()
        btn_tick_none = QPushButton("无刻度")
        btn_tick_none.clicked.connect(lambda: self.slider.setTickPosition(QSlider.NoTicks))

        btn_tick_below = QPushButton("刻度在下")
        btn_tick_below.clicked.connect(lambda: self.slider.setTickPosition(QSlider.TicksBelow))

        btn_tick_above = QPushButton("刻度在上")
        btn_tick_above.clicked.connect(lambda: self.slider.setTickPosition(QSlider.TicksAbove))

        h_layout2.addWidget(btn_tick_none)
        h_layout2.addWidget(btn_tick_below)
        h_layout2.addWidget(btn_tick_above)
        layout.addLayout(h_layout2)

        h_layout3 = QHBoxLayout()
        btn_inv_app = QPushButton("切换 invertedAppearance")
        btn_inv_app.clicked.connect(self.toggle_inv_appearance)

        btn_inv_ctrl = QPushButton("切换 invertedControls")
        btn_inv_ctrl.clicked.connect(self.toggle_inv_controls)

        h_layout3.addWidget(btn_inv_app)
        h_layout3.addWidget(btn_inv_ctrl)
        layout.addLayout(h_layout3)

        # 单步、大步修改
        h_layout4 = QHBoxLayout()
        sp_single = QSpinBox()
        sp_single.setRange(1, 20)
        sp_single.setValue(2)
        sp_single.valueChanged.connect(self.slider.setSingleStep)
        h_layout4.addWidget(QLabel("singleStep:"))
        h_layout4.addWidget(sp_single)

        sp_page = QSpinBox()
        sp_page.setRange(1, 30)
        sp_page.setValue(10)
        sp_page.valueChanged.connect(self.slider.setPageStep)
        h_layout4.addWidget(QLabel("pageStep:"))
        h_layout4.addWidget(sp_page)
        layout.addLayout(h_layout4)

        self.setLayout(layout)

        # 值变化信号
        self.slider.valueChanged.connect(self.on_value_changed)

    def on_value_changed(self, v):
        self.label_val.setText(f"当前值：{v}")

    def toggle_inv_appearance(self):
        flag = not self.slider.invertedAppearance()
        self.slider.setInvertedAppearance(flag)

    def toggle_inv_controls(self):
        flag = not self.slider.invertedControls()
        self.slider.setInvertedControls(flag)


app = QApplication(sys.argv)
window = SliderDemo()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-45.png)

## QProgressBar

`QProgressBar` 只负责显示进度，不接收用户输入。

`QProgressBar` 的常用方法包括：

- `setMinimum`：设置最小值
- `setMaximum`：设置最大值
- `setRange`：一次性设置最小、最大值
- `setValue`：设置当前进度值
- `value`：获取当前进度值
- `setFormat`：自定义文字模板
  - `%p` = 百分比，`%v` = 当前值，`%m` = 最大值
- `resetFormat`：恢复默认格式（`%p%`）
- `setTextVisible`：是否显示进度文字
- `setAlignment`：文字对齐方式
- `setInvertedAppearance`：反向填充（从右往左涨）
- `setOrientation`：水平 / 垂直方向
- `reset`：重置到最小值
- `text`：获取当前显示文本

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QHBoxLayout,
    QProgressBar, QLabel, QPushButton, QSpinBox
)
from PySide6.QtCore import Qt


class ProgressBarDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QProgressBar 常用方法演示")
        self.resize(650, 450)

        layout = QVBoxLayout()

        # 进度条本体
        self.bar = QProgressBar()
        self.bar.setRange(0, 100)
        self.bar.setValue(30)
        self.bar.setAlignment(Qt.AlignCenter)

        layout.addWidget(QLabel("进度条："))
        layout.addWidget(self.bar)

        # 状态文本
        self.label_info = QLabel(f"value={self.bar.value()}, text={self.bar.text()}")
        layout.addWidget(self.label_info)

        # 范围、setValue
        h1 = QHBoxLayout()
        btn_set_min = QPushButton("setMinimum(10)")
        btn_set_min.clicked.connect(lambda: self.bar.setMinimum(10))

        btn_set_max = QPushButton("setMaximum(90)")
        btn_set_max.clicked.connect(lambda: self.bar.setMaximum(90))

        btn_reset_range = QPushButton("setRange(0,100)")
        btn_reset_range.clicked.connect(lambda: self.bar.setRange(0, 100))

        btn_set_val = QPushButton("setValue(60)")
        btn_set_val.clicked.connect(lambda: self.bar.setValue(60))

        h1.addWidget(btn_set_min)
        h1.addWidget(btn_set_max)
        h1.addWidget(btn_reset_range)
        h1.addWidget(btn_set_val)
        layout.addLayout(h1)

        # 不确定模式、reset
        h2 = QHBoxLayout()
        btn_unknown = QPushButton("不确定模式 setRange(0,0)")
        btn_unknown.clicked.connect(lambda: self.bar.setRange(0, 0))

        btn_reset = QPushButton("reset()")
        btn_reset.clicked.connect(self.bar.reset)

        h2.addWidget(btn_unknown)
        h2.addWidget(btn_reset)
        layout.addLayout(h2)

        # 文字控制：显示隐藏、格式
        h3 = QHBoxLayout()
        btn_text_hide = QPushButton("隐藏文字 setTextVisible(False)")
        btn_text_hide.clicked.connect(lambda: self.bar.setTextVisible(False))

        btn_text_show = QPushButton("显示文字 setTextVisible(True)")
        btn_text_show.clicked.connect(lambda: self.bar.setTextVisible(True))

        btn_fmt1 = QPushButton('格式 "%v/%m"')
        btn_fmt1.clicked.connect(lambda: self.bar.setFormat("%v/%m"))

        btn_fmt2 = QPushButton('格式 "%p% 已完成"')
        btn_fmt2.clicked.connect(lambda: self.bar.setFormat("%p% 已完成"))

        btn_fmt_reset = QPushButton("resetFormat()")
        btn_fmt_reset.clicked.connect(self.bar.resetFormat)

        h3.addWidget(btn_text_hide)
        h3.addWidget(btn_text_show)
        h3.addWidget(btn_fmt1)
        h3.addWidget(btn_fmt2)
        h3.addWidget(btn_fmt_reset)
        layout.addLayout(h3)

        # 反向、方向
        h4 = QHBoxLayout()
        btn_inv = QPushButton("切换 invertedAppearance")
        btn_inv.clicked.connect(self.toggle_inv)

        btn_hori = QPushButton("水平")
        btn_hori.clicked.connect(lambda: self.bar.setOrientation(Qt.Horizontal))

        btn_vert = QPushButton("垂直")
        btn_vert.clicked.connect(lambda: self.bar.setOrientation(Qt.Vertical))

        h4.addWidget(btn_inv)
        h4.addWidget(btn_hori)
        h4.addWidget(btn_vert)
        layout.addLayout(h4)

        # 手动输入值
        h5 = QHBoxLayout()
        sp_val = QSpinBox()
        sp_val.setRange(0, 200)
        sp_val.setValue(30)
        sp_val.valueChanged.connect(self.bar.setValue)
        h5.addWidget(QLabel("修改value："))
        h5.addWidget(sp_val)
        layout.addLayout(h5)

        self.setLayout(layout)
        self.bar.valueChanged.connect(self.on_change)

    def on_change(self, v):
        self.label_info.setText(f"value={v}, text={self.bar.text()}")

    def toggle_inv(self):
        flag = not self.bar.invertedAppearance()
        self.bar.setInvertedAppearance(flag)


app = QApplication(sys.argv)
window = ProgressBarDemo()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-46.png)

## QDial

`QDial` 用于选择角度，内置角度选择器，支持自定义角度显示格式、角度范围，底层依赖 `QAngle` 对象。

`QDial` 的常用方法包括：

- `setMinimum`：设置最小值
- `setMaximum`：设置最大值
- `setRange`：一次性设置最小、最大值
- `setValue`：设置当前值
- `value`：获取当前值
- `setSingleStep`：单步增量（键盘 / 方向键）
- `setPageStep`：PageUp/PageDown 大步长
- `setWrapping`：是否循环环绕（转到尽头可以继续转，0 和最大值相连）
- `setNotchesVisible`：是否显示刻度凹槽
- `setNotchTarget`：刻度目标间距（控制刻度疏密）
- `invertedAppearance`：反向，默认顺时针增大；开启后逆时针增大

```python
import sys

from PySide6.QtWidgets import (
    QApplication, QWidget, QVBoxLayout, QHBoxLayout,
    QDial, QLabel, QPushButton, QSpinBox
)
from PySide6.QtCore import Qt


class QDialDemo(QWidget):
    def __init__(self):
        super().__init__()
        self.setWindowTitle("QDial 常用方法演示")
        self.resize(600, 450)

        layout = QVBoxLayout()

        # 旋钮本体
        self.dial = QDial()
        self.dial.setRange(0, 100)
        self.dial.setValue(50)
        self.dial.setSingleStep(1)
        self.dial.setPageStep(10)
        self.dial.setNotchesVisible(True)
        self.dial.setNotchTarget(3.0)

        layout.addWidget(QLabel("QDial 旋钮："))
        layout.addWidget(self.dial, alignment=Qt.AlignmentFlag.AlignCenter)

        self.label_val = QLabel(f"当前值：{self.dial.value()}")
        layout.addWidget(self.label_val)

        # 范围、setValue
        h1 = QHBoxLayout()
        btn_min = QPushButton("setMinimum(10)")
        btn_min.clicked.connect(lambda: self.dial.setMinimum(10))

        btn_max = QPushButton("setMaximum(90)")
        btn_max.clicked.connect(lambda: self.dial.setMaximum(90))

        btn_range = QPushButton("setRange(0,100)")
        btn_range.clicked.connect(lambda: self.dial.setRange(0, 100))

        btn_setval = QPushButton("setValue(50)")
        btn_setval.clicked.connect(lambda: self.dial.setValue(50))

        h1.addWidget(btn_min)
        h1.addWidget(btn_max)
        h1.addWidget(btn_range)
        h1.addWidget(btn_setval)
        layout.addLayout(h1)

        # 环绕 wrapping
        h2 = QHBoxLayout()
        btn_wrap_on = QPushButton("开启环绕 setWrapping(True)")
        btn_wrap_on.clicked.connect(lambda: self.dial.setWrapping(True))

        btn_wrap_off = QPushButton("关闭环绕 setWrapping(False)")
        btn_wrap_off.clicked.connect(lambda: self.dial.setWrapping(False))
        h2.addWidget(btn_wrap_on)
        h2.addWidget(btn_wrap_off)
        layout.addLayout(h2)

        # 刻度显示、刻度间距
        h3 = QHBoxLayout()
        btn_notch_show = QPushButton("显示刻度 setNotchesVisible(True)")
        btn_notch_show.clicked.connect(lambda: self.dial.setNotchesVisible(True))

        btn_notch_hide = QPushButton("隐藏刻度 setNotchesVisible(False)")
        btn_notch_hide.clicked.connect(lambda: self.dial.setNotchesVisible(False))
        h3.addWidget(btn_notch_show)
        h3.addWidget(btn_notch_hide)
        layout.addLayout(h3)

        # 反向
        h4 = QHBoxLayout()
        btn_inv = QPushButton("切换 invertedAppearance")
        btn_inv.clicked.connect(self.toggle_inv)
        h4.addWidget(btn_inv)
        layout.addLayout(h4)

        # 单步、大步、手动赋值
        h5 = QHBoxLayout()
        sp_single = QSpinBox()
        sp_single.setRange(1, 20)
        sp_single.setValue(1)
        sp_single.valueChanged.connect(self.dial.setSingleStep)
        h5.addWidget(QLabel("singleStep:"))
        h5.addWidget(sp_single)

        sp_page = QSpinBox()
        sp_page.setRange(1, 30)
        sp_page.setValue(10)
        sp_page.valueChanged.connect(self.dial.setPageStep)
        h5.addWidget(QLabel("pageStep:"))
        h5.addWidget(sp_page)

        sp_target = QSpinBox()
        sp_target.setRange(1, 10)
        sp_target.setValue(3)
        sp_target.valueChanged.connect(lambda v: self.dial.setNotchTarget(float(v)))
        h5.addWidget(QLabel("notchTarget:"))
        h5.addWidget(sp_target)
        layout.addLayout(h5)

        self.setLayout(layout)
        self.dial.valueChanged.connect(self.on_changed)

    def on_changed(self, v):
        self.label_val.setText(f"当前值：{v}")

    def toggle_inv(self):
        flag = not self.dial.invertedAppearance()
        self.dial.setInvertedAppearance(flag)


app = QApplication(sys.argv)
window = QDialDemo()
window.show()
sys.exit(app.exec())
```

演示效果如下：

![](/imgs/learn-qt/img-47.png)

## 小结

`QSlider`、`QProgressBar`、`QDial` 都属于数值可视化控件，但定位和交互场景差别明显：

- **QSlider（滑动条）**：线性拖拽交互，适合在一段区间内精细调节参数，支持水平 / 垂直、刻度显示，是最通用的调节控件；
- **QProgressBar（进度条）**：侧重状态展示，多用于任务进度加载，支持百分比自定义文本、不确定加载动画，一般只用来展示，不用于用户调节；
- **QDial（旋钮）**：环形旋转交互，自带环绕模式，适合音量、角度等环形参数调节，交互氛围感更强。
