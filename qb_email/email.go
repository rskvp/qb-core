package qb_email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_utils"
)

type EmailHelper struct {
}

var Email *EmailHelper

func init() {
	Email = new(EmailHelper)
}

// NewSender returns a new Message Sender
func (instance *EmailHelper) NewSender(settings ...interface{}) (*SmtpSender, error) {
	sender := new(SmtpSender)
	if len(settings) > 0 {
		err := sender.Configure(settings[0])
		if nil != err {
			return nil, err
		}
	}
	return sender, nil
}

// NewMessage returns a new Message that can compose an email with attachments
func (instance *EmailHelper) NewMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/plain")
}

// NewHTMLMessage returns a new Message that can compose an HTML email with attachments
func (instance *EmailHelper) NewHTMLMessage(subject string, body string) *Message {
	return newMessage(subject, body, "text/html")
}

// Send sends the message.
func (instance *EmailHelper) Send(addr string, auth smtp.Auth, m *Message) error {
	return smtp.SendMail(addr, auth, m.From.Address, m.GetToList(), m.GetBytes())
}

// SendSecure sends the message over TLS.
func (instance *EmailHelper) SendSecure(addr string, auth smtp.Auth, tlsConfig *tls.Config, m *Message) error {
	host, _, _ := net.SplitHostPort(addr)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Quit()

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}
	// To && From
	if err = c.Mail(m.From.Address); err != nil {
		return err
	}
	toList := m.GetToList()
	for _, addr := range toList {
		if len(addr) > 0 {
			if err = c.Rcpt(addr); err != nil {
				return err
			}
		}
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(m.GetBytes())
	if err != nil {
		return err
	}

	return w.Close()
}

func (instance *EmailHelper) SendMessage(host string, port int, secure bool, user string, pass string,
	from, replyTo string, to, bcc, cc string, subject string, text string, html string, attachments []interface{}) (err error) {
	// prepare server data
	servername := fmt.Sprintf("%v:%v", host, port)
	auth := smtp.PlainAuth("", user, pass, host)
	toList := qb_utils.Strings.Split(to, ";,")
	bccList := qb_utils.Strings.Split(bcc, ";,")
	ccList := qb_utils.Strings.Split(cc, ";,")

	var m *Message
	if len(html) > 0 {
		m = instance.NewHTMLMessage(subject, html)
	} else {
		m = instance.NewMessage(subject, text)
	}

	addr, err := mail.ParseAddress(from)
	if nil != err {
		m.From = &mail.Address{Address: from}
	} else {
		m.From = addr
	}

	if len(replyTo) == 0 {
		replyTo = from
	}
	m.ReplyTo = replyTo

	m.To = toList
	if len(bccList) > 0 {
		m.Bcc = bccList
	}
	if len(ccList) > 0 {
		m.Cc = ccList
	}

	for _, attachment := range attachments {
		if nil != attachment {
			if v, b := attachment.(string); b {
				err = addAttachmentString(m, v)
			} else if v, b := attachment.(map[string]interface{}); b {
				err = addAttachmentObject(m, v)
			}
		}
	}
	if nil != err {
		return
	}
	if secure {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
		err = instance.SendSecure(servername, auth, tlsconfig, m)
		return
	} else {
		err = instance.Send(servername, auth, m)
		return
	}
}

// ---------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
// ---------------------------------------------------------------------------------------------------------------------

func newMessage(subject string, body string, bodyContentType string) *Message {
	m := &Message{Subject: subject, Body: body, BodyContentType: bodyContentType}
	m.Attachments = make(map[string]*Attachment)
	return m
}

func addAttachmentString(m *Message, attachment string) error {
	filename := qb_utils.Paths.FileName(attachment, true)
	return addAttachment(m, filename, attachment)
}

func addAttachmentObject(m *Message, attachment map[string]interface{}) error {
	filename := qb_utils.Reflect.GetString(attachment, "filename")
	path := qb_utils.Reflect.GetString(attachment, "path")
	return addAttachment(m, filename, path)
}

func addAttachment(m *Message, filename, path string) error {
	if len(filename) > 0 && len(path) > 0 {
		data, err := download(path)
		if nil != err {
			return err
		}
		return m.AddAttachmentBinary(filename, data, false)
	}
	return nil // nothing to attach
}

func download(url string) ([]byte, error) {
	if len(url) > 0 {
		if strings.Index(url, "http") > -1 {
			// HTTP
			tr := &http.Transport{
				MaxIdleConns:       10,
				IdleConnTimeout:    15 * time.Second,
				DisableCompression: true,
			}
			client := &http.Client{Transport: tr}
			resp, err := client.Get(url)
			if nil == err {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if nil == err {
					return body, nil
				} else {
					return []byte{}, err
				}
			} else {
				return []byte{}, err
			}
		} else {
			// FILE SYSTEM
			path := url
			return qb_utils.IO.ReadBytesFromFile(path)
		}
	}
	return []byte{}, qb_utils.Errors.Prefix(errors.New(url), "Invalid url or path: ")
}
