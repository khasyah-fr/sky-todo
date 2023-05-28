package config

import (
	"os"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/khasyah-fr/sky-todo/entities/activities"
	"github.com/khasyah-fr/sky-todo/entities/todos"
)

type Database struct {
	connection *gorm.DB
	once       sync.Once
}

func (db *Database) lazyInit() {
	db.once.Do(func() {
		host := os.Getenv("MYSQL_HOST")
		port := os.Getenv("MYSQL_PORT")
		username := os.Getenv("MYSQL_USER")
		password := os.Getenv("MYSQL_PASSWORD")
		dbname := os.Getenv("MYSQL_DBNAME")
		dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

		conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("Cannot connect to the database")
		}

		err = conn.AutoMigrate(&activities.Activity{}, &todos.Todo{})
		if err != nil {
			panic("Failed to auto-migrate database schema")
		}

		db.connection = conn
	})
}

func (db *Database) GetConnection() *gorm.DB {
	db.lazyInit()
	return db.connection
}

var DB = &Database{}
