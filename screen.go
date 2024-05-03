package main

import (
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

var clearScreen map[string]func()

func init() {
	clearScreen = make(map[string]func())
	clearScreen["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clearScreen["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func CallClear() {
	value, ok := clearScreen[runtime.GOOS]
	if ok {
		value()
	} else {
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}

func speedMeter() {
	var oldCount uint64
	var newCount uint64

	for !exitState {
		oldCount = scannedAddress
		time.Sleep(time.Second)
		newCount = scannedAddress

		speed = newCount - oldCount
	}

}

func infoScreen() {
	for !exitState {
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
		println("Count Threads:\t ", countThreads)
		println(
			"Progress:\t ",
			strconv.FormatUint(scannedAddress, 10)+"/"+strconv.FormatUint(countAddress, 10),
			strconv.FormatFloat(percentageOfNumber(float64(scannedAddress), float64(countAddress)), 'f', 2, 64)+"%",
		)
		println("Speed:\t\t ", speed, "ip/s")
		println("Time left:\t ", timeLeftString)
		println("Allocated memory:", allocatedMemoryString)

		if !pauseState {
			println("\nControl:", "[p]pause")
		} else {
			println("\nControl: [r]resume")
		}
	}
}
