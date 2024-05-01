package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

type Db struct {
	handler *sql.DB
	name    string
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

func Database(name string) *Db {
	if !directoryExist("results") {
		os.Mkdir("results", 0755)
	}

	database, err := sql.Open("sqlite3", "results/"+name)

	if err != nil {
		log.Fatal(err.Error())
	}

	database.SetMaxOpenConns(0)
	database.SetMaxIdleConns(100)

	return &Db{handler: database, name: name}
}

func (d Db) createTable() {
	_, err := d.handler.Exec("CREATE TABLE IF NOT EXISTS results (ID INTEGER PRIMARY KEY, ip TEXT, response BLOB)")
	if err != nil {
		log.Fatal(err.Error())
	}
}

func (d Db) Add(ip string, response string) {
	// Изменяем счётчик потоков
	incCommonVar(&countDbThreads, &ctDbThMu)
	defer subCommonVar(&countDbThreads, &ctDbThMu)

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
