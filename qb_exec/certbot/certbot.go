package certbot

import (
	"os/exec"
)

// https://certbot.eff.org/instructions?ws=other&os=osx
// https://eff-certbot.readthedocs.io/en/stable/using.html#webroot

var (
	Certbot        *CertbotHelper
	certbotCommand = "certbot"
)

const wpName = "certbot"
const fsName = "./fs-certbot"

type CertbotHelper struct {
	root    string
	rootTmp string
}

func init() {
	certbotCommand = findExecPath()
	Certbot = new(CertbotHelper)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

// NewExec
// Creates new Solana command with default password
func (instance *CertbotHelper) NewExec() *CertbotExec {
	if nil != instance {
		return NewExec(certbotCommand, nil)
	}
	return nil
}

// findExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather aggressive search, which is the same in all systems.
func findExecPath() string {
	for _, path := range [...]string{
		// Unix-like
		"certbot",
		"/usr/bin/certbot",
		"/opt/homebrew/Cellar/certbot/1.29.0/bin/certbot",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "certbot"
}
