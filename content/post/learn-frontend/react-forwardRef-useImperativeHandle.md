+++
date = '2026-05-26T14:22:47+08:00'
draft = false
title = 'forwardRef - 转发 Ref 和 useImperativeHandle - 自定义 Ref'
categories = ['前端技术', 'React']
tags = ['React', 'React Hooks', 'forwardRef', 'useImperativeHandle']
toc = true
+++

forwardRef 和 useImperativeHandle 是 React 提供的两个 Hook，分别用于：

* forwardRef 是 React 提供的一个 Hook，解决父子组件传递 Ref 的问题；
* useImperativeHandle 是 React 提供的一个 Hook，用于自定义 Ref，将子组件的方法暴露给父组件。

## forwardRef

举个例子：

```tsx
import React, { useRef, forwardRef } from 'react';

const ChildInput = forwardRef((props, ref) => {
  return <input ref={ref} placeholder="子组件输入框" />;
});

const Parent = () => {
  const inputRef = useRef(null);

  const focusInput = () => {
    inputRef.current.focus();
  };

  return (
    <div>
      <ChildInput ref={inputRef} />
      <button onClick={focusInput}>聚焦子组件输入框</button>
    </div>
  );
};

export default Parent;
```

在父组件（Parent）中定义了 inputRef 的引用，并且子组件（ChildInput）使用 `forwardRef` 包裹，接收了 ref 参数，注意参数声明，如下：

```tsx
const ChildInput = forwardRef((props, ref) => {
  // ...
});
```

第 2 个参数是 ref 参数，用于接收父组件传递来的 Ref。

在子组件中，内部使用 ref 参数绑定到 input DOM 上，但“聚集”是在父组件中通过 `inputRef.current.focus()` 来实现的。

所以，forwardRef 的适用场景：

* 子组件需要在父组件操作 DOM 或其他元素；
* 子组件需要在父组件中调用方法；

## useImperativeHandle

**必须搭配 forwardRef 使用**。

举个例子：

```tsx
import React, { useRef, forwardRef, useImperativeHandle } from 'react';

const Child = forwardRef((props, ref) => {
  const hello = () => {
    alert('我是子组件的方法：你好！');
  };

  useImperativeHandle(ref, () => ({
    hello,
  }));

  return <div>我是子组件</div>;
});

export default function Parent() {
  const childRef = useRef(null);

  return (
    <div style={{ padding: 20 }}>
      <Child ref={childRef} />
      <button onClick={() => childRef.current.hello()}>
        调用子组件方法
      </button>
    </div>
  );
}
```

子组件的方法 `hello` 被暴露给父组件，父组件可以通过 `childRef.current.hello()` 来调用。

## 总结

又学了两个 React Hooks，分别是 **forwardRef** 和 **useImperativeHandle**。
