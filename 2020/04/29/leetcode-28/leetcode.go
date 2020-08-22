package main

import "fmt"

func strStr(haystack string, needle string) int {
	hLen := len(haystack)
	nLen := len(needle)
	i := 0
	for ; i < hLen - nLen + 1; i++ {
		if haystack[i:i+nLen] == needle {
			return i
		}
	}
	return -1
}

func main() {
	fmt.Println(strStr("hello", "ll"))
	fmt.Println(strStr("aaaaa", "bba"))
	fmt.Println(strStr("a", ""))
	fmt.Println(strStr("", ""))
	fmt.Println(strStr("mississippi", "a"))
	fmt.Println(strStr("aaa", "aaa"))
}