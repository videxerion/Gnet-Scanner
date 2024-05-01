package main

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func speedMeter() {
	var oldCount int
	var newCount int

	for {
		oldCount = scannedAddress
		time.Sleep(time.Second)
		newCount = scannedAddress

		speed = newCount - oldCount
	}

}

func infoScreen() {
	for {
		time.Sleep(time.Second / 4)

		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		var timeLeftString = "∞"
		var allocatedMemoryString = "None"
		if speed != 0 {
			secondsLeft := (countAddress - scannedAddress) / speed
			timeLeftString = convertSecondsToTime(secondsLeft).ToString()
		}
		allocatedMemory := m.Alloc
		if allocatedMemory != 0 {
			allocatedMemoryString = convertBytesToSize(int(allocatedMemory)).ToString()
		}

		CallClear()

		println("Count DB Threads:", countDbThreads)
		println("Count Threads:", countThreads)
		println(
			"Progress:",
			strconv.Itoa(scannedAddress)+"/"+strconv.Itoa(countAddress),
			strconv.FormatFloat(percentageOfNumber(float64(scannedAddress), float64(countAddress)), 'f', 2, 64)+"%",
		)
		println("Speed:", speed, "ip/s")
		println("Time left:", timeLeftString)
		println("Allocated memory:", allocatedMemoryString)

		if !pauseState {
			println("\nControl:", "[p]pause")
		} else {
			println("\nControl: [r]resume")
		}
	}
}
