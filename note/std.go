package note

import (
	"fmt"
	"math/rand" // generated random numbers are predictable easily, if want more secure random numbers, use "crypto/rand" package or reference about "GO cryptography"
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