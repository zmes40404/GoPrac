package note

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"goprac/util"
	"runtime"

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
