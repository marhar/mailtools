package main

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
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

func showRCG(c *client.Client) {
	mbox, err := c.Select("[Gmail]/rcgroups", false)
	check(err)
	log.Println("Flags for [Gmail]/rcgroups:", mbox.Flags)
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, 99999)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()
	for msg := range messages {
		fmt.Println("* " + msg.Envelope.Subject)
		fmt.Println(msg.Body)
	}
	err = <-done
	check(err)
}
func listMailboxes(c *client.Client) {
	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	fmt.Println("Mailboxes:")
	for m := range mailboxes {
		fmt.Println("* " + m.Name)
	}
	err := <-done
	check(err)
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

	//listMailboxes(c)
	showRCG(c)
}
