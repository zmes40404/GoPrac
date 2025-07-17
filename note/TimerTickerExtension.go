package note

import (
	"fmt"
	"time"
)

func heartbeatMonitor(stop chan bool) {	// 接收一個 channel，當你要結束這個心跳監控時，就對這個 channel 寫入一個訊號。
	ticker := time.NewTicker(5 * time.Second) // 建立一個 Ticker (主要目的: 多次觸發)，每 5 秒會觸發一次。 Ticker 底下有個 .C channel，會定期送出「時間訊號」。
	defer ticker.Stop() // 在函式結束時記得呼叫 Stop()，釋放 ticker 的資源，避免 memory leak。

	for {	// 開始一個無窮迴圈，select 是 Go 的「多 channel 監聽語法」。
		select {	
		case <-ticker.C: // 每 5 秒會執行一次，印出心跳訊息。
				fmt.Println("[Heartbeat] Service is alive...")
		case <-stop: // 當 stop channel 收到訊號，就結束函式（回傳），等於停止這個 goroutine。
			fmt.Println("Stopping heartbeat monitor...")
			return // Exits the function when a stop signal is received
		}
	}
}

func MainHeartbeatMoniter() {
	// 這段程式碼模擬一個「心跳監控系統」，每 5 秒印出一次 [Heartbeat] Service is alive... 表示服務還活著。主程式會在 12 秒後，透過 stopChan 通知監控 goroutine 結束。
	stopChan := make(chan bool) // 建立一個「停止訊號通道」。
	go heartbeatMonitor(stopChan) // 啟動一個 goroutine 執行 heartbeatMonitor，非同步執行心跳邏輯。

	time.Sleep(12 * time.Second) // 主程式等 12 秒，模擬「服務執行一段時間」的情境。
	stopChan <- true // 傳送一個訊號到 stopChan，觸發 heartbeatMonitor() 中的 <-stop case → 結束監控 goroutine。

	// Output:
	/*
	[Heartbeat] Service is alive...
	[Heartbeat] Service is alive...
	Stopping heartbeat monitor.
	（5 秒印一次心跳 → 第三次還沒印就被主程式終止）
	*/
}

func TimerTimeoutControl() {
	done := make(chan bool) // 建立一個 done channel，當操作完成時會發送訊號。

	// 模擬背景工作，3秒完成
	go func() {
		time.Sleep(3 * time.Second) // 模擬一個耗時操作，這裡假設是 3 秒。
		done <- true // 當工作完成時，發送訊號到 done channel。
	}()

	timer := time.NewTimer(5 * time.Second) // 建立一個 5 秒的計時器。

	select {
		case <- done: // 如果 done channel 收到訊號，表示工作完成。
			fmt.Println("工作完成!")
		case <- timer.C: // 如果計時器到期，表示工作超時。
			fmt.Println("等太久了，Timeout!")
	}

	/*
	背景 goroutine 用 Sleep(3秒) 模擬耗時任務 
	select 同時監聽 done（工作完成）與 timer.C（5秒 timeout） 
	因為任務 3 秒就完成 → 執行 case <-done:，輸出： 工作完成！
	如果改成6秒，則會執行 case <-timer.C:，輸出： 等太久了，Timeout!
	*/
}