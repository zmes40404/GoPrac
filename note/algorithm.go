package note

import "fmt"

// 7.1 Recursion
var fibonacciRes []int // 用一個slice來紀錄已經計算過的結果。用Space Complexity換取Time Complexity。

func fibonacci(n int) int {
	if n < 3 {
		return 1
	}

	// return fibonacci(n-2) + fibonacci(n-1) // 傳統的方式，Time Complexity特別爛 -> 尤其是n很大的情況

	if fibonacciRes[n] == 0 {
		fibonacciRes[n] = fibonacci(n-2) +  fibonacci(n-1)
	
	}
	return fibonacciRes[n]
}

func Recursion() {
	// fmt.Printf("第%v位費波那契數列為%v\n", 45, fibonacci(45)) // 傳統的方式，Time Complexity特別爛 -> 尤其是n很大的情況

	n := 45
	fibonacciRes = make([]int, n+1)
	fmt.Printf("第%v位費波那契數列為%v\n", n, fibonacci(n))
}