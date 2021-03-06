package msg

import (
	"fmt"
	"strconv"
	"encoding/json"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Message struct {
	gorm.Model
	Author int `json:author`
	Text string `json:text`
	Channel string `json:channel`
	Media string `json:media`
	signature []byte `json:signature`
}

func (m Message) Error() string {
	errstr :=  ("malformed Message object :\n\n "
}

func (m Message) Stringer() string {
	return fmt.Sprintf("%s :: %s", strconv.Itoa(m.Author), m.Text)
}

func (m Message) ToJson() ([]byte, error){
	byt, err := json.Marshal(m)
	if err != nil {
		return nil, err
	} else {
		return byt, nil
	}
}

func (m Message) FromJson(holder []byte) error {
	err := json.Unmarshal(holder, &m);
	if err != nil {
		return err
	} else {
		return nil
	}
}
