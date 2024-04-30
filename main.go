package main

import (
	"flag"
	"log"
	"net"
	"os"
	"sync"
	"time"
)
import _ "net/http/pprof"

var mu sync.Mutex
var pauseMu sync.Mutex
var dbWriteMu sync.Mutex
var ctDbThMu sync.Mutex
var countThreads int
var countAddress int
var scannedAddress int = 0
var speed int
var pauseState bool = false
var countDbThreads int = 0
var usrSave save
var saveFlag *string

const chunkSize int = 100
const limitThreads = 300

// TODO: Сделать систему сохранений

func main() {
	var err error
	log.Println("Getting flags")

	debugFlag := flag.Bool("debug", false, "enable debug for pprof on localhost:6060")
	inputNet := flag.String("network", "None", "a network")
	saveFlag = flag.String("save", "None", "path to save file")
	flag.Parse()

	if *saveFlag != "None" {
		usrSave = parseSaveFile(*saveFlag)
		scannedAddress = int(usrSave.scannedAddr)
	}

	if *saveFlag != "None" && *inputNet != "None" {
		println("It is not possible to use the --network flag together with the --save flag")
		os.Exit(0)
	}

	if *debugFlag {
		go prof()
		go leakDetector()
	}

	var ipnet *net.IPNet
	if *saveFlag == "None" {
		_, ipnet, err = net.ParseCIDR(*inputNet)
	} else {
		_, ipnet, err = net.ParseCIDR(usrSave.inputNet)

	}

	if err != nil {
		panic(err)
	}

	firstIP := ipnet.IP

	if *inputNet != "None" || *saveFlag != "None" {
		log.Println("Creating database")
		var database *db
		if *saveFlag == "None" {
			database = Database(time.Now().Format(time.DateTime) + " result.db")
		} else {
			database = Database(usrSave.dbName)
		}
		database.createTable()

		var chunk [chunkSize]string
		pointer := 0

		start := time.Now()

		if *saveFlag == "None" {
			countAddress, err = countIPAddresses(*inputNet)
		} else {
			countAddress, err = countIPAddresses(usrSave.inputNet)
		}
		if err != nil {
			panic(err)
		}

		go speedMeter()
		go infoScreen()
		go interceptingKeystrokes()
		if *saveFlag == "None" {
			go saveThread(*inputNet, database.name)
		} else {
			go saveThread(usrSave.inputNet, database.name)
		}

		var startIP net.IP

		if *saveFlag == "None" {
			startIP = firstIP.Mask(ipnet.Mask)
		} else {
			startIPint := ipToUint32(firstIP.Mask(ipnet.Mask))
			newIPInt := startIPint + uint32(scannedAddress)
			startIP = uint32ToIP(newIPInt)
		}

		for ip := startIP; ipnet.Contains(ip); nextIP(ip) {
			if !pauseState {
				if ip != nil {
					if isValidIP(ip) {
						chunk[pointer] = ip.String()
						if pointer == chunkSize-1 {
							for {
								if countThreads < limitThreads {
									mu.Lock()
									countThreads += 1
									mu.Unlock()
									go scanChunk(chunk, database)
									break
								}
							}

							for i := 0; i < chunkSize; i++ {
								chunk[i] = ""
								pointer = 0
							}
						} else {
							pointer++
						}
					} else {
						mu.Lock()
						scannedAddress++
						mu.Unlock()
					}
				} else {
					for {
						if countThreads < limitThreads {
							mu.Lock()
							countThreads += 1
							mu.Unlock()
							go scanChunk(chunk, database)
							break
						}
					}
					break
				}
			} else {
				for {
					if !pauseState {
						break
					}
				}
			}
		}

		for {
			if countThreads == 0 {
				break
			}
		}

		end := time.Now()

		println("Scan time:", end.Unix()-start.Unix(), "s.")
	} else {
		print("Please use --network flag for set network for scanning in format: address/mask")
	}
}
