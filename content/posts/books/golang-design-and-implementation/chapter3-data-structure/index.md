---
title: "Note - 第 3 章 数据结构"
date: 2023-05-23T21:27:57+08:00
draft: true
reward: false
categories:
 - Go
 - Reading Note
tags:
 - Go 语言设计与实现
---



## 数组

数组是一种由<font color="red">相同数据类型元素</font>组成的集合，计算机会为数组分配一块<font color="red">连续的内存</font>来保存其中的元素，而程序将可以通过<font color="red">索引</font>来访问数组中的元素。通常，我们从两个维度来描述数组：

1. **元素类型**
2. **数组长度**

如以下代码所示：

```go
[10]int
[20]float64
```

数组在初始化后，其大小将无法改变，并且只有**元素类型和数组大小相同的数组类型**才会被程序认为是<font color="red">同一种类型</font>。

```go
// Go 1.18.1 
// src/cmd/compile/internal/types/type.go
// NewArray returns a new fixed-length array Type.
func NewArray(elem *Type, bound int64) *Type {
	if bound < 0 {
		base.Fatalf("NewArray: invalid bound %v", bound)
	}
	t := newType(TARRAY) // TARRAY 为 Kind 类型 <uint> 的常量
	t.extra = &Array{Elem: elem, Bound: bound}
	t.SetNotInHeap(elem.NotInHeap()) // 根据元素类型判断是否分配在堆中
	if elem.HasTParam() {
		t.SetHasTParam(true)
	}
	if elem.HasShape() {
		t.SetHasShape(true)
	}
	return t
}
```

编译期间的数组类型由 `NewArray` 函数生成。

### 数组初始化

Go 中有两种方式创建数组：

1. 显示指定数组大小
2. 使用 `[...]` 声明数组

如下：

```go
myArr := [3]int{1, 2, 3}
myArr2 := [...]int{1, 2, 3}
```

<font color="red">第 2 种声明方式在编译期间会转换成第 1 种</font>，即**编译器会对数组大小进行推导**。

编译器对两种不同声明方式的不同处理：在第 1 种情况下，编译器可以直接用 `NewArray` 进行数组的创建；在第 2 种情况下，编译器则会通过遍历的方式计算数组大小。<font color="red">第 2 种声明方式是第 1 种声明方式的语法糖</font>。

### 访问与赋值

数组是否发生了越界访问，编译器可以在**静态类型检查**时进行判断（`src/cmd/compile/internal/typecheck/expr.tcIndex` 函数）。

```go
// Go 1.18.1
// tcIndex typechecks an OINDEX node.
func tcIndex(n *ir.IndexExpr) ir.Node {
	// ...
	case types.TSTRING, types.TARRAY, types.TSLICE:
		// ...
		if !n.Bounded() && ir.IsConst(n.Index, constant.Int) {
			x := n.Index.Val()
			if constant.Sign(x) < 0 {
				base.Errorf("invalid %s index %v (index must be non-negative)", why, n.Index)
			} else if t.IsArray() && constant.Compare(x, token.GEQ, constant.MakeInt64(t.NumElem())) {
				base.Errorf("invalid array index %v (out of bounds for %d-element array)", n.Index, t.NumElem())
			} else if ir.IsConst(n.X, constant.String) && constant.Compare(x, token.GEQ, constant.MakeInt64(int64(len(ir.StringVal(n.X))))) {
				base.Errorf("invalid string index %v (out of bounds for %d-byte string)", n.Index, len(ir.StringVal(n.X)))
			} else if ir.ConstOverflow(x, types.Types[types.TINT]) {
				base.Errorf("invalid %s index %v (index too large)", why, n.Index)
			}
		}
	// ...
}
```

数组赋值也会触发越界检查，赋值的过程为：<font color="red">先确定目标数组的地址，再确定目标元素的地址，最后将地址上的值替换成新值</font>。

在中间代码生成期间，编译器会插入运行时方法 `runtime.panicIndex` 调用防止发生越界错误。

## 切片

切片，即动态数组，其长度不固定，<font color="red">当对切片追加元素且容量不足时会自动扩容</font>。切片的声明不需要指定大小，如下：

```go
[]int
[]float64
```

```go
// Go 1.18.1
// src/cmd/compile/internal/types/type.go
// NewSlice returns the slice Type with element type elem.
func NewSlice(elem *Type) *Type {
	if t := elem.cache.slice; t != nil {
		if t.Elem() != elem {
			base.Fatalf("elem mismatch")
		}
		if elem.HasTParam() != t.HasTParam() || elem.HasShape() != t.HasShape() {
			base.Fatalf("Incorrect HasTParam/HasShape flag for cached slice type")
		}
		return t
	}

	t := newType(TSLICE)
	t.extra = Slice{Elem: elem}
	elem.cache.slice = t
	if elem.HasTParam() {
		t.SetHasTParam(true)
	}
	if elem.HasShape() {
		t.SetHasShape(true)
	}
	return t
}
```

编译期间的切片类型由 `NewSlice` 函数生成。