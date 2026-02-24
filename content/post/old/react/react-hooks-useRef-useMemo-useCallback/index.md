---
title: "如何高效地使用 useRef、useMemo 和 useCallback"
date: 2024-01-03T22:38:25+08:00
draft: false
reward: false
categories:
 - 前端技术
 - React
tags:
 - React Hooks
 - TypeScript
 - 翻译
---

这是一篇关于 React Hooks 的技术文章翻译，原文地址：
[https://dev.to/michael_osas/understanding-react-hooks-how-to-use-useref-usememo-and-usecallback-for-more-efficient-code-3ceh](https://dev.to/michael_osas/understanding-react-hooks-how-to-use-useref-usememo-and-usecallback-for-more-efficient-code-3ceh)，翻译不当之处请指正。

**译者评价**

文章主要介绍了 useRef、useMemo 和 useCallback 3 个 React Hook，读者可以通过此文了解 3 种 Hook 的使用方式、场景，但文章也存在一些缺点：

1. <font color="red">内容重复严重，如 useRef 的作用在文章前中后段中均有描述</font>
2. <font color="red">useCallback 的示例举例不当，容易造成读者的困惑</font>

文章的内容只适合 3 个 React Hook 的基础泛读，因为示例足够简单，所以读者读起来并不会太大的压力。

<!--more--->

**正文开始**

在不编写和使用类组件的情况下，React Hooks 是允许开发者在函数组件中管理 state（组件状态）和其它 React 特性的一组方法。React Hooks 改变了开发者编写 React 组件的方式，同时也加强了组件的可复用性和可维护性。本文详细地讨论了 React 3 个内置 Hook，分别为 <font color='blue'>useRef</font>、<font color='blue'>useMemo</font> 和 <font color='blue'>useCallback</font>。通过内置的 Hooks，我们可以很容易地在不同组件中复用相同的、带状态的逻辑，这也使得代码更具可维护性和易读性。

React Hooks 是一组可以在函数组件中管理组件状态、执行特定动作和访问上下方的<font color='blue'>钩子函数</font>，包括 useState、useEffect、useRef、useMemo 和 useCallback 等。

<font color="red">useRef</font> 提供了<font color='blue'>存储可变值</font>的方法，该引用值会在组件 render 期间持续存在。<font color='blue'>useRef 通常用于访问或修改 DOM 元素的属性值</font>，比如输入框的值（value）、容器的滚动位置（offsetLeft、offsetTop）。<font color='blue'>与 useState 不同的是，更新 useRef 的值并不会触发组件的重新渲染</font>，相应地，可以直接通过 useRef 返回引用对象的属性值直接进行访问或修改。

<font color="red">useMemo</font> 可以<font color="blue">缓存函数的返回值，当且仅当依赖改变时，才会重新计算函数的返回值</font>。useMemo 非常适合于计算密集的组件编写，比如组件的渲染依赖于某项复杂的数据结构计算结果。useMemo 的第 2 个参数是一个由依赖项组成的数组，只有其中依赖项改变时，函数才会重新计算。

<font color="red">useCallback</font> <font color="blue">缓存并返回一个函数，当且仅当依赖项改变时，才会重新创建函数实例</font>。<font color="red" style="text-decoration: underline;">useMemo 适用于向子组件传递函数 prop ，从而减少子组件渲染的场景</font>，其第 2 个参数也是一个依赖项数组，若依赖项改变时，函数才会重新创建。

一般情况下，<font color="red" style="text-decoration: underline;">React Hooks 可以将带状态逻辑（stateful logic）组件分解成更小，可组合的函数组件</font>，这使得更容易地编写可复用和可维护的代码。使用 React Hooks 可以<font color="blue">简化编写 React 组件过程，减少重复代码，使组件代码更具可读性</font>。

好了，现在让我们更深入地理解这些 Hooks 吧。

## useRef

在 React 中，useRef 可以创建并存储一个可变的引用值，<font color="blue">相当于在类组件中的 ref 属性</font>，但也有些许不同。

useRef 使用一个初始值作为参数，并返回一个可变的引用对象。引用对象的属性值可以直接地被访问和修改，并且<font color="blue">值的改变不会造成组件的重新渲染</font>。这对于要延后访问或设置 DOM 元素属性的场景十分有用

下面是使用 useRef 来保存引用对象，实现点击按钮完成输入框定焦的一个示例：

```typescript
import { useRef } from "react";

function ExampleUseRef() {
  const inputRef = useRef<HTMLInputElement>(null);
  function handleClick() {
    inputRef?.current?.focus();
  }
  return (
    <div>
      <input type="text" ref={inputRef} />
      <button onClick={handleClick}>Focus Input</button>
    </div>
  );
}
```

上述代码中，ExampleUseRef 是一个使用 useRef <font color="blue">保存 RefObject（T 为 HTMLInputElement）的引用对象</font>，引用的是一个 HTMLInputElement 对象。让我们进一步讨论下代码：

我们调用 <font color="blue">useRef\<HTMLInputElement\>(null)</font> 创建了一个 HTMLInputEelement 的引用，inputRef 的初始值为该引用。

用户点击“Focus Input”按钮时，handleClick 函数会被调用。在 handleClick 函数中，访问 inputRef 对象的 current 属性并调用 focus() 方法。

ExampleUseRef 组件返回一个包含 input 和 button 元素的 div 元素，<font color="blue">input 元素的 ref 属性设置为 inputRef，这样就会创建一个 input 元素的引用</font>。同时，把 button 元素的 onClick 属性设置为 handleClick 函数，当按钮被点击时，该函数会被调用。

综上所述，示例演示了如何在函数组件中创建一个 HTMLInputElement 的引用，然后通过 button 的 onClick 属性设置回调函数，最后在函数中使用引用对象的 current 属性调用 focus() 方法。通过使用 useRef，input 元素的引用会在组件渲染期间持续持有，这对于后期要访问或修改引用对象的属性非常有帮助。

## useMemo

useMemo <font color="blue">缓存了计算密集型函数的执行结果，从而提高程序的性能</font>。缓存技术是根据参数缓存函数的计算结果，如果传入相同的参数时，直接返回缓存结果，避免了重新计算结果。

useMemo 使用两个参数：<font color="blue">一个是返回缓存值的函数，另一个是一组包含依赖项的数组</font>。该函数只有在依赖项改变时，函数才会重新执行。如果依赖项在前后两次渲染相同时，则直接返回前一次缓存的函数返回值。下面是使用 useMemo 的详细示例：

```typescript
import { useMemo } from "react";

function ExampleUseMemo({ a, b }: { a: number; b: number }) {
  const memoizedValue = useMemo(() => {
    return a * b;
  }, [a, b]);

  return <div>Memoized Value: {memoizedValue}</div>;
}
```

在上述示例中，useMemo 缓存了 a 和 b 的计算结果，这就意味着，<font color="blue">如果 a 和 b 不发生改变，函数不会重新执行并且返回缓存值</font>。这么做可以减少程序的计算时间，从而提高程序性能。

然而，值得注意的是，不是所有的场景都适合使用 useMemo。如果计算过程并不复杂或者组件不依赖于缓存值且频繁需要渲染，使用 useMemo 的性能提升并不会太显著，相反地，可能还会造成性能的损失。所以，我们要记住：<font color="blue">如果缓存值是频繁修改的，相比于缓存计算值，直接渲染组件可能是更好的解决方案</font>。

总结下来，useMemo 是一个优化 React 程序的性能有用的工具，但是在使用时，需要综合考虑其优缺点。

## useCallback

useCallback 通过减少组件的非必要渲染从而提升程序的性能，<font color="blue">该 Hook 会缓存一个函数，只有在依赖项改变时，useCallback 才会返回一个新的函数实例</font>。

下面是 useCallback 的基础语法：

```typescript
const memoizedCallback = useCallback(
  () => {
    // 函数体
  },
  [/* 依赖数组 */]
);
```

useCallback 的语法和 useEffect 相似，第 1 个参数是你要缓存的函数，第 2 个参数是一个包含依赖项的可选数组。如果依赖项发生改变，则会重新计算要缓存的函数并返回。<font color="blue">使用 useCallback 的一个好处是，它可以帮助程序减少组件不必要的渲染次数</font>。

下面我们通过一个简单示例进行学习：

```typescript
import { useEffect, useState } from "react";

function ParentComponentWithoutUseCallback() {
  const [count, setCount] = useState(0);
  const handleClick = () => {
    setCount(count + 1);
  };
  return (
    <div>
      <button onClick={handleClick}>Increment count</button>
      <ChildComponentWithoutUseCallback />
    </div>
  );
}

function ChildComponentWithoutUseCallback() {
  const [value, setValue] = useState(0);
  const expensiveFunction = () => {
    console.log("call expensiveFunction");
    // 计算密集
  };
  useEffect(() => {
    expensiveFunction();
  }, [expensiveFunction]);
  return <div>{value}</div>;
}
```

这是一个包含两个 React 组件（ParentComponentWithoutUseCallback 和 ChildComponentWithoutUseCallback）的示例：

ParentComponentWithoutUseCallback：组件使用 useState 定义了一个状态变量 count，其初始值为 0，同时定义了一个 handleClick 的点击回调函数，其内部实现使用 setCount 实现 count 状态值的加 1 操作。

ParentComponentWithoutUseCallback 组件渲染时，页面上会渲染出一个带 button 元素的 div 元素。当 button 发生点击时，handleClick 函数会增加 count 值，但同时 ChildComponentWithoutUseCallback 组件也会发生渲染。

ChildComponentWithoutUseCallback：组件使用 useState 定义了一个状态变量 value，其初始值为 0，同时还定义了一个模拟复杂计算的函数 expensiveFunction。

当 ChildComponentWithoutUseCallback 组件完成挂载时，useEffect 内部会执行 expensiveFunction 函数。useEffect 函数使用两个参数，分别为一个函数和一个包含依赖项的数组。在这个示例中，第 1 个参数是 expensiveFunction 函数，第 2 个参数是包含 expensiveFunction 函数的数组。将 expensiveFunction 函数作为 useEffect 的依赖项，只有当 expensiveFunction 引用改变时，useEffect 才会重新执行 expensiveFunction 函数。

最后，ChildComponentWithoutUseCallback 组件返回一个带 value 状态变量的 div 元素。当然，这个示例也演示了如何使用 useState 和 useEffect 来管理组件状态和执行 efffect（副作用）。当 ParentComponentWithoutCallback 组件更新了 count 状态变量时，ChildComponentWithoutCallback 组件也会根据 useEffect 再次执行计算密集的 expensiveFunction 函数。

为了优化上述代码，我们计划使用 useCallback 来缓存计算密集的函数，只有当依赖项改变时，计算密集函数才会重新计算：

```typescript
import { useEffect, useState } from "react";

function ParentComponentUseCallback() {
  const [count, setCount] = useState(0);
  const handleClick = () => {
    setCount(count + 1);
  };
  const memoizedFunction = useCallback(() => {
    console.log("call memoizedFunction");
    // 计算密集
  }, []);
  return (
    <div>
      <button onClick={handleClick}>Increment count</button>
      <ChildComponentUseCallback expensiveFunction={memoizedFunction} />
    </div>
  );
}

function ChildComponentUseCallback({
  expensiveFunction,
}: {
  expensiveFunction: () => void;
}) {
  const [value, setValue] = useState(0);
  useEffect(() => {
    expensiveFunction();
  }, [expensiveFunction]);
  return <div>{value}</div>;
}
```

在上述示例中，除了使用 useCallback 定义 memoizedFunction 缓存函数，ParentComponentUseCallback 与之前的父组件类似。useCallback 用于缓存函数，所以不会每次渲染中重新创建。也正因如此，<font color="blue">如果将缓存函数作为子组件的 prop，就可以减少子组件不必要的渲染次数</font>。

在 ChildComponentUseCallback 子组件中，除了定义了一个 expensiveFunction prop 外，其它的与之前的子组件类似。因为 expensiveFunction 函数作为 prop 传入子组件，所以需要在父组件中定义 expensiveFunction 函数。在这个示例中，只有当 expensiveFunction 改变时，useEffect 才会重新执行 expensiveFunction。然而，在父组件中，我们使用了 useCallback 缓存计算密集函数，当父组件发生渲染，memoizedFunction 不会发生改变，从而避免了子组件的再次渲染。

最后，子组件返回一个带 value 状态变量的 div 元素。

简单点说，示例演示了在函数组件中如何使用 useCallback 才避免组件的多次渲染。同时，在示例中也<font color="blue">结合了 useEffect 技术，将计算密集函数作为组件的 prop，从而减少计算密集函数的执行</font>。

## 最终想法

总结来说，React Hooks 给开发者提供了在函数组件编写中管理组件状态更为简单的方式。

useRef 创建了一个 DOM 元素的可变的引用，在渲染中持续存在，允许开发者在造成组件重新渲染下直接访问或修改元素的属性。useMemo 可用于缓存计算密集型函数的返回值，只有当依赖项改变时才会重新计算值。useCallback 则是通过缓存一个函数，也是只当依赖项改变时，才会重新创建函数实例。

使用这些内置的 Hooks，如 useRef、useMemo 和 useCallback，开发者可以更高效管理组件状态，执行特定动作和访问上下文。通过将状态逻辑拆分到更小、可组合的函数组件中，使用 Hooks 编写出的组件更加简单，其可复用性、可读性得到加强。