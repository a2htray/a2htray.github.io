+++
date = '2026-03-23T10:28:24+08:00'
draft = false
title = '5 天学习 Playwright（Day2）：API 测试、API 模拟、认证、下载'
categories = ['后端技术', 'Playwright']
tags = ['自动化测试', '爬虫', 'Playwright']
toc = true
+++

今天是学习 Playwright 的第 2 天，文章的主要内容来自于官方的文档：

1. 学习官方的表述、加深自己的理解
2. 试着写一些自己的示例，更贴近实际

## API 测试

`Playwright can be used to get access to the REST API of your application.`

文档的第一句就很直接明了，即 Playwright 可以访问 REST API 服务。

实现上述功能的基础是 `APIRequestContext` 类，见 [APIRequestContext](https://playwright.dev/python/docs/api/class-apirequestcontext)。
每一次浏览器上下文都会关联到一个 `APIRequestContext` 实例，该实例与浏览器上下文共享 cookie 及存储信息，所以

**用 APIRequestContext 发送请求时，会自动带上浏览器上下文里的登录状态、Cookie 等，无需手动设置。**

在单个 python 文件中实现：

1. 一个地址 `/` 返回 HTML 页面，其中设置 cookie 值，value=testing
2. 一个接口 `/api/cookie`，并在接口中向终端写 cookie 值
3. 可单文件启动 web 服务

完整代码如下：

```python
from http.server import HTTPServer, BaseHTTPRequestHandler
import json


class RequestHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/":
            # 首页 - 设置 cookie
            self.send_response(200)
            self.send_header("Content-Type", "text/html; charset=utf-8")
            self.send_header("Set-Cookie", "value=testing")
            self.end_headers()
            
            html = """
            <!DOCTYPE html>
            <html>
            <head>
                <meta charset="UTF-8">
                <title>Cookie Test</title>
            </head>
            <body>
                <h1>Cookie 已设置</h1>
                <p>Cookie "value" 已设置为 "testing"</p>
                <p><a href="/api/cookie">查看 Cookie</a></p>
            </body>
            </html>
            """
            self.wfile.write(html.encode("utf-8"))
            
        elif self.path == "/api/cookie":
            # API 接口 - 打印 cookie
            cookie_header = self.headers.get("Cookie", "")
            
            # 解析 cookie
            cookies = {}
            if cookie_header:
                for item in cookie_header.split(";"):
                    if "=" in item:
                        key, value = item.strip().split("=", 1)
                        cookies[key] = value
            
            # 打印到控制台
            print(f"Received cookies: {cookies}")
            print(f"Cookie 'value': {cookies.get('value', 'not found')}")
            
            # 返回 JSON 响应
            self.send_response(200)
            self.send_header("Content-Type", "application/json")
            self.end_headers()
            
            response = {
                "cookies": cookies,
                "value": cookies.get("value", None)
            }
            self.wfile.write(json.dumps(response, ensure_ascii=False).encode("utf-8"))
            
        else:
            self.send_response(404)
            self.end_headers()
            self.wfile.write(b"Not Found")

    def log_message(self, format, *args):
        # 简化日志输出
        print(f"[{self.log_date_time_string()}] {format % args}")


def run_server(port=9090):
    server = HTTPServer(("", port), RequestHandler)
    print(f"Server running at http://localhost:{port}/")
    print(f"Press Ctrl+C to stop")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopped")


if __name__ == "__main__":
    run_server()
```

运行服务：

```bash
$ python web.py
Server running at http://localhost:9090/
Press Ctrl+C to stop
```

编写测试用例 cookie.py

```python
from playwright.sync_api import Page

def test_cookie(page: Page):
    page.goto("http://localhost:9090/")
    assert page.title() == "Cookie Test"
   
    # 发送请求 /api/cookie
    response = page.request.get("http://localhost:9090/api/cookie")
    assert response.status == 200
```

执行 `pytest cookie.py`，在 web 服务终端显示：

```bash
[24/Mar/2026 14:05:05] "GET / HTTP/1.1" 200 -
Received cookies: {'value': 'testing'}
Cookie 'value': testing
[24/Mar/2026 14:05:06] "GET /api/cookie HTTP/1.1" 200 -
```

可见，在调用 `page.request` 的方法时，之前请求的 cookie 值已自动填充到新的请求中。

## API 模拟

Playwright 提供模拟和修改网络流量的 API，格式如：

```python
page.route(url_pattern, handler)
```

* url_pattern 要模拟或修改的地址，通过通配符，如 `**/*`
* handler 为一个要执行的函数，参数为路由 Route 实例

```python
# 拦截路由
page.route(url_pattern, handler)

# 处理函数里：
async def handler(route):
    # 1. 获取真实请求
    request = route.request

    # 2. 继续请求（不改）
    await route.continue_()

    # 3. 拦截并返回假数据
    await route.fulfill(...)

    # 4. 先获取真实返回，再修改
    response = await route.fetch()
    await route.fulfill(body=新内容, ...)
```

示例如：

```python
def test_api_intercept(page: Page):
    def handle(route: Route):
        route.fulfill(json={"value": "testing-intercepted"})
   
    # 拦截 /api/cookie 请求
    page.route("http://localhost:9090/api/cookie", handle)
    page.goto("http://localhost:9090")
```

`page.goto` 的页面中，请求 `/api/cookie` 的返回会是 `{"value": "testing-intercepted"}`。

## 认证

Playwright 认证是指将登录后的状态（cookie + localStorage）保存，并在后续的测试中复用，不用每个用例都重新登录。

比如定义一个 session 级别的 Fixture，实现登录逻辑，再定义一个 function 级别的 Fixture，加载登录状态。

```python
# conftest.py （pytest 全局配置文件，固定名字）
import pytest
from playwright.sync_api import Page

# 全局只执行一次：登录并保存状态
@pytest.fixture(scope="session", autouse=True)
def login_and_save_state(page: Page):
    # 1. 登录
    page.goto("https://login-page.com")
    page.get_by_label("username").fill("your-account")
    page.get_by_label("password").fill("your-password")
    page.get_by_role("button", name="Login").click()
    
    # 等待登录成功
    page.wait_for_url("**/dashboard**")

    # 2. 保存登录状态（Cookie + localStorage 全存）
    page.context.storage_state(path="auth.json")

# 每个测试自动加载登录态
@pytest.fixture(scope="function")
def authenticated_page(page: Page):
    page.context.storage_state(path="auth.json")
    return page
```

## 下载

Playwright 处理文件下载的核心在于：

1. 自动捕获下载
2. 管理临时文件
3. 持久化保存

**自动捕获下载**

使用两个关键 API：

1. `page.expect_download()` 等待并捕获下载
2. `page.on("download")` 监听下载事件

**管理临时文件**

在默认设置中，开启下载允许（accept_downloads=True），文件存储至临时目录，当上下文关闭时，临时文件被删除。

**持久化保存**

使用 `save_as()` 将文件保存到指定路径。
