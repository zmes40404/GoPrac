package note

import (
	"fmt"
	"goprac/util"
	"math/rand"
	"sort"
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

// 7.3.3 Insertion Sort -> 將資料列假設分成已排序和未排序的兩部分，每次從未排序的資料中，挑選出一個元素，插入到已排序的資料中，直到所有的資料都已排序完成。
func InsertionSort(s []int) {
	for i:=1; i<len(s); i++ {
		tmp := s[i] 	// tmp就是當下要比較處理的element，也就是尚未排序數列中的第一個element
		j := i-1	// j就是已經排序數列中的最後一個
		for ; j>=0 && s[j] > tmp; j-- {	// 已排序的數列開始找出大小剛好要放tmp值的位置
			s[j+1] = s[j]	// 大的數字往右複製，用複製(記憶體寫入1次)的效能比用swap的好，因為swap的底層是用變數tmp來裝其中一個要交換的值，所以底層其實是記憶體寫入3次
		}
		if j != i-1 {	// 這個判斷用來避免剛好順序都排序好的片段進行額外的記憶體資料寫入
			s[j+1] = tmp	//「插入點」 = 比它小的那個右邊，如果剛好在s[0]的話，j就會是-1，-1的右邊就是0
			fmt.Println("s=", s)
		}
	}
}

// 7.3.4 Quick Sort
func QuickSort(s []int, leftIndex, rightIndex int) {
	if leftIndex >= rightIndex { // base case: 不需排序。如果你要排序的「區間」已經沒有元素（或只剩一個元素），就不需要再排序了，直接 return。
		return
	}
	
	if leftIndex < rightIndex {
		middle := s[rightIndex]		// 一開始假設最右邊的數是中位數
		var rightslice []int	// 額外的記憶體來儲存element大於pivot的數列
		l := leftIndex
		for i:=leftIndex; i<rightIndex; i++ {
			if s[i] > middle {	// 如果該element大於pivot
				rightslice = append(rightslice, s[i])	// 把該element丟到新create的slice中
			} else {
				s[l]=s[i]	// 分割邏輯中「寫入左邊區塊」的重要步驟: 小於等於 pivot 的往左擠過去
				l++
			}
		}
		s[l] = middle	// 把原本最後一個當成中位數的值給插回s中間的位置
		copy(s[l+1:], rightslice)	// 這個QuickSort是比較像merge sort要拼回去的，傳統的QuickSort都是在每個recurive裡面直接做swap
		if leftIndex < l-1 {
			QuickSort(s, leftIndex, l-1)
		}
		if l+1 < rightIndex {
			QuickSort(s, l+1, rightIndex)
		}
	}
}

// 經典版：Lomuto 分割法寫法（更常見）
func QuickSortLomuto(s []int, low, high int) {
	if low >= high {	// 如果你要排序的「區間」已經沒有元素（或只剩一個元素），就不需要再排序了，直接 return。
		return
	}
	pivotIndex := partition(s, low, high)
	QuickSortLomuto(s, low, pivotIndex-1)
	QuickSortLomuto(s, pivotIndex+1, high)
}

func partition(s []int, low, high int) int {
	pivot := s[high] // 選擇最右邊為 pivot
	i := low

	for j := low; j < high; j++ {
		if s[j] < pivot {
			s[i], s[j] = s[j], s[i]		// 直接用交換的
			i++
		}
	}
	s[i], s[high] = s[high], s[i] // 將原本選定最後一個的中位數給換回中間正確的位置
	return i // 返回 pivot 最終位置
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
	// SelectionSort(s)
	// InsertionSort(s)
	QuickSort(s, 0, len(s)-1)
	// QuickSortLomuto(s, 0, len(s)-1)
	fmt.Println("排序後:", s)
}

// 7.5.2 Binary Search
func BinarySearch(s []int, key int) int {
	startIdx := 0
	endIdx := len(s)-1
	midIdx := 0
	for startIdx <= endIdx {
		midIdx = startIdx + (endIdx-startIdx)/2	// 這種取midIdx相較於值觀的"(endIdx-startIdx)/2"，可以避免可能endIdx-startIdx的值太大超出int64導致overflow。
		if s[midIdx] < key {
			startIdx = midIdx+1
		}else if s[midIdx] > key {
			endIdx = midIdx-1
		} else {
			return  midIdx
		}				
	}
	return -1
}

func BinarySearchTest() {
	// make([]int, N)	預分配記憶體，效能較好	已知大小，會填滿資料。避免了多次擴容與內存拷貝
	s := make([] int, util.RandInt(1000)+1)	// 為什麼要+1 -> 假設 util.RandInt(1000) 回傳 0，那 0+1 = 1，這樣可以保證 slice 至少有長度 1，不會是空的
	for i:=0; i< len(s); i++ {
		s[i] = util.RandInt(1000)	// 將第 i 個元素填入隨機數值（範圍 0 到 999）
	}
	sort.Ints(s)
	i:=BinarySearch(s, 555)
	if i == -1 {
		fmt.Println("沒有找到555")
	}else {
		fmt.Printf("555的下標為%v", i)
	}
}