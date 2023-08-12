package qb_email

import (
	"crypto/tls"
	"fmt"
	"net/mail"
	"net/smtp"
	"strings"

	"github.com/rskvp/qb-core/qb_utils"
)

type SmtpSender struct {
	config *SmtpSettings
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

func (instance *SmtpSender) Configure(settings interface{}) (err error) {
	instance.config = new(SmtpSettings)
	if c, b := settings.(*SmtpSettings); b {
		instance.config = c
	} else if cc, bb := settings.(SmtpSettings); bb {
		instance.config = &cc
	} else if s, bs := settings.(string); bs && !strings.HasPrefix(s, "{") {
		text, err := qb_utils.IO.ReadTextFromFile(s)
		if nil == err && strings.HasPrefix(text, "{") {
			return instance.Configure(text)
		}
	} else {
		var config SmtpSettings
		err = qb_utils.JSON.Read(qb_utils.Convert.ToString(settings), &config)
		if nil == err && len(config.Host) > 0 {
			instance.config = &config
			return
		}
	}
	return
}

func (instance *SmtpSender) Send(subject, body string, to, bcc, cc []string, from, replyTo string, attachments []string) error {
	var message *Message
	if qb_utils.Regex.IsHTML(body) {
		message = Email.NewHTMLMessage(subject, body)
	} else {
		message = Email.NewMessage(subject, body)
	}

	// from
	if len(from) == 0 {
		from = instance.config.From
	}
	addr, err := mail.ParseAddress(from)
	if nil != err {
		message.From = &mail.Address{Address: from}
	} else {
		message.From = addr
	}

	// replyTo
	if len(replyTo) == 0 {
		replyTo = instance.config.ReplyTo
	}
	message.ReplyTo = replyTo

	if len(bcc) > 0 {
		message.Bcc = bcc
	}
	if len(cc) > 0 {
		message.Cc = cc
	}

	// to
	for _, s := range to {
		addr, err = mail.ParseAddress(s)
		if nil != err {
			message.AddTo(mail.Address{Address: s})
		} else {
			message.AddTo(*addr)
		}
	}

	// attachments
	if len(attachments) > 0 {
		for _, a := range attachments {
			err = message.AddAttachment(a)
			if nil != err {
				return err
			}
		}
	}

	return instance.SendMessage(message)
}

func (instance *SmtpSender) SendMessage(message *Message) (err error) {
	var secure bool
	var host, user, pass string
	var port int
	host = instance.config.Host
	port = instance.config.Port
	if nil != instance.config.Auth {
		secure = instance.config.Secure
		user = instance.config.Auth.User
		pass = instance.config.Auth.Pass
	}
	servername := fmt.Sprintf("%v:%v", host, port)
	auth := smtp.PlainAuth("", user, pass, host)

	if secure {
		// TLS config
		tlsconfig := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
		// try with custom TLS
		err = Email.SendSecure(servername, auth, tlsconfig, message)
		if nil != err {
			// try using SendMail
			err2 := Email.Send(servername, auth, message)
			if nil != err2 {
				err = qb_utils.Errors.Prefix(err2, fmt.Sprintf("Concatenated errors. 1-> %s; 2-> ", err))
			} else {
				err = nil
			}
		}
	} else {
		err = Email.Send(servername, auth, message)
	}
	return err
}

func (instance *SmtpSender) SendAsync(subject, body string, to, bcc, cc []string, from, replyTo string, attachments []string, callback func(error)) {
	go func() {
		err := instance.Send(subject, body, to, bcc, cc, from, replyTo, attachments)
		if nil != callback {
			callback(err)
		}
	}()
}

func (instance *SmtpSender) SendMessageAsync(message *Message, callback func(error)) {
	go func() {
		err := instance.SendMessage(message)
		if nil != callback {
			callback(err)
		}
	}()
}
