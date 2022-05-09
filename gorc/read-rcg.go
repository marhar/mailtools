package main

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
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

func showSubjects(c *client.Client) {
	mbox, err := c.Select("[Gmail]/rcgroups", false)
	check(err)
	log.Println("Flags for [Gmail]/rcgroups:", mbox.Flags)
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)

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

func showBodies(c *client.Client) {
	mbox, err := c.Select("[Gmail]/rcgroups", false)
	check(err)
	seqset := new(imap.SeqSet)
	//seqset.AddRange(1, mbox.Messages)
	_ = mbox
	seqset.AddRange(1, 1)

	var section imap.BodySectionName
	items := []imap.FetchItem{section.FetchItem()}
	messages := make(chan *imap.Message, 1)
	go func() {
		err := c.Fetch(seqset, items, messages)
		check(err)
	}()

	msg := <-messages
	if msg == nil {
		log.Fatal("Server didn't returned message")
	}

	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Server didn't returned message body")
	}
	// Create a new mail reader
	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatal(err)
	}

	header := mr.Header
	if date, err := header.Date(); err == nil {
		log.Println("Date:", date)
	}
	if from, err := header.AddressList("From"); err == nil {
		log.Println("From:", from)
	}
	if to, err := header.AddressList("To"); err == nil {
		log.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		log.Println("Subject:", subject)
	}

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
	//showSubjects(c)
	showBodies(c)
}
