package main

import (
	"net"
	"sync"
	"time"
)
import _ "net/http/pprof"

// Создаём мютексы
var (
	countThreadsMu   sync.Mutex
	scannedAddressMu sync.Mutex
	pauseMu          sync.Mutex
	dbWriteMu        sync.Mutex
	ctDbThMu         sync.Mutex
	exitMu           sync.Mutex
)

// Переменные для показателей сканированния
var (
	countDbThreads uint64 = 0
	countThreads   uint64 = 0
	countAddress   uint64 = 0
	scannedAddress uint64 = 0
	speed          uint64 = 0
)

// Инициализируем перменные под флаги
var (
	// Функциональные флаги
	debugFlag   bool
	inputNet    string
	saveFlag    string
	enableSaves bool

	// Флаги параметров
	chunkSize      uint64
	limitThreads   uint64
	connectTimeout time.Duration
	readTimeout    time.Duration
	responseSize   uint64
	pathToDb       string
	pathToSaves    string
)

// Состояния для управления
var (
	pauseState = false
	exitState  = false
)

// Константы
const fileSaveVersion uint64 = 1

var usrSave save
var usrDatabase *Db

func main() {
	var err error

	// Получение флагов
	getFlags()

	// Если указан путь к файлу сохранения то пытаемся его распарсить
	if saveFlag != "None" {
		usrSave = parseSaveFile(saveFlag)
		scannedAddress = usrSave.scannedAddr
	}

	// Проверка корректности полученных флагов
	checkValidFlags()

	// Если включён режим отладки то запускаем pprof на localhost:6060 и детектор утечек
	if debugFlag {
		go prof()
		go leakDetector()
	}

	// Парсим полученный адрес сети
	var ipnet *net.IPNet
	if saveFlag == "None" {
		_, ipnet, err = net.ParseCIDR(inputNet)
	} else {
		_, ipnet, err = net.ParseCIDR(usrSave.inputNet)

	}
	if err != nil {
		panic(err)
	}

	// Создаём и подключаем базу данных
	if saveFlag == "None" {
		usrDatabase = Database(pathToDb + time.Now().Format(time.DateTime) + " result.db")
	} else {
		usrDatabase = Database(usrSave.dbName)
	}
	usrDatabase.createTable()

	// Считаем кол-во адрессов в указанной сети
	if saveFlag == "None" {
		countAddress, err = countIPAddresses(inputNet)
	} else {
		countAddress, err = countIPAddresses(usrSave.inputNet)
	}
	if err != nil {
		panic(err)
	}

	// Запуск служебных потоков
	go speedMeter()
	go infoScreen()
	go interceptingKeystrokes()
	if enableSaves {
		if saveFlag == "None" {
			go saveThread(inputNet, usrDatabase.name)
		} else {
			go saveThread(usrSave.inputNet, usrDatabase.name)
		}
	}

	// Вычисляем IP с которого нужно начинать сканирование
	var startIP net.IP
	if saveFlag == "None" {
		startIP = ipnet.IP.Mask(ipnet.Mask)
	} else {
		startIPint := ipToUint64(ipnet.IP.Mask(ipnet.Mask))
		newIPInt := startIPint + scannedAddress
		startIP = uint64ToIP(newIPInt)
	}

	chunk := make([]string, chunkSize)
	var pointer uint64 = 0

	// Начинаем перебор
	for ip := startIP; ipnet.Contains(ip); nextIP(ip) {

		// Если кончились IP адреса, то выходим из цикла
		if ip == nil {
			// Отправляем последний чанк на сканирование
			sendChunkForScanning(chunk)
			break
		}

		// Ожидаем завершение паузы если была поставлена
		if pauseState {
			waitPauseEnd()
		}

		// Если IP действителен то записываем в чанк, иначе пропускаем
		if isValidIP(ip) {
			chunk[pointer] = ip.String()
		} else {
			incCommonVar(&scannedAddress, &scannedAddressMu)
			continue
		}

		// Если чанк заполнен до конца, то отправляем на сканирование
		if pointer == chunkSize-1 {
			sendChunkForScanning(chunk)
			clearChunk(&chunk)

			pointer = 0
		} else {
			pointer++
		}

	}

	// Ждём завершения потоков
	waitCompletionThreads()
}

// Функция отправляет чанк на сканирование
func sendChunkForScanning(chunk []string) {
	waitForThreadToFree()

	incCommonVar(&countThreads, &countThreadsMu)

	chunkCopy := make([]string, chunkSize)
	copy(chunkCopy, chunk)

	go scanChunk(chunkCopy, usrDatabase)

}

// Функция ожидаает пока не освободится место для нового потока
func waitForThreadToFree() {
	for {
		if countThreads < limitThreads {
			break
		}
	}
}

// Функция ожидает пока не завершатся все потоки
func waitCompletionThreads() {
	for {
		if countThreads == 0 {
			break
		}
	}
}

// Функция ожидает пока пауза не будет снята
func waitPauseEnd() {
	for {
		if !pauseState {
			break
		}
	}
}

// Функция очищает чанк
func clearChunk(chunkPointer *[]string) {
	chunk := *chunkPointer
	for i := range chunk {
		chunk[i] = ""
	}
}
