package main

import (
	"encoding/json"
	"fmt"
	"github.com/BrianLeishman/go-imap"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

type imapInfo struct {
	Ssl    string `json:"ssl"`
	Smtp   string `json:"smtp"`
	Imap   string `json:"imap"`
	Login  string `json:"login"`
	Passwd string `json:"passwd"`
}

func readMailConfig(cfg *imapInfo) error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	cfgfile := filepath.Join(user.HomeDir, ".imap.json")
	dat, err := os.ReadFile(cfgfile)
	check(err)
	_ = json.Unmarshal(dat, &cfg)
	return nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func imapConnect() (*imap.Dialer, error) {
	var cfg imapInfo
	if err := readMailConfig(&cfg); err != nil {
		return nil, err
	}
	//imap.Verbose = true
	im, err := imap.New(cfg.Login, cfg.Passwd, cfg.Imap, 993)
	if err != nil {
		return nil, err
	}
	return im, nil
}

func main() {

	im, err := imapConnect()
	check(err)
	defer im.Close()

	folders, err := im.GetFolders()
	check(err)
	fmt.Println(folders)

	err = im.SelectFolder("[Gmail]/rcgroups")
	check(err)

	uids, err := im.GetUIDs("ALL")
	check(err)
	fmt.Println(uids)

	emails, err := im.GetEmails(uids...)
	check(err)
	fmt.Println(len(emails))
	if len(emails) == 0 {
		return
	}
	for zz, v := range emails {
		fd, err := os.Create("/tmp/posts/" + strconv.Itoa(zz))
		check(err)
		defer fd.Close()
		fd.WriteString(v.Text)
		fmt.Println("----------------------------------------------------------------------------------------------------")
		fmt.Println(zz)
		fmt.Println(v.Subject)
		fmt.Println(v.MessageID)
		fmt.Println(v.Flags)
		//fmt.Println(v.Text)

	}

}
