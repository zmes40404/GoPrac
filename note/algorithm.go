package note

import "fmt"

// 7.1 Recursion
var fibonacciRes []int // 用一個slice來紀錄已經計算過的結果。用Space Complexity換取Time Complexity。

func fibonacci(n int) int {
	if n < 3 {
		return 1
	}

	// Resursion所調用的函數在記憶體中是獨立的 -> 往深一層調用一個function，就會分配一個記憶體空間給這個function，function用完後再GC
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

// 7.2 Closure
func closureFunc() func(int) int {	// 這個 closureFunc 的行為是： 定義了一個變數 i := 0 回傳一個匿名函數，這個匿名函數會「記住」外層的變數 i
	i := 0	// 閉包（closure）會「捕捉」其定義環境的變數（即 i）
	return func(n int) int {
		fmt.Printf("本次調用接收到n=%v\n", n)
		i++
		fmt.Printf("匿名工具函數被第%v次調用\n", i)
		return i
	}

}
func Closure() {	
	f := closureFunc()	// 每次呼叫 closureFunc()，會重新建立一個 i := 0 的環境，並把這個環境「包住」傳回去。
	f(2)
	f(4)
	f = closureFunc()
	f(6)
}