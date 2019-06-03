package db

import (
	"msg"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const DIALECT = "sqlite3"
const DATABASE = "messages.db"

func AddMsg (m msg.Message) error {
	db, err := gorm.Open(DIALECT, DATABASE)
	if err != nil {
		return err
	}
	db.AutoMigrate(&msg.Message{})
	db.Create(&msg.Messgae{Author: m.Author, Text: m.Text})
	return nil
}
