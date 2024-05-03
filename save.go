package main

import (
	"bytes"
	"encoding/binary"
	"os"
	"strings"
	"time"
)

type save struct {
	scannedAddr uint64
	dbName      string
	inputNet    string
}

func parseSaveFile(path string) save {
	content, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var retSave save
	retSave.scannedAddr = bytesToUint64(content[:8])

	dbNameLen := content[8]
	retSave.dbName = bytesToString(content[9 : dbNameLen+9])

	retSave.inputNet = bytesToString(content[9+dbNameLen+1:])

	return retSave

}

func isEmptyFile(file *os.File) bool {
	// Получаем информацию о файле
	fileInfo, err := file.Stat()
	if err != nil {
		panic(err)
	}

	// Проверяем размер файла
	if fileInfo.Size() == 0 {
		return true
	} else {
		return false
	}
}

func clearFile(file *os.File) error {
	// Установка позиции чтения/записи в начало файла
	_, err := file.Seek(0, 0)
	if err != nil {
		return err
	}

	// Очистка содержимого файла
	err = file.Truncate(0)
	if err != nil {
		return err
	}

	return nil
}

func uint64ToBytes(num uint64) []byte {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, num)
	if err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

func stringToBytes(str string) []byte {
	return []byte(str)
}

func bytesToString(b []byte) string {
	return string(b)
}

func saveThread(inputNet string, dbName string) {
	var err error

	if !directoryExist("saves") {
		os.Mkdir("saves", 0755)
	}

	if saveFlag != "None" {
		usrSave = parseSaveFile(saveFlag)
	}

	var filename string

	if saveFlag == "None" {
		filename = time.Now().Format(time.DateTime) + " " + strings.Replace(inputNet, "/", "-", -1) + ".gsv"
	} else {
		filename = saveFlag
	}

	var file *os.File

	if saveFlag == "None" {
		file, err = os.OpenFile("saves/"+filename, os.O_RDWR|os.O_CREATE, 0755)
	} else {
		file, err = os.OpenFile(saveFlag, os.O_RDWR|os.O_CREATE, 0755)
	}

	if err != nil {
		panic(err)
	}

	defer file.Close()

	for !exitState {
		time.Sleep(time.Second * 1)

		scanned := uint64ToBytes(scannedAddress)

		var wrBuf []byte

		wrBuf = append(wrBuf, scanned...)

		dbNameBytes := stringToBytes(dbName)
		length := byte(len(dbNameBytes))
		dbNameBytes = append([]byte{length}, dbNameBytes...)

		wrBuf = append(wrBuf, dbNameBytes...)

		netBytes := stringToBytes(inputNet)
		length = byte(len(netBytes))
		netBytes = append([]byte{length}, netBytes...)

		wrBuf = append(wrBuf, netBytes...)

		if !isEmptyFile(file) {
			err = clearFile(file)
			if err != nil {
				panic(err)
			}
		}

		_, err = file.Write(wrBuf)
		if err != nil {
			panic(err)
		}

	}
}
