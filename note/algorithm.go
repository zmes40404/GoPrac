package note

import (
	"fmt"
	"math/rand"
	"time"
)

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

// 7.3 Sort
// 7.3.1 Bubble Sort
func bubbleSort(s []int) { // 這是引用類型直接就地修改，不需要返回
	lastIndex := len(s)-1
	for i:=0; i<lastIndex; i++ { // 外層i: 控制輪數(總共要比幾輪) -> 一個長度為 n 的陣列，最多只要比 n-1 輪，每一輪會「把一個最大的值移到對的位置」，所以第 i 輪後，右邊的 i 個數已經是正確位置，不用再比較。
		for j:=0; j<lastIndex-i; j++ {	// 內層j: 控制當前輪要比哪些 index -> 你要從 j = 0 開始往右，讓最大值一路浮上去，但尾端的 i 個數已經排好了，所以不用再比，所以比較範圍只到：j < n - 1 - i，也就是 j < lastIndex - i
			// "j<lastIndex-i"超直觀記憶法: 「每輪比完，最右邊有 i 個不用再比！」
			fmt.Printf("第%d輪 比較 s[%d]=%d 和 s[%d]=%d\n", i, j, s[j], j+1, s[j+1])
			if s[j] > s[j+1] {	// 如果">"改成"<"就被成降序數列
				tmp := s[j+1]
				s[j+1] = s[j]
				s[j] = tmp
			}
		}
	}
}

func reverseBubbleSort(s []int) {
	lastIndex := len(s) - 1
	for i := 0; i < lastIndex; i++ {
		for j := lastIndex; j > i; j-- {
			fmt.Printf("第%d輪 比較 s[%d]=%d 和 s[%d]=%d\n", i, j, s[j], j-1, s[j-1])
			if s[j] < s[j-1] {
				s[j], s[j-1] = s[j-1], s[j] // 小的往左交換
			}
		}
	}
}

// 7.3.2 Selection Sort -> 不用每個element去進行大小交換，而是用一個max index 或 min index來記錄最大或最小的element，大小比較完之後確定max或min index的值，單圈只要交換一次即可
func SelectionSort(s []int) { // 這是引用類型直接就地修改，不需要返回
	lastIndex := len(s)-1
	for i:=0; i < lastIndex; i++ { // 外層i: 控制輪數(總共要比幾輪) -> 一個長度為 n 的陣列，最多只要比 n-1 輪，每一輪會「把一個最大的值移到對的位置」，所以第 i 輪後，右邊的 i 個數已經是正確位置，不用再比較。
		
		// 假設array最右邊都是設定最大的element
		maxIndex := lastIndex-i		// 把max index給記錄下來
		for j:=0; j<lastIndex-i; j++ {	// 
			if s[j] > s[maxIndex] {
				maxIndex = j		// 透過大小去比，如果大小有變化，就直接去更動max index，這樣跟bubble sort比起來可以減少每一個element都要去交換的過程(減少資料的寫入次數)
			}
		}
		if maxIndex != lastIndex-i {
			s[maxIndex], s[lastIndex-i] = s[lastIndex-i], s[maxIndex]	// Golang是可以直接這樣寫互換的方式
		}

		// 假設array最左邊都是設定最小的element
		// minIndex := i
		// for j:=i+1; j<=lastIndex; j++ {
		// 	if s[j] < s[minIndex] {
		// 		minIndex = j
		// 	}
		// }

		// if minIndex != i {
		// 	s[i], s[minIndex] = s[minIndex], s[i]
		// }
	}
}

func Sort() {
	n := 10
	s := make([]int, n)
	seedNum := time.Now().UnixNano()
	for i:=0; i<n; i++ {
		//rand.Seed(seedNum)
		s[i] = rand.Intn(101)
		seedNum++
	}
	fmt.Println("排序前:", s)
	// bubbleSort(s)
	// reverseBubbleSort(s)
	SelectionSort(s)
	fmt.Println("排序後:", s)
}