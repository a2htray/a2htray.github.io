+++
date = '2026-03-27T14:17:22+08:00'
draft = false
title = '5 天学习 Playwright（Day5）：改进文献信息爬取'
categories = ['后端技术', 'Playwright']
tags = ['爬虫', 'Playwright', '文献爬取']
toc = true
+++

第 4 天的实践中，留了不少的问题，还有很多提升的空间，如：

1. 启用无头模式，避免开关浏览器的开销
2. 启动用多线程获取文献详情

## 启用无头模式，避免开关浏览器的开销

针对这一改进，可通过以下设置完成：

1. 禁止无头模式检测
2. 显示设置视品大小
3. 设置 user-agent，避免被识别为爬虫

代码如下：

```python
async with async_playwright() as p:
    browser: Browser = await p.chromium.launch(
        headless=True,
        # 1. 禁用无头模式检测
        args=['--disable-blink-features=AutomationControlled'],
    )
    context: BrowserContext = await browser.new_context(
        # 2. 显示设置视口大小
        viewport={'width': 1920, 'height': 1080},
        # 3. 伪造 user-agent，避免被检测到是爬虫
        user_agent='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36',
    )
    # 其它
```

## 启动用多线程获取文献详情

多线程的改进主要体现在获取文献详情上，当然也可以多线程获取文献列表，待读者自行实现。

Playwright 异步 API 天然支持并发，所以我们使用 `asyncio.gather` 可以轻松实现，如下：

```python
for i in range(0, len(publications), batch_size):
    await asyncio.gather(*[fullfill_publication(publication, await context.new_page()) for publication in publications[i:i + batch_size]])
```

* batch_size 设置单个线程处理的文献数

## 执行效果

```bash
$ python get_publications_by_pubmed_improved.py
开始获取文献基本信息（标题、详情地址）
获取第1页文献基本信息（标题、详情地址）
获取第2页文献基本信息（标题、详情地址）
获取第3页文献基本信息（标题、详情地址）
获取第4页文献基本信息（标题、详情地址）
获取第5页文献基本信息（标题、详情地址）
获取第6页文献基本信息（标题、详情地址）
获取第7页文献基本信息（标题、详情地址）
获取第8页文献基本信息（标题、详情地址）
获取第9页文献基本信息（标题、详情地址）
获取第10页文献基本信息（标题、详情地址）
获取文献基本信息（标题、详情地址）完成
开始获取文献详细信息（作者、期刊、摘要、PMID、DOI、出版日期）
获取文献详细信息（作者、期刊、摘要、PMID、DOI、出版日期）完成
耗时：179.30069017410278 秒
```

较于上一个版本，时间加快了 720s，即 12 分钟，效果很明显。

## 遇到的问题

> Future exception was never retrieved
future: <Future finished exception=TargetClosedError('Target page, context or browser has been closed')>
playwright._impl._errors.TargetClosedError: Target page, context or browser has been closed

出现上述问题的原因于：

1. 浏览器打开过多页面导致浏览器崩溃自动自闭
2. 当然多个线程共用了一个 Page 实例

改进：

1. 合理设置 batch_size 值
2. fullfill_publication 传入新的 Page 实例

## 完整代码

```python
import asyncio
import typing
import csv
import time

from pydantic import BaseModel
from playwright.async_api import async_playwright, Browser, Page, BrowserContext


class Publication(BaseModel):
    title: str
    authors: typing.Optional[str] = None
    journal: typing.Optional[str] = None
    first_affiliation: typing.Optional[str] = None
    abstract: typing.Optional[str] = None
    pmid: typing.Optional[str] = None
    doi: typing.Optional[str] = None
    publish_date: typing.Optional[str] = None
    pubmed_url: str


def save_publications(publications: typing.List[Publication], file: str):
    with open(file, 'w', newline='', encoding='utf-8') as f:
        writer = csv.writer(f)
        writer.writerow(publications[0].model_dump().keys())
        for publication in publications:
            writer.writerow(publication.model_dump().values())


async def fullfill_publication(publication: Publication, page: Page):
    await page.goto(publication.pubmed_url)
    await page.wait_for_load_state('networkidle')

    # 取 authors
    author_eles = await page.query_selector_all('span.authors-list-item a.full-name')
    publication.authors = ', '.join([await author_ele.text_content() for author_ele in author_eles])

    # 取 journal
    journal_ele = await page.wait_for_selector('#full-view-journal-trigger', timeout=15000)
    publication.journal = (await journal_ele.text_content()).strip()

    # 取 first_affiliation
    await page.locator('#toggle-authors').click()
    first_affiliation_ele = page.locator('#full-view-affiliation-1')
    # [2:] 去掉前两个字符，因为前两个字符是 "1 "
    publication.first_affiliation = (await first_affiliation_ele.text_content()).strip()[2:]

    # 取 abstract
    abstract_ele = await page.query_selector('#eng-abstract')
    publication.abstract = (await abstract_ele.text_content()).strip()

    # 取 pmid
    pmid_ele = await page.query_selector('#full-view-identifiers strong.current-id')
    publication.pmid = (await pmid_ele.text_content()).strip()

    # 取 doi
    doi_ele = await page.query_selector('a[data-ga-action="DOI"]')
    publication.doi = (await doi_ele.text_content()).strip()

    # 取 publish_date
    publish_date_ele = await page.query_selector('span.cit')
    publication.publish_date = (await publish_date_ele.text_content()).strip().split(';')[0]


async def run(query: str, *, size: int, save_to: str, batch_size: int = 10):
    base_url = 'https://pubmed.ncbi.nlm.nih.gov/'
    async with async_playwright() as p:
        browser: Browser = await p.chromium.launch(
            headless=True,
            # 1. 禁用无头模式检测
            args=['--disable-blink-features=AutomationControlled'],
        )
        context: BrowserContext = await browser.new_context(
            # 2. 显示设置视口大小
            viewport={'width': 1920, 'height': 1080},
            # 3. 伪造 user-agent，避免被检测到是爬虫
            user_agent='Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36',
        )
        page: Page = await context.new_page()

        print('开始获取文献基本信息（标题、详情地址）')
        publications = []
        for page_num in range(1, size // 10 + 1):
            print(f'获取第{page_num}页文献基本信息（标题、详情地址）')
            # term 查询关键词
            # sort 排序方式，pubdate 表示按出版日期排序
            # page 页码
            # filter=simsearch1.fha 表示文献要摘要可访问
            # filter=datesearch.y_5 表示搜索范围为近 5 年文献
            await page.goto(f'{base_url}?term={query}&sort=pubdate&page={page_num}&filter=simsearch1.fha&filter=datesearch.y_5')

            title_eles = await page.query_selector_all('article a.docsum-title')
            for title_ele in title_eles:
                title = await title_ele.text_content()
                href = await title_ele.get_attribute('href')
                publications.append(Publication(
                    title=title.strip(),
                    # [1:] 去掉第一个字符，因为第一个字符是 /
                    pubmed_url=f'{base_url}{href.strip()[1:]}',
                ))

        print('获取文献基本信息（标题、详情地址）完成')

        # 采用多线程并发，获取文献详细信息
        print('开始获取文献详细信息（作者、期刊、摘要、PMID、DOI、出版日期）')

        for i in range(0, len(publications), batch_size):
            await asyncio.gather(*[fullfill_publication(publication, await context.new_page()) for publication in publications[i:i + batch_size]])

        print('获取文献详细信息（作者、期刊、摘要、PMID、DOI、出版日期）完成')

        save_publications(publications, save_to)

        await browser.close()

if __name__ == '__main__':
    start_time = time.time()
    size = 100
    batch_size = 10
    save_to = 'publications_improved.csv'

    asyncio.run(run(
        'maize',
        size=size,
        save_to=save_to,
        batch_size=batch_size,
    ))

    end_time = time.time()
    print(f'耗时：{end_time - start_time} 秒')
```