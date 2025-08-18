package main

import (
	"goprac/note"
	// "goprac/note/factory"
)


func main() {
	note.LeveldbIterate()

	// note.MainHeartbeatMoniter() // function to demonstrate heartbeat monitoring using ticker in "TimerTickerExtension.go"
	// note.TimerTimeoutControl() // function to demonstrate timer timeout control in "TimerTickerExtension.go"

	// m := &factory.mes{}
	// m := factory.NewMes()
	// m.SetPwd("")
	
}