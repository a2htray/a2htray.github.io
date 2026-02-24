---
title: "递归版本的归并排序 Go & Python 实现"
date: 2023-04-23T23:38:47+08:00
draft: false
reward: false
categories:
 - 算法
 - 实现
tags:
 - 排序算法
 - 归并排序
 - 递归
 - Go
 - Python
---

递归版本的归并排序。

<!--more-->

## Go

```go
/**
 * 递归版本的归并排序
 */
package main

import "fmt"

func rMerge(left []int, right []int) []int {
	ret := make([]int, 0, len(left)+len(right))
	i := 0
	j := 0
	for i < len(left) && j < len(right) {
		if left[i] < right[j] {
			ret = append(ret, left[i])
			i++
		} else {
			ret = append(ret, right[j])
			j++
		}
	}

	if i < len(left) {
		for _, v := range left[i:] {
			ret = append(ret, v)
		}
	}

	if j < len(right) {
		for _, v := range right[j:] {
			ret = append(ret, v)
		}
	}
	return ret
}

func rMergeSort(s []int) []int {
	if len(s) <= 1 {
		return s
	}
	mid := len(s) / 2
	left := rMergeSort(s[:mid])
	right := rMergeSort(s[mid:])
	return rMerge(left, right)
}

func main() {
	s1 := []int{9, 8, 6, 6, 5, 4, 3, 2, 1}
	fmt.Println(rMergeSort(s1))
	s2 := []int{3, 1, 4, 4, 6, 5, 8}
	fmt.Println(rMergeSort(s2))
}
```

## Python

```python
"""
递归版本的归并排序
"""
from typing import List


def r_merge(left: List[int], right: List[int]) -> List[int]:
    ret = []
    i = 0
    j = 0
    while i < len(left) and j < len(right):
        if left[i] < right[j]:
            ret.append(left[i])
            i += 1
        else:
            ret.append(right[j])
            j += 1

    if i < len(left):
        ret = ret + left[i:]
    if j < len(right):
        ret = ret + right[j:]

    return ret


def r_merge_sort(arr: List[int]) -> List[int]:
    if len(arr) <= 1:
        return arr

    mid = len(arr) // 2
    left = r_merge_sort(arr[:mid])
    right = r_merge_sort(arr[mid:])

    return r_merge(left, right)


if __name__ == '__main__':
    arr1 = [9, 8, 6, 6, 5, 4, 3, 2, 1]
    print(r_merge_sort(arr1))

    arr2 = [3, 1, 4, 4, 6, 5, 8]
    print(r_merge_sort(arr2))
```

