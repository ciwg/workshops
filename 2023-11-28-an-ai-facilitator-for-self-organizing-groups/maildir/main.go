package main

import (
	"bytes"
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

var sysMsgSend = "you are %s. take action on your inbox content to accomplish the goal quickly with minimal administrative overhead, then send a message to whoever should act next.  you can invite new participants into the team by emailing them.  any file you create must be included as a text/plain MIME attachment.  your answer must be in the form of a properly-formatted rfc822 email message as shown.  include no markdown wrappers."

const TokenLimit = 8192.0 // XXX get from grokker model

func main() {
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
		status := os.Args[3]
		tokenLimit, err := strconv.Atoi(os.Args[4])
		Ck(err)
		// summarize the inbox
		_, err = summarizeInbox(homedir, status, tokenLimit)
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
	Ck(err)

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
		recipient := extractEmail(rawRecipient)
		recipients = append(recipients, recipient)
	}

	// XXX use git

	// write the message to the recipient maildirs
	for _, recipient := range recipients {
		deliver(recipient, msg)
	}

	// put a copy in the sender's inbox/cur
	sender := extractEmail(msg.From)
	mdmsg, err := deliver(sender, msg)
	Ck(err)
	// move the message from inbox/new to inbox/cur
	key, err := mdmsg.Process()
	Ck(err, "failed to move message %s", key)

	return
}

// extractEmail extracts the email address from an email address string
func extractEmail(in string) (out string) {
	// use a regex to extract the email address
	// - the email address must include @
	// - the email address may or may not be enclosed in <>
	re := regexp.MustCompile(`([\w\-\.]+\@[\w\-\.]+)`)
	m := re.FindStringSubmatch(in)
	Assert(len(m) == 2, "invalid address: %s", in)
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

	// detach files from messages
	for _, key := range newkeys {
		mdmsg := newMap[key]
		err = detachFiles(homedir, mdmsg)
		Ck(err)
	}

	// summarize inbox/new to reduce redundancy
	_, err = summarizeInbox(homedir, "new", TokenLimit/2.0)

	// process each message
	for _, key := range newkeys {
		msg := newMap[key]
		process(homedir, inbox, msg)
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

// process processes a message for a single recipient. It moves the
// message from inbox/new to inbox/cur.
func process(homedir string, inbox *maildir.Maildir, mdmsg *lib.Message) (err error) {
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

// respond responds to msg.  It looks in the cur inbox, considers msg,
// creates a new message, and sends it to the given io.Writer.
func respond(homedir string, msg *email.Email, out io.Writer) (err error) {
	defer Return(&err)

	Fpf(os.Stderr, "%s responding to message: %s\n", homedir, msg.Subject)

	// XXX move grokker setup from cli.go to api.go so we can call
	// functions in api.go directly instead of going through Cli()
	config := grokker.NewCliConfig()

	// get msg token count
	msgBuf, err := msg.Bytes()
	Ck(err)
	msgTxt := string(msgBuf)
	tcMsg := tokenCount(msgTxt)

	// calculate the maximum token count for the summary
	// - summaryTokenLimit + tcMsg <= TokenLimit/2
	summaryTokenLimit := TokenLimit/2.0 - tcMsg

	// summarize inbox into an in-memory mbox
	mbox, err := summarizeInbox(homedir, "cur", summaryTokenLimit)
	Ck(err)

	// add msg to the mbox
	err = mbox.append(msg)
	Ck(err)

	// ask OpenAI to create a new message
	config.Stdin = strings.NewReader(mbox.String())
	config.Stdout = out
	sysmsg := Spf(sysMsgSend, extractEmail(homedir))
	args := []string{"msg", sysmsg}
	Fpf(os.Stderr, "grokker args=%v\n", args)
	_, err = grokker.Cli(args, config)
	Ck(err)

	// close the output
	if closer, ok := out.(io.Closer); ok {
		err = closer.Close()
		Ck(err)
	}

	return
}

// summarizeInbox summarizes the messages in the given user's
// inbox/{status}.  It returns a string containing the summary.
func summarizeInbox(homedir, status string, tokenLimit int) (mbox *mboxMem, err error) {
	defer Return(&err)

	// read the messages in the inbox
	inboxPath := filepath.Join(homedir, "inbox")
	inbox := maildir.NewMaildir(inboxPath)
	curMap, err := inbox.List(status)
	Ck(err)

	// convert the map to a sorted slice of keys
	var curkeys []string
	for key := range curMap {
		curkeys = append(curkeys, key)
	}
	sort.Strings(curkeys)
	Debug("curkeys=%v", curkeys)

	// collect the messages into mbox format
	mbox = &mboxMem{}
	for _, key := range curkeys {
		mdmsg := curMap[key]
		msgTxt, err := mdmsg.GetData()
		Ck(err)
		msg, err := email.NewEmailFromReader(strings.NewReader(msgTxt))
		Ck(err)
		err = mbox.append(msg)
		Ck(err)
	}
	Debug("mbox len=%d", len(mbox.String()))

	// keep summarizing until summary token count is less than
	// tokenLimit
	mbox, err = summarizeUntil(homedir, status, mbox, tokenLimit)
	Ck(err)

	// if we didn't replace the mbox, we're done
	if !mbox.replaced {
		Debug("mbox not replaced")
		return
	}
	Debug("mbox replaced")

	// put the summary mbox messages into inbox/new
	msgs, err := mbox.messages()
	Ck(err)
	for _, msg := range msgs {
		buf, err := msg.Bytes()
		Ck(err)
		// put msg in inbox/new
		mdmsg, err := inbox.Add(string(buf))
		Ck(err)
		if status == "cur" {
			// move msg from inbox/new to inbox/cur
			_, err = mdmsg.Process()
			Ck(err)
		}
	}

	// move old messages from inbox to archive/new
	for _, key := range curkeys {
		mdmsg := curMap[key]
		err = moveToArchive(homedir, mdmsg)
		Ck(err)
	}

	return
}

// summarizeUntil recursively summarizes the given mbox until the
// token count is less than tokenLimit.  It returns the resulting
// mbox.
func summarizeUntil(homedir, status string, mbox *mboxMem, tokenLimit int) (*mboxMem, error) {
	// if mbox is less than tokenLimit, return mbox
	tcMbox := tokenCount(mbox.String())
	Debug("tcMbox=%d", tcMbox)
	if tcMbox < tokenLimit {
		return mbox, nil
	}
	// otherwise, use grokker to summarize and recurse
	Fpf(os.Stderr, "%s summarizing inbox/%s, token count=%d\n", homedir, status, tcMbox)
	config := grokker.NewCliConfig()
	config.Stdin = strings.NewReader(mbox.String())
	config.Stdout = bytes.NewBuffer(nil)
	args := []string{"msg", Spf(sysMsgSummarize, extractEmail(homedir))}
	Fpf(os.Stderr, "grokker args=%v\n", args)
	_, err := grokker.Cli(args, config)
	if err != nil {
		return nil, err
	}
	summary := config.Stdout.(*bytes.Buffer).String()
	mbox.replaceWith(summary)
	return summarizeUntil(homedir, status, mbox, tokenLimit)
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
	mbox := mbox.NewWriter(&m.buf)
	Ck(err)
	// write the message to the mbox
	msgWriter, err := mbox.CreateMessage(msg.From, time.Now())
	Ck(err)
	msgbuf, err := msg.Bytes()
	Ck(err)
	_, err = msgWriter.Write(msgbuf)
	Ck(err)
	// close the mbox file
	err = mbox.Close()
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
	mbox := mbox.NewReader(&m.buf)
	// read each message
	for {
		msg, err := mbox.NextMessage()
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
