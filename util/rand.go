package util

import (
	"math/rand"
	"time"
)

/* 
	- rand.NewSource(seed)：建立一個新的亂數種子
	- rand.New(...)：基於該 seed 建立一個 區域性的亂數生成器實例（不是共用的 global rand）。
	- 呼叫 RandInt(max) 就能獲得良好的隨機數而不影響全域狀態
	- 使用傳統的rand.Seed(...) 會設定 package-level 的全域亂數種子，不適合多 goroutine 使用。
	- 使用rand.New(...) 能讓每個模組、每個 goroutine 擁有自己的亂數來源，更安全也更可控。
*/
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// get the random number from (0, max)
func RandInt(max int) int {
	return  r.Intn(max)
}