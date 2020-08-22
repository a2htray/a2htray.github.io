package main

import "fmt"
import "math"

func divide(dividend int, divisor int) int {
	if dividend == 0 {
		return 0
	}
	flag := 1
	if (dividend > 0 && divisor < 0) || dividend < 0 && divisor > 0 {
		flag = -1
	}

	dividend = int(math.Abs(float64(dividend)))
	divisor = int(math.Abs(float64(divisor)))

	i := 0
	for {
		if divisor > dividend {
			break
		}

		dividend -= divisor
		i++
	}
	if flag == 1 {
		return i
	}

	return -i
}

func main() {
	fmt.Println(divide(10, 3))
	fmt.Println(divide(7, -3))
	fmt.Println(divide(1, 1))
	fmt.Println(divide(-2147483648, -1))
	
}