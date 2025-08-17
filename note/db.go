package note

import (
	"fmt"

	"github.com/syndtr/goleveldb/leveldb"
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
 