package db

import (
	"../msg"
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
	defer db.Close()
	db.AutoMigrate(&msg.Message{})
	db.Create(&msg.Messgae{Author: m.Author, Text: m.Text})
	return nil
}

//this present the last 20 messages. Will use this to test sending a slice of messages
func Last20 () (msgs []msg.Message) error {
	var holder msg.Message
	db, err := gorm.Open(DIALECT, DATABASE)
	if err != nil {
		return err
	}
	//open 20 into slice
}
