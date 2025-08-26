package util

import (
	"context"
	"github.com/go-redis/redis/v8"
	"sync"
)

var (
	redisClient *redis.Client
	once 		sync.Once
	// context.Context 是從 redis v8 開始被強制要求的。v8 的所有 API 都強制帶入 ctx，這是 Go 開發趨勢，也符合 idiomatic Go 設計。
	ctx 		= context.Background()	// 這是一個空白的 context，永不逾時、不會取消。context.Context 是 Go 語言設計用來在 API 間傳遞 deadline（截止時間）、cancel signal（取消信號）、request-scoped value（請求範圍變數） 的核心工具
)

// InitRedisClient() for initialize Redis client (only once)
func InitRedisClient() {
	// 使用once.Do(...)的好處:
	// 確保初始化邏輯只執行一次: 1.無論在程式中呼叫InitRedisClient()幾次，裡面的redis.NewClient(...)只會被執行一次。2.對於"需要共用資源"如資料庫連線、設定讀取等初始化行為非常重要
	// 執行緒安全(thread-safe): "sync.Once"內建機制鎖，保證即使在多執行緒(goroutines)的環境下，也只會初始化一次，不會有race conditon問題
	// 適用於單例(Singleton)設計模式: 一種Golang實現Singleton(單一實例)的常見方式，確保一份global的Redis Client可被各處使用，但又不會重複建立
	
	// 沒用 once 的情況會怎樣？
	// 每次呼叫InitRedisClient()都會重新建立一個新連線
	// 如果在不同地方不小心又呼叫到一次，可能產生多的Redis client，出現浪費資源、競爭或資源釋放不當的問題
	// 在多Goroutine的環境下，更可能出現race condition(資料競爭)或initialization ordering bug
	// 結論: 初始化全域共享資源(Ex. DB, Redus, Logger等)一定要用"sync.Once"包起來
	once.Do(func() {
		// "redis.NewClient(...)"接收的是一個指標型別 *Options，而不是一個實體值 Options。這代表你必須傳入一個 Options 結構的記憶體位置（也就是指標）
		// 為什麼 NewClient() 要用指標型別？
		// 1. 效能優化: Options結構可能包含很多欄位(Ex. TLS, PoolSize、DB等)，若用傳值(by value)會複製整個struct，比較耗能。2. 彈性設計	
		redisClient = redis.NewClient(&redis.Options{	
			Addr: "localhost:6379",		// "6379"是安裝redis server預設的port
		})
	})
}

// GetRedisClient return the initialized instance of Redis client
func GetRedisClient() *redis.Client {
	if redisClient == nil {
		InitRedisClient()
	}
	return redisClient
}

// GetRedisContext() returns the commonly used context
func GetRedisContext() context.Context {
	return ctx
}