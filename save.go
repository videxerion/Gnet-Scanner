package main

import (
	"os"
	"strings"
	"time"
)

type save struct {
	version     uint64
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
	// Получаем версию сохранения
	retSave.version = bytesToUint64(content[:8])
	// Получаем кол-во отсканированных адресов
	retSave.scannedAddr = bytesToUint64(content[8:16])
	// Получаем путь до базы данных
	dbNameLen := content[16]
	retSave.dbName = bytesToString(content[17 : dbNameLen+17])
	// Получаем сеть и маску
	retSave.inputNet = bytesToString(content[18+dbNameLen:])

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

func saveThread(inputNet string, dbName string) {
	var err error

	if !directoryExist("saves") && pathToSaves == "saves/" {
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
		file, err = os.OpenFile(pathToSaves+filename, os.O_RDWR|os.O_CREATE, 0755)
	} else {
		file, err = os.OpenFile(saveFlag, os.O_RDWR|os.O_CREATE, 0755)
	}

	if err != nil {
		panic(err)
	}

	defer file.Close()

	for !exitState {
		time.Sleep(time.Second * 1)

		// Инициализируем переменную с будущим содержимым файла
		var wrBuf []byte

		// Записываем номер версии
		versionBytes := uint64ToBytes(fileSaveVersion)
		wrBuf = append(wrBuf, versionBytes...)

		// Записываем сколько адресов было отсканировано
		scanned := uint64ToBytes(scannedAddress)
		wrBuf = append(wrBuf, scanned...)

		// Записываем путь до базы данных
		dbNameBytes := stringToBytes(dbName)
		length := byte(len(dbNameBytes))
		dbNameBytes = append([]byte{length}, dbNameBytes...)
		wrBuf = append(wrBuf, dbNameBytes...)

		// Записываем маску и сеть
		netBytes := stringToBytes(inputNet)
		length = byte(len(netBytes))
		netBytes = append([]byte{length}, netBytes...)
		wrBuf = append(wrBuf, netBytes...)

		// Если файл не пустой то очищаем
		if !isEmptyFile(file) {
			err = clearFile(file)
			if err != nil {
				panic(err)
			}
		}

		// Записываем новые данные в файл
		_, err = file.Write(wrBuf)
		if err != nil {
			panic(err)
		}

	}
}
