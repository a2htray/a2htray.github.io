---
title: "Windows 下移动文件的坑"
date: 2022-04-24T17:16:44+08:00
draft: false
reward: false
categories:
  - 后端技术
  - Go
tags:
 - file
 - I/O
 - Stackoverflow
---

在 Windows 下，Go 的 `os` 标准库提供的 `Rename` 方法不能跨磁盘移动文件。下面通过问题重现，提供两种解决方案。

<!--more-->

## 问题重现

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	err := os.Rename("D:\\black.txt", "E:\\black-new.txt")
	if err != nil {
		fmt.Println(err)
	}
}
```

执行上面的代码后报错。

```bash
rename D:\black.txt E:\black-new.txt: The system cannot move the file to a different disk drive.
```

从源码中可知，Windows 平台有专门实现的 `file_windows.go` ：

```go
// file_windows.go
func rename(oldname, newname string) error {
	e := windows.Rename(fixLongPath(oldname), fixLongPath(newname))
	if e != nil {
		return &LinkError{"rename", oldname, newname, e}
	}
	return nil
}
```

```go
// syscall_windows.go
func Rename(oldpath, newpath string) error {
	from, err := syscall.UTF16PtrFromString(oldpath)
	if err != nil {
		return err
	}
	to, err := syscall.UTF16PtrFromString(newpath)
	if err != nil {
		return err
	}
	return MoveFileEx(from, to, MOVEFILE_REPLACE_EXISTING)
}
```

```go
// zsyscall_windows.go
func MoveFileEx(from *uint16, to *uint16, flags uint32) (err error) {
	r1, _, e1 := syscall.Syscall(procMoveFileExW.Addr(), 3, uintptr(unsafe.Pointer(from)), uintptr(unsafe.Pointer(to)), uintptr(flags))
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}
```

可见，最终移的动操作是通过系统调用完成的，其中 `Rename` 方法调用了两次 `syscall.UTF16PtrFromString` 方法，返回了两个 `*uint16` 类型的值，再使用 `MoveFileEx` 方法完成移动。

## syscall

查看 `zsyscall_windows.go` 提供的方法，还有一个 `MoveFile` 可以尝试，所以就有了以下代码：

```go
package main

import (
	"fmt"
	"syscall"
)

func main() {
	oldpath := "D:\\black.txt"
	newpath := "E:\\black-new.txt"

	from, _ := syscall.UTF16PtrFromString(oldpath)
	to, _ := syscall.UTF16PtrFromString(newpath)
	fmt.Println(*from, *to)
	err := syscall.MoveFile(from, to)
	if err != nil {
		panic(err)
	}
}
```

```bash
68 69
```

移动操作是成功完成的。

## stackoverflow

同时在 stackoverflow 上，也有开发者提供了移动实现。该实现的过程是借助一个中间文件，比如要移动文件到 E 盘，则先在 E 盘创建一个目标文件（`os.Create`），再把源文件的内容写入到目标文件（`os.Copy`），最后删除源文件（os.Remove），代码如下：

```go
import (
    "fmt"
    "io"
    "os"
)

func MoveFile(sourcePath, destPath string) error {
    inputFile, err := os.Open(sourcePath)
    if err != nil {
        return fmt.Errorf("Couldn't open source file: %s", err)
    }
    outputFile, err := os.Create(destPath)
    if err != nil {
        inputFile.Close()
        return fmt.Errorf("Couldn't open dest file: %s", err)
    }
    defer outputFile.Close()
    _, err = io.Copy(outputFile, inputFile)
    inputFile.Close()
    if err != nil {
        return fmt.Errorf("Writing to output file failed: %s", err)
    }
    // The copy was successful, so now delete the original file
    err = os.Remove(sourcePath)
    if err != nil {
        return fmt.Errorf("Failed removing original file: %s", err)
    }
    return nil
}
```

## 跨平台支持

如果希望应用能够在多个平台上正常运行，可以创建 `file.go` 和 `file_windows.go`，分别为不同平台要编译的源代码，也可以创建一个方法，在方法中对平台进行判断。

```go
func MoveFile(src string, dst string) error {
	if runtime.GOOS == "windows" {
		from, _ := syscall.UTF16PtrFromString(src)
		to, _ := syscall.UTF16PtrFromString(dst)
		return syscall.MoveFile(from, to)
	} else {
		return os.Rename(src, dst)
	}
}
```

## 参考

1. [stackoverflow: Move a file to a different drive with Go](https://stackoverflow.com/questions/50740902/move-a-file-to-a-different-drive-with-go)