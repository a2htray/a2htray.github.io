---
title: "Go 的错误处理"
date: 2022-08-22T23:40:01+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - 翻译
---

原文：[Error handling and Go](https://go.dev/blog/error-handling-and-go)



# 介绍

如果你写过 Go 的代码，就一定遇到过 Go 的内置类型 error。一个 error 类型的值可用于指明程序的某种不正常的状态，比如，当打开文件失败时，os.Open 函数会返回一个非 nil 的 error 值。

<!--more-->

```go
func Open(name string) (file *File, err error)
```

下面代码演示了：使用 os.Open 打开一个文件失败时，用 log.Fatal 来打印错误信息和停止程序运行。

```go
f, err := os.Open("filename.ext")
if err != nil {
    log.Fatal(err)
}
// do something with the open *File f
```

只要知道上述一点关于 error 的内容，在 Go 中就可以做很多事，但在这篇文章中，我们会进一步地讨论 error 类型以及 Go 处理错误的最佳实践。

# error 类型

error 类型是一种接口类型，error 值可以是任意被字符串所能表示的值，下方是接口的定义：

```go
type error interface {
    Error() string
}
```

与其它内置类型一样，error 类型也是提前定义，并且全局有效的。

使用的最多的 error 实现是 errors 包中的未导出类型 errorString：

```go
// errorString is a trivial implementation of error.
type errorString struct {
    s string
}

func (e *errorString) Error() string {
    return e.s
}
```

你可以使用 errors.New 函数构建一个 error 值，该函数会把一个字符串转换成一个 errorString，返回值是一个 error 类型。

```go
// New returns an error that formats as the given text.
func New(text string) error {
    return &errorString{text}
}
```

下面是使用 errors.New 的示例：

```go
func Sqrt(f float64) (float64, error) {
    if f < 0 {
        return 0, errors.New("math: square root of negative number")
    }
    // implementation
}
```

当调用者传入一个负数时，函数会返回一个非 nil 的 error 值（实际类型为 errorString）。同时，调用者可以使用 error 的 Error 方法得到错误字符串，或者直接打印：

```go
f, err := Sqrt(-1)
if err != nil {
    fmt.Println(err)
}
```

fmt 包内部会调用 error 的 Error 方法。

error 值有必要对一次代码执行做一次总结，应当尽可能描述错误的细节，如“open /etc/passwd: permission denied”，而不是“permission denied.”。Sqrt 的错误返回需要表明函数传递了一个无效的参数，并且要具体到无效的值。

为了加个无效参数的信息，fmt 包的 Errorf 方法派上了用场。该方法会按照 Printf 方法的规则格式化字符串并返回一个 error 类型。

```go
if f < 0 {
    return 0, fmt.Errorf("math: square root of negative number %g", f)
}
```

在大部分场景下，使用 fmt.Errorf 已经足够了，但由于 error 是一个接口类型，所以任意数据结构都可以作为 error 的值（只要实现了 Error 方法）。这样的话，调用者可以尽可能知道错误的具体信息。

假设，调用者想通过 recover 函数捕获到传入 Sqrt 函数的无效参数，那可以自定义一种错误类型，而不是使用默认的 errors.errorString。

```go
type NegativeSqrtError float64

func (f NegativeSqrtError) Error() string {
    return fmt.Sprintf("math: square root of negative number %g", float64(f))
}
```

在捕获到错误时，我们可以使用类型断言还得到无效参数并对其进行适当的处理。与之相对的是，获得错误时仅仅使用 fmt.Println 或 log.Fatal 进行打印，这并不会改变程序的行为。

正如另一个示例，json 包中定义了 SyntaxError 类型，该类型在 json.Decode 函数解析到错误的 JSON 语法时返回。

```go
type SyntaxError struct {
    msg    string // description of error
    Offset int64  // error occurred after reading Offset bytes
}

func (e *SyntaxError) Error() string { return e.msg }
```

在 Error 方法中，Offset 字段没有被使用到。但调用者可以组织其它信息（文件、行）来创建新的错误消息。

```go
if err := dec.Decode(&val); err != nil {
    if serr, ok := err.(*json.SyntaxError); ok {
        line, col := findLine(f, serr.Offset)
        return fmt.Errorf("%s:%d:%d: %v", f.Name(), line, col, err)
    }
    return err
}
```

实现 error 接口只需要定义一个 Error 方法，当然 error 实现也可以有其它的方法。比如在 net 包中，Error 实现了 error 接口，同时自己也是一个接口并定义了其它方法：

```go
package net

type Error interface {
    error
    Timeout() bool   // Is the error a timeout?
    Temporary() bool // Is the error temporary?
}
```

客户端可以为 net.Error 做测试并通过类型断言来判定当前的网络错误是不是暂时的。比如，web 爬虫遇到临时的网络错误时可以选择休眠或者停止执行。

```go
if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
    time.Sleep(1e9)
    continue
}
if err != nil {
    log.Fatal(err)
}
```

# 简化重复错误处理

在 Go 中，错误处理非常重要。同时，Go 的语言设计和约定也鼓励开发者显式地处理错误。在某些情况下，显式处理错误会使代码冗余，但幸运的是，一些编码技术可以最小化重复错误处理的代码。

试想一个 HTTP handler 的应用引擎（App Engine），handler 从数据库中获得数据并渲染到视图：

```go
func init() {
    http.HandleFunc("/view", viewRecord)
}

func viewRecord(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
    record := new(Record)
    if err := datastore.Get(c, key, record); err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    if err := viewTemplate.Execute(w, record); err != nil {
        http.Error(w, err.Error(), 500)
    }
}
```

上述方法处理了由 datastore.Get 函数和 viewTemplate.Execute 函数返回的 error。在这两种情况下，服务端给用户返回一个 HTTP 状态为 500 的消息。这看上去是组织良好的代码，但如果添加更多的 handler，你会发现存在大量的重复代码。

为了减少重复，我们可以定义 HTTP appHandler 类型，该类型为一个函数，函数的返回值是一个 error。

```go
type appHandler func(http.ResponseWriter, *http.Request) error
```

然后，我们修改 viewRecord 函数：

```go
func viewRecord(w http.ResponseWriter, r *http.Request) error {
    c := appengine.NewContext(r)
    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
    record := new(Record)
    if err := datastore.Get(c, key, record); err != nil {
        return err
    }
    return viewTemplate.Execute(w, record)
}
```

相较于上一个版本，该版本会更加简单。但是 http 包并不能正确理解返回 error 的 handler，为了解决这一问题，该类型可以实现 http.Hanlder 接口的 ServeHTTP 方法：

```go
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if err := fn(w, r); err != nil {
        http.Error(w, err.Error(), 500)
    }
}
```

ServeHTTP 方法调用了 appHandler 函数，如果有错误则返回给用户。请注意，这里方法的接收都也是一个函数 fn。

现在我们要注册 handler 函数时，不使用默认的 http.HanlderFunc 类型。

```go
func init() {
    http.Handle("/view", appHandler(viewRecord))
}
```

在这种方式下，可以对 500 的错误进行统一管理。程序可以给用户一个友好的 500 消息，同时我们还需要更好地记录错误信息。

我们可以定义 appError 结构，该结构包含了一个 error 以及其它字段：

```go
type appError struct {
    Error   error
    Message string
    Code    int
}
```

接下来，我们修改 appHandler 的返回值类型：

```go
type appHandler func(http.ResponseWriter, *http.Request) *appError
```

然后，让 appHandler.ServeHTTP 方法给用户显示 appError.Message，并把错误信息打印在终端：

```go
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if e := fn(w, r); e != nil { // e is *appError, not os.Error.
        c := appengine.NewContext(r)
        c.Errorf("%v", e.Error)
        http.Error(w, e.Message, e.Code)
    }
}
```

最后，我们更新 viewRecord：

```go
func viewRecord(w http.ResponseWriter, r *http.Request) *appError {
    c := appengine.NewContext(r)
    key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
    record := new(Record)
    if err := datastore.Get(c, key, record); err != nil {
        return &appError{err, "Record not found", 404}
    }
    if err := viewTemplate.Execute(w, record); err != nil {
        return &appError{err, "Can't display record", 500}
    }
    return nil
}
```

这个版本的 viewRecord 与上一个版本有相同行数的代码，但不同行错误处理又包含其具体的含义。
不止于此，在程序中，我们可以进一步提高错误处理的效率，如：

* 以 HTML 模板的方式展示错误；
* 使用栈追踪技术更好地进行调试；
* 为 appError 编写一个构造函数来存储栈信息；
* 使用 recover 和 panic 设计错误恢复机制；

# 结论

合适的错误处理是好软件的必要项，运用本文介绍的技术一定可以编写更可靠的 Go 代码。
