package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type db struct {
	handler *sql.DB
	name    string
}

func addDbThread() {
	ctDbThMu.Lock()
	countDbThreads += 1
	ctDbThMu.Unlock()
}

func subDbThread() {
	ctDbThMu.Lock()
	countDbThreads -= 1
	ctDbThMu.Unlock()
}

func directoryExist(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			return false
		}
	} else {
		return fileInfo.IsDir()
	}
}

func Database(name string) *db {
	if !directoryExist("results") {
		os.Mkdir("results", 0755)
	}

	database, err := sql.Open("sqlite3", "results/"+name)

	if err != nil {
		log.Fatal(err.Error())
	}

	database.SetMaxOpenConns(0)
	database.SetMaxIdleConns(100)

	return &db{handler: database, name: name}
}

func (d db) createTable() {
	_, err := d.handler.Exec("CREATE TABLE IF NOT EXISTS results (ID INTEGER PRIMARY KEY, ip TEXT, response BLOB)")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (d db) Add(ip string, response string) {
	// Изменяем счётчик потоков
	addDbThread()
	defer subDbThread()

	// Блокируем для записи
	dbWriteMu.Lock()
	defer dbWriteMu.Unlock()

	// Создаём транзакцию
	var tx *sql.Tx
	var err error

	tx, err = d.handler.Begin()
	if err != nil {
		println(ip)
		panic(err)
	}

	defer tx.Rollback()

	// Формируем запрос
	query := "INSERT INTO results (ip, response) VALUES (?, ?)"
	var stmt *sql.Stmt

	stmt, err = tx.Prepare(query)
	if err != nil {
		println(ip)
		panic(err)
	}
	defer stmt.Close()

	// Выполняем запрос
	_, err = stmt.Exec(ip, response)
	if err != nil {
		println(ip)
		panic(err)
	}

	// Завершаем транзакцию
	err = tx.Commit()
	if err != nil {
		println(ip)
		panic(err)
	}
}
