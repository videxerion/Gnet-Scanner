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
			setBoolCommonVar(&exitState, &exitMu, true)
			break
		} else if key == keyboard.KeyCtrlE {
			os.Exit(0)
		} else if char == 'p' || char == 'P' {
			setBoolCommonVar(&pauseState, &pauseMu, true)
		} else if char == 'r' || char == 'R' {
			setBoolCommonVar(&pauseState, &pauseMu, false)
		}
	}

}
