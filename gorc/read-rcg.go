package main

import (
	"encoding/json"
	"github.com/emersion/go-imap/client"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type imapInfo struct {
	Ssl    string `json:"ssl"`
	Smtp   string `json:"smtp"`
	Imap   string `json:"imap"`
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	user, err := user.Current()
	check(err)
	cfgfile := filepath.Join(user.HomeDir, ".imap.json")
	dat, err := os.ReadFile(cfgfile)
	check(err)
	cfg := imapInfo{}
	_ = json.Unmarshal(dat, &cfg)
	c, err := client.DialTLS(cfg.Imap, nil)
	check(err)
	err = c.Login(cfg.Login, cfg.Passwd)
	check(err)
	log.Println("Logged in")
	defer c.Logout()

}
