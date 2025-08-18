package note

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
	leveldbUtil "github.com/syndtr/goleveldb/leveldb/util"
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
		&leveldbUtil.Range{Start: []byte("user-3"), Limit: []byte("user-8")}, nil) // 從user-3開始遍歷到user-7
	
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