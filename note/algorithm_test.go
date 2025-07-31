package note

import (
	"math/rand"
	"time"
	"testing"
)

func NewSlice() []int {
	n := 500000
	s := make([]int, n)
	seedNum := time.Now().UnixNano()
	for i:= 0; i < n; i++ {
		//  如果你是為了「用特定種子取得固定的亂數序列」（例如單元測試或重現某個情況) -> r := rand.New(rand.NewSource(seed))
		s[i] = rand.Intn(10001)
		seedNum++
	}
	return  s
}

func BenchmarkBubbleSort(log *testing.B) {
	bubbleSort(NewSlice())
}

func BenchmarkSelectionSort(log *testing.B) {
	SelectionSort(NewSlice())
}

func BenchmarkInsertionSort(log *testing.B) {
	InsertionSort(NewSlice())
}

func BenchmarkQuickSort(log *testing.B) {
	s := NewSlice()
	// QuickSort(s, 0, len(s)-1)
	QuickSortLomuto(s, 0, len(s)-1)
}

