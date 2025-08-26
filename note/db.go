package note

import (
	_"context"
	"fmt"
	"time"
	"github.com/syndtr/goleveldb/leveldb"
	leveldbUtil "github.com/syndtr/goleveldb/leveldb/util"
	"github.com/go-redis/redis/v8"	// go-redis/redis/v8 是 Go 語言的 Redis 用戶端函式庫版本，與 Redis Server 本身的版本無直接關聯
	"goprac/util"
)

// 11.1 Basic Usage of LevelDB
func LeveldbBasic() {
	db, err := leveldb.OpenFile("leveldb", nil)	// 這個"OpenFile"如果路徑有填寫資料夾的話，資料夾不存在則新建，資料夾已經存在的話則打開。
	if err != nil {
		panic(err)
	}
	defer db.Close()	 // 一般常見寫法會將打開db寫在init中，關閉db寫一個Close func在main中執行

	db.Put([]byte("user-1"), []byte("{\"username\":\"1\"}"), nil)
	// db.Delete([]byte("user-1"), nil)
	data, _ := db.Get([]byte("user-1"), nil)
	fmt.Println("data=", string(data))

	// 判斷資料是否存在
	ok, _ := db.Has([]byte("user-1"), nil)
	fmt.Println("Has \"user-1\" ?", ok)

	//批量操作
	batch := new(leveldb.Batch)
	batch.Put([]byte("user-2"), []byte("{\"username\":\"2\"}"))
	batch.Put([]byte("user-3"), []byte("{\"username\":\"3\"}"))
	batch.Delete([]byte("user-1"))
	batch.Put([]byte("user-1"), []byte("{\"username\":\"n1\"}"))
	err = db.Write(batch, nil) // 將記憶體中的緩存資料寫入硬碟，db的修改都建議在寫入硬碟後再修改，不要在緩存記憶體中時就進行修改
	if err != nil {
		panic(err)
	}
	data, _ = db.Get([]byte("user-3"), nil)
	fmt.Println("data=", string(data))
}

// 11.1.4 LevelDB for loop
func LeveldbIterate() {
	db, err := leveldb.OpenFile("leveldb", nil)
	if err != nil{
		panic(err)
	}
	defer db.Close()
	batch := new(leveldb.Batch)
	for i:=1; i<11; i++ {
		batch.Put(
			[]byte(fmt.Sprintf("user-%v", i)),
			[]byte(fmt.Sprintf("{\"name\":\"u%v\"}", i)))
	}
	db.Write(batch, nil)

	// 遍歷指定範圍的資料、&leveldbUtil.Range處填nil為完整資料庫
	iter := db.NewIterator(
		&leveldbUtil.Range{
			Start: []byte("user-3"), 
			Limit: []byte("user-8"),
		}, nil) // 從user-3開始遍歷到user-7
	
	for iter.Next(){
		fmt.Printf("%v=%v\n", string(iter.Key()), string(iter.Value()))
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
	//遍歷指定前綴(開頭)的資料
	iter = db.NewIterator(leveldbUtil.BytesPrefix([]byte("user-")), nil)	// 篩選所有 key 以 "user-" 為開頭的資料項目。
	/*這邊的iterator是一個interface型別，這個 Iterator 是來自 github.com/syndtr/goleveldb/leveldb/iterator.Iterator，針對 LevelDB 優化設計，支援範圍查詢、前綴匹配、延遲加載等功能。
	type Iterator interface {
	Next() bool
    Prev() bool
    Seek(key []byte) bool
    First() bool
    Last() bool
    Key() []byte
    Value() []byte
    Error() error
    Release()
	}
	*/
	for iter.Next() {
		fmt.Printf("%v=%v\n", string(iter.Key()), string(iter.Value()))
	}
	iter.Release()
	err = iter.Error()
	if err != nil {
		fmt.Println(err)
	}
}

// 11.1.5~11.1.6 LevelDB Transaction and Snapshot
// Snapshot 是一種資料庫在某個時間點的唯讀快照，它允許你從該時間點開始讀取資料，即使資料庫之後有改變，Snapshot 讀到的內容仍然保持不變。
// Transaction（交易） 是一個原子性的操作集合，在 LevelDB 中透過 db.OpenTransaction() 開始，用來進行多筆資料的更新。
func LeveldbTransactionAndSnapshot() {
	db, err := leveldb.OpenFile("leveldb", nil)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	ss, err := db.GetSnapshot()
	if err != nil {
		panic(err)
	}
	defer ss.Release()	// Snapshot存在內存中
	t, err := db.OpenTransaction()
	if err != nil {
		panic(err)
	}
	batch := new(leveldb.Batch)
	for i := 1; i < 11; i++ {
		batch.Put(
			[]byte(fmt.Sprintf("cat-%v", i)),
			[]byte(fmt.Sprintf("{\"name\":\"c%v\"}", i)))
	}
	t.Write(batch, nil)
	// t.Discard()	// Discard 則會放棄整筆交易
	t.Commit()	// commit 才會寫入資料庫
	ok, _ := db.Has([]byte("cat-1"), nil)
	fmt.Println("db Has \"cat-1\" ?", ok)
	ok, _ = ss.Has([]byte("cat-1"), nil)
	fmt.Println("ss Has \"cat-1\" ?", ok)
}

// 11.2  Basic Operations of Redis
func RedisBasic() {
	// opt := redis.Options {
	// 	Addr: "localhost:6379",	// "6379"是安裝redis server預設的port
	// }
	// db := redis.NewClient(&opt)
	// context.Context 是從 redis v8 開始被強制要求的。v8 的所有 API 都強制帶入 ctx，這是 Go 開發趨勢，也符合 idiomatic Go 設計。
	// ctx := context.Background()	// 這是一個空白的 context，永不逾時、不會取消。context.Context 是 Go 語言設計用來在 API 間傳遞 deadline（截止時間）、cancel signal（取消信號）、request-scoped value（請求範圍變數） 的核心工具
	db := util.GetRedisClient()
	ctx := util.GetRedisContext()

	db.Do(ctx, "set", "k1", "v1")
	res, err := db.Do(ctx, "get", "k1").Result() // 正常會顯示"res= v1", 若將"k1"改為"k2"則為"該key不存在"
	if err != nil {
		if err == redis.Nil {
			fmt.Println("該key不存在")
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println("res=", res.(string))
	}

	// Basic type handling
	db.Do(ctx, "set", "b1", true)
	db.Do(ctx, "set", "b2", 0)
	// b, err := db.Do(ctx, "get", "b2").Bool()
	b, err := db.Do(ctx, "mget", "b1", "b2").BoolSlice()
	if err == nil {
		fmt.Println("b=", b)
	}

	db.Set(ctx, "t1", time.Now(), 0)
	// t, err := db.Get(ctx, "t1").Time()	// return一個time的類型
	t := db.Get(ctx, "t1").Val()	// 直接return string
	if err == nil {
		fmt.Println("t=", t)
	}
}

// 11.2.6 Redis Pipeline
func RedisPipeline() {
	db := util.GetRedisClient()
	ctx := util.GetRedisContext()

	pipe := db.Pipeline()	// 開啟一個管道
	t1 := pipe.Get(ctx, "t1")	// 這邊只是先設定好命令，但尚未執行，所以這邊回傳的參數t1也會是空的。回傳的t1為 *redis.StringCmd
	fmt.Println("pipe執行前的t1=", t1)
	for i := 0; i < 10; i++ {	// 用for loop來設定要批量執行的命令
		pipe.Set(ctx, fmt.Sprintf("p%v", i), i, 0)
	}
	_, err := pipe.Exec(ctx)	// 這邊才是真正去執行管道
	if err != nil {
		panic(err)
	}
	fmt.Println("pipe執行後的t1=", t1)

	cmds, err := db.Pipelined(ctx, func(pipe redis.Pipeliner) error{	// 回傳的cmds是一個[]cmder的切片，Cmder 是一個 interface（介面）
		for i:=0; i<10; i++ {
			pipe.Get(ctx, fmt.Sprintf("p%v", i))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	for i, cmd:= range cmds {
		fmt.Printf("p%v=%v\n", i, cmd.(*redis.StringCmd).Val())
		// 為什麼要 .(*redis.StringCmd)？->因為 Pipelined() 回傳的是通用介面 Cmder，你需要轉型才能呼叫對應方法(*redis.StringCmd（對應 GET 等會回傳字串的 Redis 指令） *redis.IntCmd（像是 INCR, DECR 等會回傳整數） *redis.StatusCmd *redis.BoolCmd 其他類型...)
		// 為什麼不能直接印 string？-> cmd 是介面不是字串，要呼叫 .Val() 才能取出具體值，且須先轉型
	}
}

// 11.2.6 Redis Transaction
func RedisTransaction() {
	db := util.GetRedisClient()
	ctx := util.GetRedisContext()
	// 以下是模擬一個「錢包轉帳」的邏輯->示範 Redis 的 Transaction：扣 p0 加 p1，具原子性，確保一致性
	// 把金額 100 從帳戶 p0 扣掉 
	// 把金額 100 增加到帳戶 p1 
	// 並保證這兩個動作是原子性 (atomic) 的：要嘛兩個都成功、要嘛兩個都失敗（不能只做一半）
	// 重點: 使用 WATCH 監控 key、用 pipeline 執行多指令、EXEC 提交交易

	for i:=0; i < 10; i++ {
		err := db.Watch(ctx, func(tx *redis.Tx) (err error) {	// 使用 Redis 的 WATCH + MULTI/EXEC 機制來進行樂觀鎖控制
			pipe := tx.Pipeline()	// 建立一個 transaction pipeline
			err = pipe.IncrBy(ctx, "p1", 100).Err()	// p1 加值 100
			if err != nil {
				return 
			}
			err = pipe.DecrBy(ctx, "p0", 100).Err()	// p0 減值 100
			if err != nil {
				return
			}
			_, err = pipe.Exec(ctx)	// 提交整個 transaction（如果中間有 key 被其他 client 改變，就會返回 TxFailedErr）
			return
		}, "p0")	// 也可寫成"}, "p0", "p1")"，這樣就會同時監視 p0 和 p1，當任一個 key 在交易前被其他 client 修改時，這次交易就會失敗。
		if err == nil {
			fmt.Println("Transaction commit成功")
			break
		} else if err == redis.TxFailedErr {	// 如果 transaction 失敗就重試（最多 10 次）
			fmt.Println("Transaction執行失敗, 這次是第", i, "次嘗試")
			continue
		} else {
			panic(err)
		}
	}
}

// 11.2.8 Redis Iteration
func RedisIterate() {
	db := util.GetRedisClient()
	ctx := util.GetRedisContext()

	// Scan
	iter := db.Scan(ctx, 0, "p*", 0).Iterator()
	// 第一個參數 0 是 cursor（游標），初始為 0，會自動由 Iterator 控制。
	// 第二個參數 "p*" 是 pattern，這邊代表匹配前綴為 "p" 的所有鍵（例如 p0, p1…）
	// 第三個參數 0 是 count，表示讓 Redis 自行決定回傳的筆數（可視為優化的 hint）
	for iter.Next(ctx) {
		fmt.Printf("key=%v, value=%v\n", iter.Val(), db.Get(ctx, iter.Val()).Val())
	}

	// if 語句中的短變數宣告（short variable declaration）。可以在 if 語句裡面進行短變數宣告後，用邏輯運算子 && 或 || 結合其他條件
	// 短變數宣告只能宣告一個變數，不能同時宣告多個。
	// 好處: 
	// 限定作用域：err 僅在這個 if 區塊內有效，不會污染整個函式或其他區塊的變數命名空間。
	// 更簡潔：常見於 error handling 的寫法，Go 社群也鼓勵這種方式。
	if err := iter.Err(); err != nil {	
	}

	// HScan
	// db.Set / db.Scan 跟 db.HSet / db.HScan 的差異其實反映了 Redis 不同資料結構（data type）在使用與底層儲存上的根本設計差異
	// Set, Get, Scan String（最基本型別） 儲存單一 key 對應的 value（字串）
	// HSet, HGet, HScan Hash（類似 dict/map） 一個 key 對應多個 field:value 欄位的集合 

	db.HSet(ctx, "h1", "f1", "v1", "f2", "v2", "f3", "v3")	// db.HSet(...): 新增一個 hash 鍵 "h1"，包含 f1:v1, f2:v2, f3:v3 三個欄位。
	iter = db.HScan(ctx, "h1", 0, "*", 0).Iterator()	// db.HScan(...)：針對 hash 鍵 "h1" 的欄位做遍歷（非阻塞）
	for i:=0 ; iter.Next(ctx); i++ {
		if i % 2 ==0 {	// iter.Val()：會依序回傳 field、value、field、value，所以要用 index i 判斷偶數是 field、奇數是 value。
			fmt.Printf("field=%v, ", iter.Val())
		}else {
			fmt.Printf("field=%v\n", iter.Val())
		}
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

}

// 11.2.9 Scan Redis Hash to Structure
type RedisHash struct {	// 這邊的開頭要大寫for public，因為要傳給redis包(exported，跨套件存取)
	Name string	`redis:"name"`
	Id int `redis:"id"`
	Online bool `redis:"online"`
}

// 把一個 RedisHash 結構寫入 Redis 的 hash 結構。
// 再從 Redis 把該 hash 資料取出來，轉成 RedisHash 結構
func RedisHashToStruct() {
	db := util.GetRedisClient()
	ctx := util.GetRedisContext()
	var rh1 = RedisHash{
		Name: "rhName",
		Id: 123,
		Online: true,
	}

	db.Pipelined(ctx, func(pipe redis.Pipeliner) error {	// 使用 Redis 的 pipeline 技術，一次發送多個 HSET 命令來提高效能
		pipe.HSet(ctx, "rh1", "name", rh1.Name)
		pipe.HSet(ctx, "rh1", "id", rh1.Id)
		pipe.HSet(ctx, "rh1", "online", rh1.Online)
		return nil	// 如果想要回報一個 Redis 相關錯誤（如 Key 或 Field 不存在等），才會 return redis.Nil
	})

	var rh2 RedisHash
	// HGetAll(): 一次抓取整個 Redis hash 中的 所有 field-value 對，回傳 map[string]string。適合用來讀取整個物件，像你這段 Scan(&rh2) 就是搭配整筆資料轉 struct
	// HGet():只抓取指定 field 的值，例如：HGet(ctx, "rh1", "name") 只會回傳 "rhName"。 回傳值是單一 StringCmd（不是 map）
	if err := db.HGetAll(ctx, "rh1").Scan(&rh2); err != nil{	// 傳入指標（pointer）&rh2 才能讓函式真的去改變外部變數的內容，如果你寫成 Scan(rh2)，那只是一份 copy，無法把值寫入 rh2 這個變數本身。結論：Scan() 一定要傳指標（如 &yourStruct）
		panic(err)
	}	// HGetAll返回的是一個map[string]string(redis.StringStringMapCmd)
	fmt.Printf("rh2=%+v\n", rh2)	// "%+v":印出變數內容，並包含struct欄位名稱，Ex.{Name:test Id:1 Online:true}。"%#v": 印出該值的Go語法表示形式(含型別)，Ex.main.RedisHash{Name:"test", Id:1, Online:true}
}