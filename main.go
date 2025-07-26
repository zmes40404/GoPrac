package main

import (
	"goprac/note"
)


func main() {
	note.PackageSync()
	// note.MainHeartbeatMoniter() // function to demonstrate heartbeat monitoring using ticker in "TimerTickerExtension.go"
	// note.TimerTimeoutControl() // function to demonstrate timer timeout control in "TimerTickerExtension.go"
}