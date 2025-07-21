package note

import (
	"errors"
	"testing"
)

// "go test ./note -v" to run tests
func TestIsNotNegative(log *testing.T) {// .T是測bug, .B是測速度效能
	err := errors.New("Is Negative")
	if IsNotNegative(0) {
		log.Log("OK")
	} else {
		// log.Error(err)
		log.Fatal(err) // log.Fatal會直接退出程式，並且不會執行defer函數
	}

	if IsNotNegative(1) {
		log.Log("OK")
	} else {
		log.Error(err)
	}
}

// "go test -bench . -benchmem ./note" to run benchmarks
func BenchmarkIfCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsNotNegative(i)
	}
}

func BenchmarkMathCheck(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = i >= 0 // 直接比較
	}
}