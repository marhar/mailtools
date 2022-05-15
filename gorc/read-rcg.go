package main

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	_ "github.com/emersion/go-message/charset"
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
	fmt.Println("Flags for [Gmail]/rcgroups:", mbox.Flags)
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

func doSearch(c *client.Client) {
	_, err := c.Select("[Gmail]/rcgroups", false)
	check(err)
	x, err := c.Search(&imap.SearchCriteria{Body: []string{"FliteTest"}})
	check(err)
	fmt.Println(x)
	x, err = c.UidSearch(&imap.SearchCriteria{Body: []string{"FliteTest"}})
	check(err)
	fmt.Println(x)
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
		fmt.Println("Date:", date)
	}
	if from, err := header.AddressList("From"); err == nil {
		fmt.Println("From:", from)
	}
	if to, err := header.AddressList("To"); err == nil {
		fmt.Println("To:", to)
	}
	if subject, err := header.Subject(); err == nil {
		fmt.Println("Subject:", subject)
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

func connectToMail() (*client.Client, error) {
	user, err := user.Current()
	check(err)
	cfgfile := filepath.Join(user.HomeDir, ".imap.json")
	dat, err := os.ReadFile(cfgfile)
	check(err)
	cfg := imapInfo{}
	_ = json.Unmarshal(dat, &cfg)
	c, err := client.DialTLS(cfg.Imap+":993", nil)
	check(err)
	err = c.Login(cfg.Login, cfg.Passwd)
	check(err)
	return c, nil
}

type rcgMessage struct {
	msgno int
	body  string
}

func fetchRcgMessages(c *client.Client) []rcgMessage {
	mbox, err := c.Select("[Gmail]/rcgroups", false)
	check(err)
	seqset := new(imap.SeqSet)
	seqset.AddRange(1, mbox.Messages)

	messages := make(chan *imap.Message, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	}()
	var results []rcgMessage
	for msg := range messages {
		fmt.Println("=======================================================================================")
		fmt.Println("* " + msg.Envelope.Subject)
		fmt.Println(msg.Envelope.MessageId)
		msg.BodyStructure.Walk(func(path []int, part *imap.BodyStructure) (walkChildren bool) {
			fmt.Println(path)
			_ = part
			return true
		})
	}
	err = <-done
	check(err)
	return results
}

func main() {

	c, err := connectToMail()
	check(err)
	defer c.Logout()
	fetchRcgMessages(c)

	/*	connect
			get list of rcg messages
		   	loop over messages
				extract url
				open url
			delete message
	*/
	//listMailboxes(c)
	//showSubjects(c)
	//showBodies(c)
	//doSearch(c)
}
