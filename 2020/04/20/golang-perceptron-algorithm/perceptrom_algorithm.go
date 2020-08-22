package main

import (
	"fmt"
	"math/rand"
)

// Samples 生成样本数据
// 因为感知器算法需要数据集线性可分
// 为方便起见，使用给定直接函数 f
// 若 f(x) > 0 标记为 1
// 若 f(x) < 0 标记为 -1
// * 10 的目的在于让数据稍微大一些
func Samples(n int, m int) [][]float32 {
	// 给定分隔数据集函数，用于判断
	b := -rand.Float32() * 10
	w := make([]float32, 0)
	for i := 0; i < m; i++ {
		if i % 2 == 0 {
			w = append(w, rand.Float32() * 10)
		} else {
			w = append(w, -rand.Float32() * 10)
		}
	}

	f := func (x ...float32) float32 {
		var result float32
		for i := 0; i < m; i++ {
			result += x[i] * w[i]
		}
		result += b
		return result
	}

	samples := make([][]float32, 0)

	for i := 0; i < n; i++ {
		sample := make([]float32, 3)
		sample[0] = rand.Float32() * 10
		sample[1] = rand.Float32() * 10
		if f(sample[0:3]...) > 0 {
			sample[2] = 1
		} else {
			sample[2] = -1
		}
		samples = append(samples, sample)
	}

	return samples
}

// 判断算法是否结束
func IsFinish(indicator []int) bool {
	for _, i := range indicator {
		if i == 0 {
			return false
		}
	}
	return true
}

// 指示器 Slice 全部设置为 0
func Zeros(indicator []int) []int {
	for i, _ := range indicator {
		indicator[i] = 0
	}
	return indicator
}

func DisplayResult(samples [][]float32, w []float32, b float32)  {
	fmt.Printf("Finished: w1 = %f, w2 = %f, b = %f\n", w[0], w[1], b)
	flag := ">"
	for _, sample := range samples {
		value := w[0] * sample[0] + w[1] * sample[1] + b
		if value > 0 {
			flag = ">"
		} else {
			flag = "<"
		}
		fmt.Printf(
			"w1*x1 + w2*x2 + b = %f*%f + %f*%f + %f = %f %s 0, and y = %f\n",
			w[0], sample[0], w[1], sample[1], b,
			value, flag, sample[2])
	}
}

func main()  {
	samples := Samples(100, 2)
	indicator := make([]int, 100)

	// 输出数据集所有样本
	for _, sample := range samples {
		fmt.Println(sample)
	}

	// 随机生成权重和偏置
	w := make([]float32, 0)
	w = append(w, rand.Float32())
	w = append(w, rand.Float32())
	b := rand.Float32()

	// 输出初始权重与偏置
	fmt.Printf("initialized: w1 = %f, w2 = %f, b = %f\n", w[0], w[1], b)

	i := 0
	for {

		if samples[i][0] * w[0] + samples[i][1] * w[1] + b > 0 && samples[i][2] == -1 {
			w[0] = w[0] - samples[i][0]
			w[1] = w[1] - samples[i][1]
			b = b - 1

			i += 1
			if i == 100 {
				i = 0
			}
			indicator = Zeros(indicator)
			continue
		} else if samples[i][0] * w[0] + samples[i][1] * w[1] + b < 0 && samples[i][2] == 1 {
			w[0] = w[0] + samples[i][0]
			w[1] = w[1] + samples[i][1]
			b = b + 1
			i += 1

			if i == 100 {
				i = 0
			}
			indicator = Zeros(indicator)
			continue
		} else {
			indicator[i] = 1

			i += 1
			if i == 100 {
				i = 0
			}
		}

		if IsFinish(indicator) {
			break
		}
	}

	DisplayResult(samples, w, b)
}