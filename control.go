package main

import (
	"github.com/eiannone/keyboard"
	"os"
)

func interceptingKeystrokes() {
	err := keyboard.Open()
	if err != nil {
		panic(err)
	}
	defer keyboard.Close()

	for {
		char, key, _ := keyboard.GetKey()

		if key == keyboard.KeyCtrlC {
			os.Exit(0)
		} else if char == 'p' || char == 'P' {
			pauseMu.Lock()
			pauseState = true
			pauseMu.Unlock()
		} else if char == 'r' || char == 'R' {
			pauseMu.Lock()
			pauseState = false
			pauseMu.Unlock()
		}
	}

}
