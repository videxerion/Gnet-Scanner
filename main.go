package main

import (
	"flag"
	"net"
	"os"
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
)

// Состояния для управления
var pauseState = false

var usrSave save

func main() {
	var err error

	// Получение флагов
	debugFlagObj := flag.Bool("debug", false, "Enable debug for pprof on localhost:6060")
	inputNetObj := flag.String("network", "None", "Network for scanning")
	saveFlagObj := flag.String("save", "None", "Path to save file")
	chunkSizeObj := flag.Uint64("ChunkSize", 100, "Number of addresses in the chunk")
	limitThreadsObj := flag.Uint64("threads", 50, "Number of scanning threads")
	connectTimeoutObj := flag.Uint64("ConnectTimeout", 100, "Sets the length of time to wait for a connection (ms)")
	readTimeoutObj := flag.Uint64("ReadTimeout", 250, "Sets the length of time to wait for reading (ms)")
	responseSizeObj := flag.Uint64("ResponseSize", GigaByte*2, "Sets the maximum response size (bytes)")
	flag.Parse()
	debugFlag = *debugFlagObj
	inputNet = *inputNetObj
	saveFlag = *saveFlagObj
	chunkSize = *chunkSizeObj
	limitThreads = *limitThreadsObj
	connectTimeout = time.Duration(*connectTimeoutObj)
	readTimeout = time.Duration(*readTimeoutObj)
	responseSize = *responseSizeObj

	// Если не указан ни один из обязательных флагов то сообщаем об этом
	if inputNet == "None" && saveFlag == "None" {
		println("To start scanning, you must specify the network using the flag: --network {ip/mask}")
		os.Exit(0)
	}

	// Если указаны сразу оба обязательных флага то сообщаем об этом
	if saveFlag != "None" && inputNet != "None" {
		println("It is not possible to use the --network flag together with the --save flag")
		os.Exit(0)
	}

	// Если установленный размер ответа превышает 2 гигабайта то сообщаем об этом
	if responseSize > GigaByte*2 {
		println("The maximum response size cannot exceed 2 gigabytes because it corresponds to the BLOB type in sqlite")
		os.Exit(0)
	}

	// Если указан путь к файлу сохранения то пытаемся его распарсить
	if saveFlag != "None" {
		usrSave = parseSaveFile(saveFlag)
		scannedAddress = usrSave.scannedAddr
	}

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
	var database *Db
	if saveFlag == "None" {
		database = Database(time.Now().Format(time.DateTime) + " result.db")
	} else {
		database = Database(usrSave.dbName)
	}
	database.createTable()

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
		go saveThread(inputNet, database.name)
	} else {
		go saveThread(usrSave.inputNet, database.name)
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
								go scanChunk(chunk, database)
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
						go scanChunk(chunk, database)
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
