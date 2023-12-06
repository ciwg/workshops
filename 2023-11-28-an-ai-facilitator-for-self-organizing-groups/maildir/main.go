package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/amalfra/maildir/v3"
	"github.com/amalfra/maildir/v3/lib"
	"github.com/emersion/go-mbox"
	"github.com/jordan-wright/email"
	. "github.com/stevegt/goadapt"
	"github.com/stevegt/grokker/v3"
)

var sysMsgSummarize = "you are %s.  rewrite the given mbox messages, reducing the number of messages without losing information.  your answer must be in the form of a series of properly-formatted mbox messages as shown.  include no markdown wrapper."

var sysMsgSend = "you are %s.  the requirements are below.  do what you can to meet them, then send a message to whoever should act next.  you can invite new participants into the team by emailing them.  your answer must be in the form of a MIME-formatted email message with any file attached as a valid text/plain MIME attachment with no encoding.  include no markdown wrappers.\n\n%s"

const TokenLimit = 8192.0 // XXX get from grokker model

var Requirements string

func main() {
	// load requirements
	buf, err := ioutil.ReadFile("requirements.md")
	Ck(err)
	Requirements = string(buf)

	// get subcommand
	Assert(len(os.Args) > 1, "missing subcommand")
	cmd := os.Args[1]
	switch cmd {
	case "distribute":
		// read a message from stdin and write it to the recipients'
		// inboxes
		recipients := distribute(os.Stdin)
		// print recipients
		for _, recipient := range recipients {
			Pl(recipient)
		}
	case "send":
		// format and send a message to stdout
		// get args
		from := os.Args[2]
		tostr := os.Args[3]
		ccstr := os.Args[4]
		subject := os.Args[5]
		// body is on stdin
		body, err := ioutil.ReadAll(os.Stdin)
		Ck(err)
		// parse the address strings and add domain if needed
		froms := parseAddrs(from)
		from = froms[0]
		to := parseAddrs(tostr)
		cc := parseAddrs(ccstr)
		// ensure we always cc the sender
		// cc = appendIfMissing(cc, from)
		// cc = appendIfMissing(cc, "facilitator@example.com")
		send(os.Stdout, from, to, cc, subject, body)
	case "summarize":
		// summarize and replace the given user's inbox/{status}
		homedir := os.Args[2]
		ratio, err := strconv.ParseFloat(os.Args[3], 64)
		Ck(err)
		// summarize the inbox
		_, err = summarizeInbox(homedir, ratio)
		Ck(err)
	case "readMail":
		// read, process, and respond to all new messages for the given user
		homedir := os.Args[2]
		err := readMail(homedir)
		Ck(err)
	case "cycle":
		// one cycle of readMail for each user
		err := cycle()
		Ck(err)
	case "users":
		// get the list of users
		users, err := getUsers()
		Ck(err)
		// print the list of users
		for _, user := range users {
			Pl(user)
		}
	default:
		Assert(false, "unknown subcommand: %s", cmd)
	}
}

// cycle reads, processes, and responds to all new messages for each
// user.  It returns an error if any user fails.
func cycle() (err error) {
	defer Return(&err)

	// get the list of users
	users, err := getUsers()
	Ck(err)

	// process each user
	for _, user := range users {
		err = readMail(user)
		Ck(err)
	}

	return
}

// getUsers returns a slice of strings containing the names of all
// users.  It reads the list of users from the home directory names.
func getUsers() (users []string, err error) {
	defer Return(&err)

	// get the list of home directories
	allDirs, err := ioutil.ReadDir(".")
	Ck(err)

	// get the list of users
	for _, dir := range allDirs {
		if !dir.IsDir() {
			continue
		}
		// a home directory is a directory that has @ in its name
		if !regexp.MustCompile(`\@`).MatchString(dir.Name()) {
			continue
		}
		users = append(users, dir.Name())
	}

	return
}

// appendIfMissing appends an element to a slice if it is not already
// in the slice
func appendIfMissing(slice []string, s string) []string {
	for _, ele := range slice {
		if ele == s {
			return slice
		}
	}
	return append(slice, s)
}

// parseAddrs parses a comma-separated list of email addresses and
// returns a slice of strings.  It appends "@example.com" to any
// address that doesn't already have an @.
func parseAddrs(in string) (out []string) {
	re := regexp.MustCompile(`\s*,\s*`)
	addrs := re.Split(in, -1)
	for _, addr := range addrs {
		// trim whitespace
		addr = strings.TrimSpace(addr)
		// skip empty addresses
		if addr == "" {
			continue
		}
		if !regexp.MustCompile(`\@`).MatchString(addr) {
			addr += "@example.com"
		}
		out = append(out, addr)
	}
	return
}

// send formats a message and writes it to the given io.Writer.
// This is not normally called by an agent.  It is used by the
// user interface to format and send a message to get things started.
func send(out io.Writer, from string, to, cc []string, subject string, body []byte) {

	// Create a new email message
	e := email.NewEmail()
	e.From = from
	e.To = to
	e.Cc = cc
	e.Subject = subject
	e.Text = body

	// write the email to output
	msgbuf, err := e.Bytes()
	Ck(err)
	_, err = out.Write(msgbuf)
	Ck(err)
}

// distribute reads an rfc822 message from an io.Reader and writes it to the
// inbox maildir(s) for the recipient(s)
func distribute(in io.Reader) (recipients []string) {
	// read the message
	msg, err := email.NewEmailFromReader(in)
	if err != nil {
		Fpf(os.Stderr, "error reading message: %s\n", err.Error())
		return
	}

	// show the message content on stderr
	buf, err := msg.Bytes()
	Ck(err)
	Fpf(os.Stderr, "-----------------------------\n%s\n--------------------------\n", string(buf))

	Fpf(os.Stderr, "distributing message: %s\n", msg.Subject)

	// get the recipient(s)
	var rawRecipients []string
	rawRecipients = append(rawRecipients, msg.To...)
	rawRecipients = append(rawRecipients, msg.Cc...)
	// extract the email address from each recipient
	for _, rawRecipient := range rawRecipients {
		recipient, err := extractEmail(rawRecipient)
		if err != nil {
			Fpf(os.Stderr, err.Error())
			continue
		}
		recipients = append(recipients, recipient)
	}

	// XXX use git

	// write the message to the recipient maildirs
	for _, recipient := range recipients {
		deliver(recipient, msg)
	}

	// put a copy in the sender's inbox/cur
	sender, err := extractEmail(msg.From)
	if err != nil {
		Fpf(os.Stderr, err.Error())
	} else {
		mdmsg, err := deliver(sender, msg)
		Ck(err)
		// move the message from inbox/new to inbox/cur
		key, err := mdmsg.Process()
		Ck(err, "failed to move message %s", key)

		// detach any files
		err = detachFiles(sender, mdmsg)
		Ck(err)
	}

	return
}

// extractEmail extracts the email address from an email address string
func extractEmail(in string) (out string, err error) {
	// use a regex to extract the email address
	// - the email address must include @
	// - the email address may or may not be enclosed in <>
	re := regexp.MustCompile(`([\w\-\.]+\@[\w\-\.]+)`)
	m := re.FindStringSubmatch(in)
	if len(m) != 2 {
		return "", fmt.Errorf("invalid address: %s", in)
	}
	out = m[1]
	return
}

// deliver receives a message for a single recipient.  It writes
// the message to the recipient's inbox.
func deliver(homedir string, msg *email.Email) (mdmsg *lib.Message, err error) {
	defer Return(&err)

	Fpf(os.Stderr, "delivering message to %s: %s\n", homedir, msg.Subject)

	// create the recipient's home directory if it doesn't exist
	err = os.MkdirAll(homedir, 0755)
	Ck(err)

	// ensure there's an inbox
	inboxPath := filepath.Join(homedir, "inbox")
	inbox := maildir.NewMaildir(inboxPath)

	// write the message to inbox/new
	buf, err := msg.Bytes()
	Ck(err)
	mdmsg, err = inbox.Add(string(buf))
	Ck(err)
	return
}

// readMail reads and processes all new messages for the given user.
func readMail(homedir string) (err error) {
	defer Return(&err)

	Fpf(os.Stderr, "%s reading mail\n", homedir)

	// ensure there's an inbox
	inboxPath := filepath.Join(homedir, "inbox")
	inbox := maildir.NewMaildir(inboxPath)

	// get all new messages
	newMap, err := inbox.List("new")
	Ck(err)

	// convert the map to a sorted slice of keys
	var newkeys []string
	for key := range newMap {
		newkeys = append(newkeys, key)
	}
	sort.Strings(newkeys)

	// summarize inbox/new to reduce redundancy
	_, err = summarizeInbox(homedir, 0.5)

	// XXX don't return mboxNew above, do re-read from inbox/new in
	// respond()

	// respond to entire mbox with one message
	// create a pipe from respond to distribute
	pr, pw := io.Pipe()
	go respond(homedir, pw)
	// send the response
	_ = distribute(pr)

	// move old messages from inbox/new to inbox/cur
	newMap, err = inbox.List("new")
	for _, mdmsg := range newMap {
		_, err = mdmsg.Process()
		Ck(err)
	}

	return
}

// detachFiles detaches files from the given message.  It saves the
// files to the recipient's home directory and removes the files from
// the message.  It overwrites any existing file from an attachment of
// the same name.  // XXX use git
func detachFiles(homedir string, mdmsg *lib.Message) (err error) {
	defer Return(&err)

	// get the message from the maildir
	data, err := mdmsg.GetData()
	Ck(err)
	msg, err := email.NewEmailFromReader(strings.NewReader(data))
	Ck(err)

	// save any attachments to the recipient's home directory
	for _, att := range msg.Attachments {
		Fpf(os.Stderr, "%s saving attachment: %s\n", homedir, att.Filename)
		// create the file
		fh, err := os.Create(homedir + "/" + att.Filename)
		Ck(err)
		// write the attachment to the file
		_, err = fh.Write(att.Content)
		Ck(err)
		// close the file
		err = fh.Close()
		Ck(err)
	}
	return
}

/*
// process processes a message for a single recipient. It moves the
// message from inbox/new to inbox/cur.
func XXXprocess(homedir string, inbox *maildir.Maildir, mdmsg *lib.Message) (err error) {
	defer Return(&err)

	// get the message from the maildir
	data, err := mdmsg.GetData()
	Ck(err)
	msg, err := email.NewEmailFromReader(strings.NewReader(data))
	Ck(err)

	Fpf(os.Stderr, "%s processing message: %s\n", homedir, msg.Subject)

	// respond to the message if it's not from us
	sender := extractEmail(msg.From)
	if sender != homedir {
		// create a pipe from respond to distribute
		pr, pw := io.Pipe()
		go respond(homedir, msg, pw)
		// send the response
		_ = distribute(pr)
	}

	// move the message from inbox/new to inbox/cur
	key, err := mdmsg.Process()
	Ck(err, "failed to process message %s", key)

	return
}
*/

// respond responds to msg.  It looks in the given mboxMem,
// creates a new message, and sends it to the given io.Writer.
func respond(homedir string, out io.Writer) (err error) {
	defer Return(&err)

	Fpf(os.Stderr, "%s responding\n", homedir)

	// read the messages in the inbox into mboxIn
	inboxPath := filepath.Join(homedir, "inbox")
	inbox := maildir.NewMaildir(inboxPath)
	curMap, err := inbox.List("new")
	Ck(err)
	// convert the map to a sorted slice of keys
	var curkeys []string
	for key := range curMap {
		curkeys = append(curkeys, key)
	}
	sort.Strings(curkeys)
	Debug("curkeys=%v", curkeys)
	// collect the messages into mbox format
	mboxIn := &mboxMem{}
	for _, key := range curkeys {
		mdmsg := curMap[key]
		msgTxt, err := mdmsg.GetData()
		Ck(err)
		msg, err := email.NewEmailFromReader(strings.NewReader(msgTxt))
		Ck(err)
		err = mboxIn.append(msg)
		Ck(err)
	}

	msgs, err := mboxIn.messages()
	Ck(err)
	if len(msgs) == 0 {
		Fpf(os.Stderr, "%s no messages to respond to -- sending presence message\n", homedir)
		subject := Spf("%s is available", homedir)
		send(out, homedir, []string{"facilitator@example.com"}, nil, subject, []byte("How can I help?"))
	} else {
		for _, msg := range msgs {
			Fpf(os.Stderr, "%s responding to message: %s\n", homedir, msg.Subject)
		}

		// XXX move grokker setup from cli.go to api.go so we can call
		// functions in api.go directly instead of going through Cli()
		config := grokker.NewCliConfig()

		// ask OpenAI to create a new message
		config.Stdin = strings.NewReader(mboxIn.String())
		config.Stdout = out
		user, err := extractEmail(homedir)
		Ck(err)
		sysmsg := Spf(sysMsgSend, user, Requirements)
		args := []string{"msg", sysmsg}
		// Fpf(os.Stderr, "grokker args=%v\n", args)
		_, err = grokker.Cli(args, config)
		Ck(err)
	}

	// close the output
	if closer, ok := out.(io.Closer); ok {
		err = closer.Close()
		Ck(err)
	}

	return
}

// summarizeInbox summarizes the messages in the given user's
// inbox/{status}.  It returns a string containing the summary.
func summarizeInbox(homedir string, ratio float64) (mboxOut *mboxMem, err error) {
	defer Return(&err)

	// read the messages in the inbox
	inboxPath := filepath.Join(homedir, "inbox")
	inbox := maildir.NewMaildir(inboxPath)
	curMap, err := inbox.List("new")
	Ck(err)

	// convert the map to a sorted slice of keys
	var curkeys []string
	for key := range curMap {
		curkeys = append(curkeys, key)
	}
	sort.Strings(curkeys)
	Debug("curkeys=%v", curkeys)

	// collect the messages into mbox format
	mboxOut = &mboxMem{}
	for _, key := range curkeys {
		mdmsg := curMap[key]
		msgTxt, err := mdmsg.GetData()
		Ck(err)
		msg, err := email.NewEmailFromReader(strings.NewReader(msgTxt))
		Ck(err)
		err = mboxOut.append(msg)
		Ck(err)
	}
	Debug("mbox len=%d", len(mboxOut.String()))

	// keep summarizing until summary token count is less than
	// tokenLimit/2
	mboxOut, err = summarizeUntil(homedir, mboxOut, ratio)
	Ck(err)

	// if we didn't replace the mbox, we're done
	if !mboxOut.replaced {
		Debug("mbox not replaced")
		return
	}
	Debug("mbox replaced")

	// put the summary mbox messages into inbox/new
	msgs, err := mboxOut.messages()
	Ck(err)
	Fpf(os.Stderr, "%s has %d messages\n", homedir, len(msgs))
	for _, msg := range msgs {
		buf, err := msg.Bytes()
		Ck(err)
		// put msg in inbox/new
		_, err = inbox.Add(string(buf))
		Ck(err)
		Fpf(os.Stderr, "%s added message to inbox/new: %s\n", homedir, msg.Subject)
	}

	// move old messages from inbox/new to inbox/cur
	for _, key := range curkeys {
		mdmsg := curMap[key]
		// err = moveToArchive(homedir, mdmsg)
		_, err = mdmsg.Process()
		Ck(err)
	}

	return
}

// summarizeUntil recursively summarizes the given mbox until the
// token count is less than tokenLimit.  It returns the resulting
// mbox.
func summarizeUntil(homedir string, mboxIn *mboxMem, ratio float64) (*mboxMem, error) {
	// if mbox is less than tokenLimit/2, return mbox
	tcMbox := tokenCount(mboxIn.String())
	Debug("tcMbox=%d", tcMbox)
	if float64(tcMbox) < float64(TokenLimit)*ratio {
		Fpf(os.Stderr, "%s no more summarizing needed, token count=%d\n", homedir, tcMbox)
		return mboxIn, nil
	}
	Fpf(os.Stderr, "%s summarizing, token count=%d\n", homedir, tcMbox)
	// if mbox is greater than TokenLimit, split it in half and recurse
	if tcMbox > TokenLimit {
		// split mbox in half
		mbox1 := &mboxMem{}
		mbox2 := &mboxMem{}
		msgs, err := mboxIn.messages()
		Ck(err)
		start2 := len(msgs) / 2
		for i := 0; i < start2; i++ {
			err = mbox1.append(msgs[i])
			Ck(err)
		}
		for i := start2; i < len(msgs); i++ {
			err = mbox2.append(msgs[i])
			Ck(err)
		}
		// summarize each half
		mbox1, err = summarizeUntil(homedir, mbox1, 0.25)
		Ck(err)
		mbox2, err = summarizeUntil(homedir, mbox2, 0.25)
		Ck(err)
		// combine the summaries
		mboxOut := &mboxMem{}
		mboxOut.buf.Write(mbox1.buf.Bytes())
		mboxOut.buf.Write(mbox2.buf.Bytes())
		return summarizeUntil(homedir, mboxOut, 0.5)
	}

	// otherwise, use grokker to summarize and recurse
	config := grokker.NewCliConfig()
	config.Stdin = strings.NewReader(mboxIn.String())
	config.Stdout = bytes.NewBuffer(nil)
	user, err := extractEmail(homedir)
	Ck(err)
	args := []string{"msg", Spf(sysMsgSummarize, user)}
	// Fpf(os.Stderr, "grokker args=%v\n", args)
	Fpf(os.Stderr, "%s summarizing, waiting for OpenAI...\n", homedir)
	_, err = grokker.Cli(args, config)
	if err != nil {
		return nil, err
	}
	summary := config.Stdout.(*bytes.Buffer).String()
	Fpf(os.Stderr, "%s summarizing done, token count=%d\n", homedir, tokenCount(summary))
	mboxIn.replaceWith(summary)
	return summarizeUntil(homedir, mboxIn, 0.5)
}

// tokenCount returns the number of tokens in the given string.
func tokenCount(in string) (tc int) {
	// use grokker to count tokens
	config := grokker.NewCliConfig()
	config.Stdin = strings.NewReader(in)
	config.Stdout = &bytes.Buffer{}
	args := []string{"tc"}
	_, err := grokker.Cli(args, config)
	Ck(err)
	tcStr := config.Stdout.(*bytes.Buffer).String()
	tcStr = strings.TrimSpace(tcStr)
	tc, err = strconv.Atoi(tcStr)
	Ck(err, "invalid tc output: '%s'", tcStr)
	Debug("in len=%d tc=%d", len(in), tc)
	return
}

// mboxMem is a memory-based mbox file
type mboxMem struct {
	buf      bytes.Buffer
	replaced bool
}

// append appends the given message to the in-memory mbox
func (m *mboxMem) append(msg *email.Email) (err error) {
	defer Return(&err)
	w := mbox.NewWriter(&m.buf)
	Ck(err)
	// write the message to the mbox
	msgWriter, err := w.CreateMessage(msg.From, time.Now())
	Ck(err)
	msgbuf, err := msg.Bytes()
	Ck(err)
	_, err = msgWriter.Write(msgbuf)
	Ck(err)
	// close the writer
	err = w.Close()
	Ck(err)
	return
}

// String returns the mbox contents as a string
func (m *mboxMem) String() string {
	return m.buf.String()
}

// messages returns the mbox contents as a slice of email messages
func (m *mboxMem) messages() (msgs []*email.Email, err error) {
	defer Return(&err)
	// create an mbox reader
	r := mbox.NewReader(&m.buf)
	// read each message
	for {
		msg, err := r.NextMessage()
		if err == io.EOF {
			break
		}
		Ck(err)
		// convert the message to an email
		emailMsg, err := email.NewEmailFromReader(msg)
		Ck(err)
		msgs = append(msgs, emailMsg)
	}
	return
}

// replaceWith replaces the mbox contents with the given string
func (m *mboxMem) replaceWith(in string) {
	m.buf.Reset()
	m.buf.WriteString(in)
	m.replaced = true
}

// moveToArchive moves the given message from inbox/cur to archive/new
func moveToArchive(homedir string, mdmsg *lib.Message) (err error) {
	defer Return(&err)

	// ensure there's an archive
	archivePath := filepath.Join(homedir, "archive")
	archive := maildir.NewMaildir(archivePath)

	// get the message from the maildir
	data, err := mdmsg.GetData()
	Ck(err)
	msg, err := email.NewEmailFromReader(strings.NewReader(data))
	Ck(err)

	// write the message to archive/new
	buf, err := msg.Bytes()
	Ck(err)
	_, err = archive.Add(string(buf))
	Ck(err)

	// delete the message from inbox/cur
	inboxPath := filepath.Join(homedir, "inbox")
	inbox := maildir.NewMaildir(inboxPath)
	err = inbox.Delete(mdmsg.Key())
	Ck(err)

	return
}
