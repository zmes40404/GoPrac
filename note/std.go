package note

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"goprac/util"
	"runtime"
	"sort"
	"sync"

	// "goprac/util"
	"math/rand" // generated random numbers are predictable easily, if want more secure random numbers, use "crypto/rand" package or reference about "GO cryptography"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// 6.1 ramdon number
func RandomNumber() {	// 跨包調用函式開頭需要大寫
	// fmt.Println(rand.Intn(10)) // Generates a random number between 0 and 9 (0~n-1)
	
	// r := rand.New(rand.NewSource(1234))	// Creates a new random number generator with a specific seed (fixed seed: 1234 in this case)
	// fmt.Println(r.Intn(10)) // seed is fixed, so it will always generate the same sequence of random numbers.

	for i:=0; i<10; i++ {
		rand.Seed(time.Now().UnixNano()) // "time.Now().UnixNano()" returns the current time in nanoseconds -> used to seed the random number generator, "rand.Seed" is deprecated: As of Go 1.20 there is no reason to call Seed with
										// a random value.
		fmt.Println(rand.Intn(10) - 9) // Generates a random number between -9 and 0
	}
}

//6.2 String type conversion
func StrConv() {
	num1 := 123
	s1 := "MatonWang"
	s2 := fmt.Sprintf("%d@%s", num1, s1) // Converts num1 to a string and concatenates it with S1
	fmt.Println("S2=", s2)

	var (
		i2 int
		s3 string
	)

	n, err := fmt.Sscanf(s2, "%d@%s", &i2, &s3)
	if err != nil {
		panic(err) // If there is an error during scanning, it will panic
	}
	fmt.Printf("Parse %d nums of data successfully \n", n) // Prints the number of items scanned
	fmt.Println("i2=", i2, "s3=", s3) // Prints the parsed values of i2 and s3	

	s4 := strconv.FormatInt(123, 4)
	fmt.Println("s4=", s4) // Converts the integer 123 to a string in base 4 representation
	u1, err := strconv.ParseUint(s4, 4, 0) // Parses the string s4 as an unsigned integer in base 4
	if err != nil {
		panic(err) // If there is an error during scanning, it will panic
	}
	fmt.Println("u1=", u1) 
}

//6.3 Common functions in string package
func PackageStr() {
	fmt.Println(strings.Contains("hello", "ll"))
	fmt.Println(strings.Count("hello", "ll")) // Counts the number of occurrences of "l" in "hello"
	fmt.Println(strings.Replace("hello", "l", "dd", -1)) // Replaces all occurrences of "l" with "dd" in "hello"
	fmt.Println(strings.Repeat("mia", 3)) // Repeats the string "mia" three times
	fmt.Println(strings.Fields("mia mia\n mia\tmia")) // Splits the string into fields based on whitespace characters (spaces, tabs, newlines)
	fmt.Println(strings.Trim("#*\naaa.aaa.aaa&%#", "#*\n&%")) // Trims the specified characters from both ends of the string
}

//6.4 "utf8" package common functions
func PackageUtf8() {
	// fmt.Println(utf8.RuneCountInString("hello,world")) // Counts the number of runes (characters) in the string "hello,world"
	str := "hello,世界"
	fmt.Println(utf8.ValidString(str)) // Checks if the string is a valid UTF-8 encoded string -> true
	fmt.Println(utf8.ValidString(str[:len(str)-1])) // false: 因為中文通常都占用3個 bytes，所以這裡的str[:len(str)-1]會切掉最後一個字元，導致不完整的UTF-8編碼
}

//6.5 "time" package common functions
func PackageTime(){
	fmt.Println("\n 6.5 time section")
	for range 5 {
		fmt.Print(".")
		time.Sleep(time.Microsecond * 100) // Sleeps for 100 microseconds
	}
	fmt.Println()

	d1, err:= time.ParseDuration("1000s")
	if err != nil {
		panic(err) // If there is an error during parsing, it will panic
	}
	fmt.Println("d1=", d1) // Prints the parsed duration of 1000 seconds

	fmt.Println("\n 6.5.2 time zone")
	l1, err := time.LoadLocation("Asia/Taipei") // Loads the time zone for Asia/Taipei
	if err != nil {
		panic(err) 
	}
	fmt.Println(l1.String()) // Prints the string representation of the loaded time zone

	fmt.Println("\n 6.5.5 time")
	fmt.Println(time.Now().Format("2006年1月2日, 15點04分")) // time.Format 的格式字串必須用「2006-01-02 15:04:05」這組數字（Go 的基準時間）來代表各個欄位。但 Go 只認得 15 代表小時、04 代表分鐘。
	t2, err := time.ParseInLocation("2006年1月2日, 15點04分", "2100年12月31日, 17點14分", l1) // Parses the date string in the specified format and time zone
	if err != nil {
		panic(err) 
	}
	fmt.Println(t2)
	fmt.Println(t2.Location()) // Prints the location of the parsed time

	fmt.Println("\n 6.5.6 ticker")
	intChan := make(chan int, 1) // Creates a buffered channel to receive integers
	go func() {
		time.Sleep(time.Second)
		intChan<-1 // Sends a value to the intChan channel after 1 second
	}()
	TickerFor:
	for {
		select {
			case <-intChan:
				fmt.Println()
				break TickerFor // Breaks the loop if a value is received from the intChan channel
			case <-time.NewTicker(100 * time.Millisecond).C: // Creates a ticker that ticks every 100 milliseconds
				fmt.Print(".") // 每次 select 進入時都會重新創建一個新的 Ticker，但你沒有 Stop() 它-> 每 100 毫秒建立一個新的 goroutine(時間越短會印越多...)，每個 goroutine 都在 tick，結果：你印出 . 的速度越來越快 → 變成一堆點點...
		}
	}
	
	fmt.Println("\n 6.5.7 timer")
	select {
	case <-intChan:	// 這邊的intChan沒有接收任何資料，所以會block在這裡，導致只會執行下面那個case的情況
		fmt.Println("收到使用者發送的驗證碼")
	case <-time.NewTimer(time.Second).C:
		fmt.Println("驗證碼已過期")
	}
}

// Common file operations
func FileOperation() {
	// util.MKdirWithFilePath("d1/d2/fil2")
	fmt.Println("\n 6.6.5 Folder Operation")
	dirEntrys, err := os.ReadDir("C:/GoPrac/util") // dirEntrys is a slice of os.DirEntry, which represents the entries in the specified directory
	if err != nil {
		panic(err)
	}
	for _, v := range dirEntrys {
		fmt.Println(v.Name())
	}
	 

	fmt.Println("\n 6.6.6 File Operation")
	file, err := os.OpenFile("f1", os.O_RDWR|os.O_CREATE, 0665)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	fmt.Println("\n 6.6.7 Unbuffered reading and writing (suitable for small files)")
	data, err := os.ReadFile("f1") // Reads the entire content of the file "f1" into memory
	if err != nil {
		panic(err)
	}
	fmt.Println("f1中的資料為", string(data)) // Converts the byte slice to a string and prints it
	err = os.WriteFile("f2", data, 0775)
	if err != nil {
		panic(err)
	}
}

// 6.7 Files Reading and Writing
func FileReadAndWrite() {
	f5, err := os.OpenFile("f5", os.O_WRONLY|os.O_CREATE, 0666) // Opens the file "f5" for writing, creating it if it doesn't exist
	if err != nil {
		panic(err)
	}
	defer f5.Close() // Ensures the file is closed when the function exits
	writer := bufio.NewWriter(f5) // Creates a buffered writer for the file
	fmt.Println("buffer size: ", writer.Size()) 
	for i := range 5 {	// 假如有f1~f4四個文件，這裡會讀取f1~f4的內容並寫入到f5中
		fileName := fmt.Sprintf("f%v", i)
		data, err := os.ReadFile(fileName) // Reads the content of the file "f%v" into memory
		if err != nil {
		panic(err)
		}
		data=append(data, '\n') // Appends a newline character to the data
		writer.Write(data) // 寫入記憶體RAM
	}
	writer.Flush() // 確保所有記憶體的資料都被寫入到硬碟(f5這個文件)中
} 

// 6.8 Error
func Errors(){
	defer func() {
		err := recover() // recover() 用來捕捉 panic 的錯誤
		fmt.Println("捕捉到了錯誤:", err) 
	}()
	err1 := errors.New("可愛的錯誤") // Creates a new error with the message "可愛的錯誤"
	fmt.Println("err1=", err1)
	err2 := fmt.Errorf("%s的錯誤", "溫柔")
	fmt.Println("err2=", err2) // Creates a formatted error with the message "溫柔的錯誤"
	panic(err1) // Panics with the error message, which will stop the program execution
}

// 6.9 Logging
func Log() {
	defer func() {
		err := recover() // recover() 用來捕捉 panic 的錯誤
		fmt.Println("捕捉到了錯誤:", err)
	}()
	err := errors.New("這是一個錯誤")
	util.INFO.Println(err)
	util.DEBUG.Println(err)
	// util.WARN.Panicln(err)
	// util.ERROR.Fatalln(err) // This will log the error and then call os.Exit(1), terminating the program, 0代表正常退出，1代表異常退出。Fatal只能退出，連defer recover都不會執行。
}

// 6.10 Unit Test
func IsNotNegative(n int) bool {
	return n > -1
}

// 6.11 Command Prompt Arguments
func CmdArgs() {
	fmt.Printf("接受到了%v個參數\n", len(os.Args))
	for i, v := range os.Args {
		fmt.Printf("第%v個參數是%v\n", i, v)
	}
	fmt.Println()

	/*
	Input: {go run . fsed -fdv sdc -sdf=sv "csf dss "}
	Output: 第0個參數是C:\Users\Maton_Wang\AppData\Local\go-build\02\02b1f4d14ad032f98f0cf97257051459bbe434718b07811823a026922748fbf4-d\goprac.exe -> 當前可執行文件的位置(緩存位置)
			第1個參數是fsed
			第2個參數是-fdv
			第3個參數是sdc
			第4個參數是-sdf=sv
			第5個參數是csf dss
	*/

	vPtr := flag.Bool("v", false, "Go版本號") // flag.Bool(...) 的回傳值是*bool，這是一個指向 bool 的指標，指向的值會被 flag 包裝起來。代表要存取真正的布林值，要透過解引用(dereferencing): *vPtr，才能拿到布林值的實際內容(true or false) 
	var userName string
	flag.StringVar(&userName, "u", "", "用戶名")
	flag.Func("f", "", func(s string) error {
		fmt.Println("s=", s) // Output: s= 444
		return nil // Returns nil to indicate no error occurred
	}) // flag.Func(...) allows you to define a custom function to handle the flag, here it just prints the value of the flag	
	flag.Parse() // Parses the command-line flags and arguments
	if *vPtr { // If the -v flag is set, it will print the Go version	
		fmt.Println("Go版本是 V0.0.0")
	}
	fmt.Println("當前用戶為:", userName)
	for i, v := range flag.Args() { // flag.Args() returns the non-flag command-line arguments
		fmt.Printf("第%v個無flag參數是%v\n", i, v)
	}

	/*
	Input: go run . -u fang -v f s sfsdf
	Output: 
		第0個無flag參數是f
		第1個無flag參數是s
		第2個無flag參數是sfsdf
	*/

}

// 6.12 "builtin" package
func PackageBuiltin() {
	c1 := complex(12.34, 45.67) // Creates a complex number with real part 12.34 and imaginary part 45.67
	fmt.Println("c1=", c1)
	r1 := real(c1) // Gets the real part of the complex number
	i1 := imag(c1) // Gets the imaginary part of the complex number
	fmt.Println("r1=", r1, "i1=", i1) // Prints the real and imaginary parts of the complex number
}

// 6.13 runtime package
func PackageRuntime() {
	cpuNum := runtime.NumCPU() // Gets the number of logical CPUs usable by the current process
	fmt.Println("當前系統有", cpuNum, "個 CPU")
	if runtime.NumCPU() > 7 {
		oriNum := runtime.GOMAXPROCS(runtime.NumCPU() - 1)	 // Sets the maximum number of CPUs that can be executing simultaneously to one less than the number of available CPUs
		fmt.Println("原本的最大 CPU 數量是", oriNum)
	}
	// runtime.Goexit() // 這一定會報錯，這行code導致沒有一個線程能繼續執行，因為它會終止當前的 goroutine。這個函式通常用於在 goroutine 中提前退出，而不會影響其他 goroutine 的執行。
}

// 6.14 sync package
func PackageSync() {
	// 計算從 2 到 100000 之間有多少個質數（prime number），並使用 Goroutine 並行處理提升效能。
	// 但因為多個 Goroutine 都會對共享變數 c 做加一的動作（c++），所以 必須使用 mutex 保護臨界區（critical section），以避免數據競爭（race condition）。
	fmt.Println("\n 6.14.1 Mutex互斥鎖 / 6.14.2 WaitGroup") // "WaitGroup": 通過計數器來獲得阻塞能力
	var c int
	var mutex sync.Mutex // Mutex（互斥鎖）保護 c++ 的區塊，保證同一時間只有一個 Goroutine 在修改 c。
	var wg sync.WaitGroup // 在 Go 中，一旦 main() 執行完畢，整個程式就會結束，不管還有多少 Goroutine 沒執行完。

	primeNum:=func (n int) {
		defer wg.Done()	// Goroutine 執行完會通知「我結束了」
		for i:=2; i<n; i++ {
			if n%i==0 {
				return 
			}
		}
		// 如果是質數，執行 mutex.Lock() → c++ → mutex.Unlock()。
		mutex.Lock()
		c++
		mutex.Unlock()
	}
	
	for i:= 2; i<100001; i++ {	// 共產生 99999 個 Goroutine 同時判斷質數。 並行處理，大幅提升效能。
		wg.Add(1)	// 告訴主程式「我有一個 Goroutine 開始了」
		go primeNum(i)
	}
	// 沒有 WaitGroup 的話會怎樣？-> 程式可能在還沒算完質數時就印出 0 或提早結束，因為主線程直接跑到 fmt.Println(...)。
	wg.Wait()	// 主程式等到所有 Goroutine 都 Done() 後才會繼續執行
	fmt.Printf("\n總共找到%v的質數\n", c)

	fmt.Println("\n 6.14.3 Cond")	// Cond提供了同時控制多個thread阻塞的能力
	cond := sync.NewCond(&mutex)
	for i:=0; i<10; i++ {
		go func (n int)  {
			cond.L.Lock()	// 保護與「條件」相關的 共享狀態（像 ready、任務佇列、緩衝區資料 等等）。
			cond.Wait()	// 阻塞，等待條件被喚醒（釋放鎖，之後再重新取得）
			fmt.Printf("協程%v被喚醒了\n", n)
			cond.L.Unlock()
		}(i)
	}
	for i:=0; i < 15; i++ {
		time.Sleep(100 * time.Millisecond)
		fmt.Print(".")
		if i==4 {
			fmt.Println()
			cond.Signal()	// 喚醒其中一個在 Wait() 的 Goroutine
		}
		if i==9 {
			fmt.Println()
			cond.Broadcast()	// 喚醒所有等待中的 Goroutine
		}
	}

	fmt.Println("\n 6.14.4 Once") // "Once"確保一個函數只能被執行一次
	var once sync.Once
	for range 10 {
		wg.Add(1)
		go func(){
			once.Do(func(){
				fmt.Println("只有一次機會")
			})
			wg.Done()
		}()
	}
	wg.Wait()

	fmt.Println("\n 6.14.5 Map") // "Map": 併發安全的鍵值對
	var m sync.Map
	m.Store(1, 100)
	m.Store(2, 200)
	m.Store(3, 300)
	m.Range(func(key, value interface{}) bool {
		fmt.Printf("m[%v]=%v\n", key, value.(int))
		return true
	})
}

// 7.4 Sort Package
type Person struct {
	Name string
	Age int
}

// 實現出3種sort interface: Len() int, Less(i, j int)bool, Swap(i, j int)
type PersonSlice []Person
func (ps PersonSlice) Len() int {
	return len(ps)
}
func (ps PersonSlice) Less(i, j int) bool {
	return ps[i].Age > ps[j].Age // 降序排序條件
}
func (ps PersonSlice) Swap(i ,j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func PackageSort() {
	fmt.Println("\n 7.4.1 Sorting common types")
	intSlice:=[]int{2, 4, 8, 10}
	v:=6
	i:=sort.SearchInts(intSlice, v)
	fmt.Printf("%v適合插入在%v的%v\n", v, intSlice, i)

	fmt.Println("\n 7.4.2 自訂義排序")
	p:=[]Person{{"小小", 18}, {"小方", 5}, {"小米", 50}}
	
	// func Slice(slice any, less func(i, j int) bool): 它使用你提供的 less 函數（兩個 index 比大小），在內部做 quicksort
	sort.Slice(p, func(i, j int) bool {	//  i 和 j 是 slice 中的索引值
		return p[i].Age < p[j].Age	// 該函數必須回傳一個 bool: 若回傳 true，代表 p[i] 要排在 p[j] 前面。若回傳 false，代表 p[i] 要排在 p[j] 後面。排序過程呼叫: p[0].Age < p[1].Age → 18 < 5 → false, p[1].Age < p[2].Age → 5 < 50 → true
	})
	fmt.Println("p=", p)

	fmt.Println("\n 7.4.3 自訂義查找")
	i = sort.Search(len(intSlice), func(i int) bool {
		return intSlice[i] >= v
	})
	fmt.Printf("%v中第一次出現不小於%v的位置是%v\n", intSlice, v, i)

	fmt.Println("\n 7.4.4 sort.interface")
	// sort.Sort(PersonSlice(p)) // 就邊Sort()出來的結果會使用前面定義的Less() interface，從大到小排列
	sort.Sort(sort.Reverse(PersonSlice(p)))
	fmt.Println("p=", p)
}

// 9.1 Common Operation on "JSON"
func PackageJSON() {
	// bool -> 將encode成JSON booleans
	// float, int -> 將encode為JSON number
	// string -> JSON string
	// slice, array -> JSON array
	// struct, map -> JSON object
	// nil -> JSON null

	type user struct {
		Name string		`json:"name"`
		Age int			`json:"age,omitempty"`
		Email string	`json:"-"`	// 想用"-"當作別名，在後面加個逗號即可，Ex"-,"
		Job map[string]string
	}

	u1:=user{
		Name: "Diamond",
		Age: 3,
		Email: "123@asc.com",
		Job: map[string]string{
			"早班": "警衛",
			"午班": "洗碗工",
			"晚班": "外送員",
		},
	}
	data, _ := json.Marshal(u1) // 在struct的情況下，直接傳address->(&u1)效果是一樣的，但是傳指標的整體效能更好，因為"json.Marshal(u1)"拷貝一份資料 → 傳入函式，通常效能稍低（尤其 struct 很大時）
	// 所以只有在「具體型別」的 struct 時可以兩種都用，interface 的情況下只能傳值（不能取指標）。
	fmt.Println(string(data))	// 這邊印出來是沒有分行的

	// 格式化
	buf:= new(bytes.Buffer)
	json.Indent(buf, data, "", "\t")
	fmt.Println(buf.String())

	var u2 user
	json.Unmarshal(data, &u2)
	fmt.Println("u2=", u2)
}