package main

import "fmt"

func in(x byte, s []byte) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

func equationsPossible(equations []string) bool {
	m := make(map[byte][]byte, 0)
	notEquations := make([]string, 0)
	for _, equation := range equations {
		if equation[1] == '=' {
			left := equation[0]
			right := equation[1]

			if m[left] == nil {
				m[left] = make([]byte, 0)
			}
			if m[right] == nil {
				m[right] = make([]byte, 0)
			}

			if !in(left, m[right]) {
				m[right] = append(m[right], left)
			}

			if !in(right, m[left]) {
				m[left] = append(m[left], right)
			}

		} else {
			notEquations = append(notEquations, equation)
		}
	}
	fmt.Println(m)
	for _, equation := range notEquations {
		left := equation[0]
		right := equation[3]
		if left == right {
			return false
		}
		leftSet, leftIn := m[left]
		rightSet, rightIn := m[right]

		if leftIn && rightIn {
			if in(left, rightSet) && in(right, leftSet) {
				return false
			}
		}
	}
	return true
}

func main() {
	fmt.Println(equationsPossible([]string{
		"a==b",
		"b!=c",
		"c==a",
	}))

	fmt.Println(equationsPossible([]string{
		"c==c","b==d","x!=z",
	}))

	fmt.Println(equationsPossible([]string{
		"b==b","b==e","e==c","d!=e",
	}))
}