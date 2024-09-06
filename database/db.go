package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

const dbPath = "database/scheduler.db"

// СheckDB проверяет существет ли файл базы данных
// Создает, если не существует
func СheckDB() (*sql.DB, error) {
	appPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	dbFile := filepath.Join(appPath, dbPath)

	_, err = os.Stat(dbFile)

	var install bool
	if err != nil {
		install = true
	}
	// если install равен true, после создания БД требуется выполнить
	// sql-запрос с CREATE TABLE и CREATE INDEX
	if install {
		if err = CreateDBFile(dbPath); err != nil {
			return nil, err
		}
		fmt.Printf("Database file created successfully in %s!\n", dbPath)
		database, err := sql.Open("sqlite", dbPath)
		if err != nil {
			return nil, err
		}
		defer database.Close()
		if err = CreateTable(database); err != nil {
			return nil, err
		}
		fmt.Println("Table created successfully!")
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// CreateDBFile создает файл базы данных в db/scheduler.db
func CreateDBFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	file.Close()

	return nil
}

// CreateTable создает талицу с полями и индексирует поле дата
func CreateTable(db *sql.DB) error {
	scheduler_table := `CREATE TABLE scheduler (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "date" CHAR(8) NOT NULL DEFAULT "",
        "title" TEXT NOT NULL DEFAULT "",
        "comment" TEXT NOT NULL DEFAULT "",
        "repeat" VARCHAR(128) NOT NULL DEFAULT "");
		CREATE INDEX scheduler_date ON scheduler (date);`
	query, err := db.Prepare(scheduler_table)
	if err != nil {
		return err
	}
	_, err = query.Exec()
	if err != nil {
		return err
	}

	return nil
}
