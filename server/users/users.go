package users

import (
	"fmt"
	"os"
	"crypto/rsa"
	"encoding/json"
	"crypto/rsa"
	"crypto/sha256"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

//constants
const TAGDIR = "datapile/tags"
const USERDIR = "datapile/users"

var ROOTUSER = Tag {
	"ROOTUSER",
	"datapile/icons/ROOTUSER.svg",
	[4]uint8(0, 0, 0, 0)
	PrivaledgePage{
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true,
		true
	}

}

func checkErr (e error) {
	if e != nil {
		panic e
	}
}

type Tag struct {
	Title string
	IconPath string
	Color [4]uint8 //RGBA
	Privs PrivaledgePage
}

func (t *Tag) Export () []byte {
	var eslice []byte
	eslice = append(eslice, []byte(t.Title))
	eslice = append(eslice, []byte("|"))
	eslice = append(eslice, []byte(t.IconPath))
	eslice = append(eslice, []byte("|"))
	eslice = append(eslice, []byte(t.Color))
	eslice = append(eslice, []byte("|"))
	eslice = append(eslice, t.Privs.Export())
}

func (t *Tag) CheckExist () (bool, error) {
	if 
}

type PrivaledgePage struct {
	Owner bool
	Admin bool
	VeiwAudit bool
	ManageServer bool
	ManageTags bool
	MangeChannels bool
	Kick bool
	Ban bool
	SendMsg bool
	SendMedia bool
	InternalPing bool
	ExternalPing bool
	Blacklist bool
}

func Btoi (b bool) uint8 {

}

func (p * PrivaledgePage) Export () []byte {
	var bits = uint8

}

type User struct {
	PublicKey []byte
	Username string
	PublicName string
	Tags []Tag
	PingAddress []byte
}

func BundleInfo (PublicKey []byte, Username string, PublicName string, Tags []Tag, PingAddress []byte) (User, error) {
	var u = User{
		"",
		"",
		"",

	}
}

func (u User) NewUser (u User) error {

}
