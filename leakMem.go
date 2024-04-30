package main

import (
	"fmt"
	"github.com/gen2brain/beeep"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

func prof() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
}

func downloadDump() {
	// URL, с которого нужно скачать файл
	url := "http://127.0.0.1:6060/debug/pprof/heap"

	// Создаем HTTP-запрос
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		panic(err)
	}

	// Создаем файл для записи
	outFile, err := os.Create("heap")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	// Копируем содержимое ответа в файл
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		panic(err)
	}
}

func leakDetector() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		if m.Alloc >= GigaByte {
			_ = beeep.Alert("MEMORY LEAK DETECTED", fmt.Sprintf("Program allocated %d MegaBytes", int(m.Alloc/GigaByte)), "")
			downloadDump()
			return
		}
		time.Sleep(time.Millisecond * 100)
	}
}
