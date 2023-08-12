package getssl

import (
	"os/exec"

	"github.com/rskvp/qb-core/qb_utils"
)

// https://github.com/srvrco/getssl
// https://github.com/srvrco/getssl/wiki/Guide-to-getting-a-certificate-for-example.com-and-www.example.com

/**
GetSSL was written in standard bash ( so it can be run on a server, a desktop computer, or even a virtualbox) and
add the checks, and certificates to a remote server ( providing you have a ssh with key, sftp or ftp access to the
remote server).

FEATURES
------------------------
* Bash - It runs on virtually all unix machines, including BSD, most Linux distributions, macOS.
* Get certificates for remote servers - The tokens used to provide validation of domain ownership, and the certificates themselves can be automatically copied to remote servers (via ssh, sftp or ftp for tokens). The script doesn't need to run on the server itself. This can be useful if you don't have access to run such scripts on the server itself, e.g. if it's a shared server.
* Runs as a daily cron - so certificates will be automatically renewed when required.
* Automatic certificate renewals
* Checks certificates are correctly loaded - After installation of a new certificate it will test the port specified ( see Server-Types for options ) that the certificate is actually being used correctly.
* Automatically updates - The script can automatically update itself with bug fixes etc if required.
* Extensively configurable - With a simple configuration file for each certificate it is possible to configure it exactly for your needs, whether a simple single domain or multiple domains across multiple servers on the same certificate.
* Supports http and dns challenges - Full ACME implementation
* Simple and easy to use
* Detailed debug info - Whilst it shouldn't be needed, detailed debug information is available.
* Reload services - After a new certificate is obtained then the relevant services (e.g. apache/nginx/postfix) can be reloaded.
* ACME v1 and V2 - Supports both ACME versions 1 and 2 (note ACMEv1 is deprecated and clients will automatically use v2)
*/

var (
	GetSSL        *GetSslHelper
	getsslCommand = "getssl"
)

const wpName = "getssl"
const fsName = "./fs-getssl"

type GetSslHelper struct {
	root    string
	rootTmp string
}

func init() {
	getsslCommand = findExecPath()
	GetSSL = new(GetSslHelper)
}

// ---------------------------------------------------------------------------------------------------------------------
//	p u b l i c
// ---------------------------------------------------------------------------------------------------------------------

// NewExec
// Creates new Solana command with default password
func (instance *GetSslHelper) NewExec() *GetSslExec {
	if nil != instance {
		return NewExec(getsslCommand, nil)
	}
	return nil
}

// findExecPath tries to find the Chrome browser somewhere in the current
// system. It performs a rather aggressive search, which is the same in all systems.
func findExecPath() string {
	currentDir := qb_utils.Paths.Absolute(".")
	for _, path := range [...]string{
		// Unix-like
		"getssl",
		currentDir + "/bin/getssl",
		currentDir + "/fs-getssl/bin/getssl",
	} {
		found, err := exec.LookPath(path)
		if err == nil {
			return found
		}
	}
	// Fall back to something simple and sensible, to give a useful error
	// message.
	return "getssl"
}
