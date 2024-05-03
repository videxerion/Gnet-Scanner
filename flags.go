package main

import (
	"flag"
	"os"
	"time"
)

func pathExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}

func getFlags() {
	// Получение флагов
	debugFlagObj := flag.Bool("debug", false, "Enable debug for pprof on localhost:6060")
	inputNetObj := flag.String("network", "None", "Network for scanning")
	saveFlagObj := flag.String("save", "None", "Path to save file")
	chunkSizeObj := flag.Uint64("ChunkSize", 100, "Number of addresses in the chunk")
	limitThreadsObj := flag.Uint64("threads", 50, "Number of scanning threads")
	connectTimeoutObj := flag.Uint64("ConnectTimeout", 100, "Sets the length of time to wait for a connection (ms)")
	readTimeoutObj := flag.Uint64("ReadTimeout", 250, "Sets the length of time to wait for reading (ms)")
	responseSizeObj := flag.Uint64("ResponseSize", GigaByte, "Sets the maximum response size (bytes)")
	pathToDbObj := flag.String("PathToBD", "results/", "Sets the path where the usrDatabase will be created")
	flag.Parse()
	debugFlag = *debugFlagObj
	inputNet = *inputNetObj
	saveFlag = *saveFlagObj
	chunkSize = *chunkSizeObj
	limitThreads = *limitThreadsObj
	connectTimeout = time.Duration(*connectTimeoutObj)
	readTimeout = time.Duration(*readTimeoutObj)
	responseSize = *responseSizeObj
	pathToDb = *pathToDbObj
}

func checkValidFlags() {
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
	if responseSize > GigaByte {
		println("The maximum response size cannot exceed 2 gigabytes because it corresponds to the BLOB type in sqlite")
		os.Exit(0)
	}

	if !pathExist(pathToDb) {
		println("The specified path to the directory where the usrDatabase should be saved does not exist")
		os.Exit(0)
	}

	if usrSave.version != fileSaveVersion {
		println("The received version of the save file is not compatible with the current version of the program")
		os.Exit(0)
	}

}
