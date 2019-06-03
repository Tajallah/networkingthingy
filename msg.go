package msg

import (
	"fmt"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Message struct {
	gorm.Model
	Author int `json:author`
	Text string `string`
}

func (m Message) String() string {
	return fmt.Sprintf("%s :: %s", m.Author, m.Text)
}
