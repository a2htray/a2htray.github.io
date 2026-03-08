+++
date = '2026-03-04T10:29:27+08:00'
draft = false
title = 'JavaScript 中的 async/await'
categories = ['前端技术', 'JavaScript']
tags = ['JavaScript', 'async/await']
toc = true
+++

`async/await` 是基于 Promise、以同步方式编写异步代码的语法糖。

简单来说：

async 放置在函数定义前，函数总是返回一个 Promise，返回值会自动包裹在已解析的 Promise 中。

await 让 JavaScript 等待 Promise 处理完成并返回结果。

## async 函数

async 关键字用于声明异步函数。

异步函数声明会创建一个 AsyncFunction 对象。每次调用异步函数时，它都会返回一个新的 Promise —— 该 Promise 会在异步函数返回值时兑现（resolve），或在异步函数内部抛出未捕获异常时拒绝（reject）。

异步函数可包含零个或多个 await 表达式。await 表达式会暂停函数执行，直至其等待的 Promise 兑现（fulfilled）或拒绝（rejected），从而让返回 Promise 的函数表现得如同同步函数一般。该 Promise 的兑现值会被视作 await 表达式的返回值。使用 async 和 await 能够在异步代码中使用常规的 try/catch 代码块（捕获异常）。

## await 语法

使用格式：

```javascript
await expression
```

参数：

* expression：Promise 对旬，或 thenable 对象，或其它须等待的值

返回值：

Promise 或 thenable 对象的兑现值，若该表达式并非 thenable 类型，则为表达式自身的值。

## 记住 3 点

1. async 函数总是返回一个 Promise
2. await 只能在 async 函数中使用
3. await 等待一个 Promise 完成

## 使用 async/await 的好处

1. 提高代码的可读性
2. 易于调试、编写测试用例
3. 便于错误处理
4. 减少内嵌

## 资源

1. [面试官：说说async/await？我差点翻车！原来还可以这么用 ](https://juejin.cn/post/7563547711006408738)
2. [JavaScript.info Async/await](https://javascript.info/async-await)
3. [When to Use Async/Await vs Promises in JavaScript](https://www.freecodecamp.org/news/when-to-use-asyncawait-vs-promises-in-javascript/)
4. [MDN async](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Statements/async_function)
5. [MDN await](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/await)
