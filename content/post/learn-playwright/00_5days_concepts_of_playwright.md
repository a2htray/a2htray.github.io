+++
date = '2026-03-20T20:04:14+08:00'
draft = false
title = '5 天学习 Playwright（Day1）：Locator、Action、Assertion、Fixture'
categories = ['后端技术', 'Playwright']
tags = ['自动化测试', '爬虫', 'Playwright']
toc = true
+++

在一次面试的过程中，问到是否接触过爬虫、是否使用过 Playwright，我其实经验不多，爬取网页数据是有做过
的，但 Playwright 没有用过。在一些岗位的任职要求中，也有描述这一工具，所以我打算趁着这段时间，学习
一个 Playwright。

学习 Playwright，对我应该也会有所益处：

1. 数据来源：我自己其实想做个数据 + AI 的一个产品，Playwright 可以帮助我更好、更快地获得数据
2. 增加技能：Playwright 应该是爬虫、自动化测试的常用工具，多项技能也好匹配更多的岗位

Playwright 的安装、基本概念，我就不赘述了，看官方的[介绍页](https://playwright.dev/python/docs/intro)就十分明了。
我主要想把其中的一些重要的概念摘出来，逐一给出描述和示例。

## Locator

**locators represent a way to find element(s) on the page at any moment.**

Locator 表示在页面中**寻找元素**的一种方式。Playwright 内置的 Locator 的方法名格式：

```bash
page.get_by_xxx
```

> page 是 Page 的实例，表示一个已打开的页面。

包括：

* page.get_by_role()：根据元素角色查找元素
* page.get_by_text()：根据文本内容查找元素
* page.get_by_label()：根据 label 值查找对应的表单元素
* page.get_by_placeholder()：根据占位符查找 input 元素
* page.get_by_alt_text()：根据 alt 属性值查找元素，通常用于查找 image 元素
* page.get_by_title()：根据 title 属性值查找元素
* page.get_by_test_id()：根据 data-testid 属性值查找元素

### page.get_by_role()

在 W3C 的规范中，每一个元素都可以计算出一个 role（角色），元素与角色的对应列表见 [Document conformance requirements for use of ARIA attributes in HTML](https://www.w3.org/TR/html-aria/#docconformance)。

如：

* 带 href 属性的 a 标签的角色是 link
* 不带 href 属性的 a 标签的角色是 generic
* article 标签的角色是 article

page.get_by_role() 的第一个参数是要寻找元素的角色，其它参数则用于过滤元素，如：

* name：`page.get_by_role("button", name="提交")` 找提交按钮
* level：`page.get_by_role("heading", level=2)` 找 h2 标签

方法声明如下：

```python
def get_by_role(
    self,
    role: Literal[
        "alert",
        "alertdialog",
        "application",
        "article",
        "banner",
        "blockquote",
        "button",
        "caption",
        "cell",
        "checkbox",
        "code",
        "columnheader",
        "combobox",
        "complementary",
        "contentinfo",
        "definition",
        "deletion",
        "dialog",
        "directory",
        "document",
        "emphasis",
        "feed",
        "figure",
        "form",
        "generic",
        "grid",
        "gridcell",
        "group",
        "heading",
        "img",
        "insertion",
        "link",
        "list",
        "listbox",
        "listitem",
        "log",
        "main",
        "marquee",
        "math",
        "menu",
        "menubar",
        "menuitem",
        "menuitemcheckbox",
        "menuitemradio",
        "meter",
        "navigation",
        "none",
        "note",
        "option",
        "paragraph",
        "presentation",
        "progressbar",
        "radio",
        "radiogroup",
        "region",
        "row",
        "rowgroup",
        "rowheader",
        "scrollbar",
        "search",
        "searchbox",
        "separator",
        "slider",
        "spinbutton",
        "status",
        "strong",
        "subscript",
        "superscript",
        "switch",
        "tab",
        "table",
        "tablist",
        "tabpanel",
        "term",
        "textbox",
        "time",
        "timer",
        "toolbar",
        "tooltip",
        "tree",
        "treegrid",
        "treeitem",
    ],
    *,
    checked: typing.Optional[bool] = None,
    disabled: typing.Optional[bool] = None,
    expanded: typing.Optional[bool] = None,
    include_hidden: typing.Optional[bool] = None,
    level: typing.Optional[int] = None,
    name: typing.Optional[typing.Union[typing.Pattern[str], str]] = None,
    pressed: typing.Optional[bool] = None,
    selected: typing.Optional[bool] = None,
    exact: typing.Optional[bool] = None,
) -> "Locator":
```

### page.get_by_text()

**根据文本内容寻找元素**

方法声明如下：

```python
def get_by_text(
    self,
    text: typing.Union[str, typing.Pattern[str]],
    *,
    exact: typing.Optional[bool] = None,
) -> "Locator":
```

比较简单，参数 `exact` 用于是否精确匹配。

### page.get_by_label()

通常，表单元素都会有一个带有文本的 label 标签，`page.get_by_label()` 就是通过 label 标签的值查找到对应的表单元素，如下方 HTML 代码：

```html
<label>Password <input type="password" /></label>
```

要找到 `<input type="password" />`，则可以使用：

```python
page.get_by_label("Password")
```

方法签名如下：

```python
def get_by_label(
    self,
    text: typing.Union[str, typing.Pattern[str]],
    *,
    exact: typing.Optional[bool] = None,
) -> "Locator":
```

参数 `exact` 用于是否精确匹配。

### page.get_by_placeholder()


如果 input 元素设置了占位符，即 placeholder 属性，则可以通过 `page.get_by_placeholder()` 方法进行查找，如：

```html
<input type="email" placeholder="name@example.com" />
```

```python
page.get_by_placeholder("name@example.com")
```

方法声明如下：

```python
def get_by_placeholder(
    self,
    text: typing.Union[str, typing.Pattern[str]],
    *,
    exact: typing.Optional[bool] = None,
) -> "Locator":
```

参数 `exact` 用于是否精确匹配。

### page.get_by_alt_text()

img 标签一般都建议带上 alt 属性，`page.get_by_alt_text()` 可以按 alt 属性值查找 img 元素，如：

```html
<img alt="playwright logo" src="/img/playwright-logo.svg" width="100" />
```

```python
page.get_by_alt_text("playwright logo")
```

方法声明如下：

```python
def get_by_alt_text(
    self,
    text: typing.Union[str, typing.Pattern[str]],
    *,
    exact: typing.Optional[bool] = None,
) -> "Locator":
```

### page.get_by_title()

title 属性是一个通用属性，即所有的 html 标签都可以带有 title 属性，Playwright 提供 `page.get_by_title()` 方法进行查找元素，如：

```html
<span title='Issues count'>25 issues</span>
```

```python
page.get_by_title("Issues count")
```

### page.get_by_test_id()

`data-testid` 属性的场景是内部测试，避免某一些元素的角色或文本发生了改变，使用方式如下：

```html
<button data-testid="directions">Itinéraire</button>
```

```python
page.get_by_test_id("directions")
```

方法签名如下：

```python
def get_by_test_id(
    self, test_id: typing.Union[str, typing.Pattern[str]]
) -> "Locator":
```

## Action

Action 指的是用户与页面之间的交互，如用户点击了一个按钮、输入框获取了焦点等。Playwright 提供了：

* 文本输入：`locator.fill()`
* 单选多选：`locator.check()`、`locator.uncheck()`
* 下拉框选择：`locator.set_option()`
* 鼠标点击：`locator.click()`、`locator.dblclick()`
* 文本按字符逐一输入：`locator.press_sequentially()`
* 按键与快捷键：`locator.press()`
* 文件上传：`locator.set_input_files()`
* 获得焦点：`locator.focus()`
* 拖拽事件：`locator.drag_to()`
  * 手动拖拽：`locator.hover()`、`mouse.down()`、`mouse.move()`、`mouse.up()`
* 滚动事件：通常不需要手动显示滚动，对于无限列表，可以使用 `locator.scroll_into_view_if_needed()`

的交互效果。

上述所有 Action 的前提，是要先查找到一个 Locator。

## Assertion

Assertion，断言，声明程序按预期的逻辑执行，Playwright 中的断言格式如下：

```python
except(locator).to_xxx()
```

| 断言方法                                    | 描述（中文）                                         | 备注                                                                   |
|-----------------------------------------|------------------------------------------------|----------------------------------------------------------------------|
| `expect(locator).to_be_checked()`       | 检查复选框或单选框是否被选中                                 | 适用于 `<input type="checkbox/radio">` 元素                               |
| `expect(locator).to_be_disabled()`      | 检查元素是否处于禁用状态                                   | 匹配 `disabled` 属性，如禁用的按钮、输入框                                          |
| `expect(locator).to_be_editable()`      | 检查元素是否可编辑（可输入内容）                               | 适用于输入框、文本域等可交互输入元素                                                   |
| `expect(locator).to_be_empty()`         | 检查元素是否为空（无子元素/无文本内容）                           | 常用于 `<div>`/`<span>`/`<ul>` 等容器元素                                    |
| `expect(locator).to_be_enabled()`       | 检查元素是否处于启用状态（非禁用）                              | `to_be_disabled()` 的反向断言                                             |
| `expect(locator).to_be_focused()`       | 检查元素是否获得焦点                                     | 适用于输入框、按钮等可获焦元素                                                      |
| `expect(locator).to_be_hidden()`        | 检查元素是否隐藏（不可见，且不占页面空间）                          | 区别于 `to_be_visible()`，匹配 `display: none` 等完全隐藏的情况                    |
| `expect(locator).to_be_visible()`       | 检查元素是否可见（在页面上显示且占空间）                           | 最常用的断言之一，会等待元素加载完成后验证                                                |
| `expect(locator).to_contain_text()`     | 检查元素包含指定文本（模糊匹配，不区分大小写，可指定 `exact=True` 精确匹配）  | 支持多语言文本，可传字符串或正则表达式                                                  |
| `expect(locator).to_have_attribute()`   | 检查元素拥有指定属性及对应值                                 | 如验证 `href`/`src`/`class` 等属性，例：`to_have_attribute("href", "/home")`  |
| `expect(locator).to_have_class()`       | 检查元素拥有指定的 CSS 类名                               | 可匹配单个类或多个类（支持部分匹配）                                                   |
| `expect(locator).to_have_count()`       | 检查匹配的元素数量等于指定值                                 | 适用于列表、表格行等批量元素计数，例：`expect(page.get_by_role("li")).to_have_count(5)` |
| `expect(locator).to_have_css()`         | 检查元素拥有指定的 CSS 样式属性及值                           | 例：`to_have_css("color", "rgb(255, 0, 0)")`（验证文字红色）                   |
| `expect(locator).to_have_id()`          | 检查元素的 `id` 属性等于指定值                             | 精准匹配元素 ID，例：`to_have_id("submit-btn")`                               |
| `expect(locator).to_have_js_property()` | 检查元素的 JavaScript 属性等于指定值                       | 访问元素的 JS 内置属性，例：`to_have_js_property("value", "123")`（输入框值）          |
| `expect(locator).to_have_text()`        | 检查元素的文本内容完全匹配指定值（默认精确匹配，可设 `exact=False` 模糊匹配） | 与 `to_contain_text()` 对比：前者是完全匹配，后者是包含匹配                             |
| `expect(locator).to_have_value()`       | 检查输入框/选择框的 `value` 属性等于指定值                     | 适用于 `<input>`/`<select>` 等表单元素                                       |
| `expect(page).to_have_title()`          | 检查页面标题（`<title>` 标签）匹配指定值                      | 支持模糊/精确匹配，例：`to_have_title("首页", exact=False)`                       |
| `expect(page).to_have_url()`            | 检查页面 URL 匹配指定值                                 | 支持字符串或正则表达式，例：`to_have_url(re.compile(r"\/user\/\d+"))`              |
| `expect(locator).not_to_*`              | 反向断言：检查上述所有条件不成立                               | 例：`expect(locator).not_to_be_visible()`（元素不可见）                       |

详情见文档：[Assertions](https://playwright.dev/python/docs/test-assertions)。

## Fixture

Fixture 的本质是提供一套可复用的前置或后置的逻辑块，通过装饰器函数 `@pytest.fixture()` 将一个函数标记为 Fixture：

* `yield` 关键字用于区分前置逻辑和后置逻辑
  * `yield` 之前，测试执行前的准备工作，如启动浏览器、打开页面
  * `yield` 之后，测试执行后的准备工作，如关闭页面、清理缓存
* `scope` 参数：控制 Fixture 的作用域
  * function：每一个测试用例执行一次，适用于页面、上下文
  * class：每一个测试类执行一次，适用于共享数据的类级测试
  * module：每一个模块执行一次，适用于测试器实例
  * session：整个测试会话执行一次，适用于全局资源

例如：

```python
import pytest
from playwright.sync_api import Playwright, Browser, Page

# ========== 1. 自定义 Fixture：管理浏览器启动/关闭 ==========
@pytest.fixture(scope="module")  # scope：作用域（module=整个模块仅执行一次）
def browser(playwright: Playwright) -> Browser:
    """自定义浏览器Fixture：启动Chromium，测试完成后关闭"""
    # 测试前置：启动浏览器（无头模式）
    browser = playwright.chromium.launch(headless=True)
    yield browser  # 关键：yield 前是前置逻辑，后是后置逻辑
    # 测试后置：关闭浏览器
    browser.close()

# ========== 2. 自定义 Fixture：基于浏览器创建页面 ==========
@pytest.fixture(scope="function")  # scope：function=每个测试用例执行一次
def page(browser: Browser) -> Page:
    """自定义页面Fixture：每个测试用例创建新页面，自动登录"""
    # 测试前置：创建上下文 + 页面，模拟登录
    context = browser.new_context()
    page = context.new_page()
    page.goto("https://example.com")  # 打开测试页面
    yield page  # 将page实例传递给测试用例
    # 测试后置：关闭页面和上下文
    page.close()
    context.close()
```