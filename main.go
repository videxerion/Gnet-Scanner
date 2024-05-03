package main

import (
	"net"
	"sync"
	"time"
)
import _ "net/http/pprof"

// Создаём мютексы
var (
	countThreadsMu sync.Mutex
	pauseMu        sync.Mutex
	dbWriteMu      sync.Mutex
	ctDbThMu       sync.Mutex
	exitMu         sync.Mutex
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
	debugFlag bool
	inputNet  string
	saveFlag  string

	// Флаги параметров
	chunkSize      uint64
	limitThreads   uint64
	connectTimeout time.Duration
	readTimeout    time.Duration
	responseSize   uint64
	pathToDb       string
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
	if saveFlag == "None" {
		go saveThread(inputNet, usrDatabase.name)
	} else {
		go saveThread(usrSave.inputNet, usrDatabase.name)
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

	for ip := startIP; ipnet.Contains(ip); nextIP(ip) {
		if !pauseState {
			if ip != nil {
				if isValidIP(ip) {
					chunk[pointer] = ip.String()
					if pointer == chunkSize-1 {
						for {
							if countThreads < limitThreads {
								incCommonVar(&countThreads, &countThreadsMu)
								chunkCopy := make([]string, chunkSize)
								copy(chunkCopy, chunk)
								go scanChunk(chunkCopy, usrDatabase)
								break
							}
						}

						for i := uint64(0); i < chunkSize; i++ {
							chunk[i] = ""
							pointer = 0
						}
					} else {
						pointer++
					}
				} else {
					incCommonVar(&scannedAddress, &countThreadsMu)
				}
			} else {
				for {
					if countThreads < limitThreads {
						incCommonVar(&countThreads, &countThreadsMu)
						go scanChunk(chunk, usrDatabase)
						break
					}
				}
				break
			}
		} else {
			waitPauseEnd()
		}
	}
	waitCompletionThreads()
}

func waitCompletionThreads() {
	for {
		if countThreads == 0 {
			break
		}
	}
}

func waitPauseEnd() {
	for {
		if !pauseState {
			break
		}
	}
}
