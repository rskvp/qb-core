//go:build freebsd || netbsd || openbsd || dragonfly || solaris
// +build freebsd netbsd openbsd dragonfly solaris

package qb_sys

import (
	"bytes"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const hostidPath = "/etc/hostid"

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func shutdown(adminPsw string) error {
	// echo <password> | sudo -S shutdown -h now
	if err := exec.Command("/bin/sh", "-c", "echo "+adminPsw+" | sudo -S shutdown -h now").Run(); err != nil {
		return err
	}
	return nil
}

func getInfo() *InfoObject {
	out := _getInfo()
	for strings.Index(out, "broken pipe") != -1 {
		out = _getInfo()
		time.Sleep(500 * time.Millisecond)
	}
	osStr := strings.Replace(out, "\n", "", -1)
	osStr = strings.Replace(osStr, "\r\n", "", -1)
	osInfo := strings.Split(osStr, " ")
	gio := &InfoObject{Kernel: osInfo[0], Core: osInfo[1], Platform: runtime.GOARCH, OS: osInfo[2], GoOS: runtime.GOOS, CPUs: runtime.NumCPU()}
	gio.Hostname, _ = os.Hostname()
	return gio
}

func _getInfo() string {
	cmd := exec.Command("uname", "-sri")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		//fmt.Println("getInfo:", err)
	}
	return out.String()
}

// machineID returns the uuid specified at `/etc/hostid`.
// If the returned value is empty, the uuid from a call to `kenv -q smbios.system.uuid` is returned.
// If there is an error an empty string is returned.
func machineID() (string, error) {
	id, err := readHostid()
	if err != nil {
		// try fallback
		id, err = readKenv()
	}
	if err != nil {
		return "", err
	}
	return id, nil
}

func readHostid() (string, error) {
	buf, err := _readFile(hostidPath)
	if err != nil {
		return "", err
	}
	return _trim(string(buf)), nil
}

func readKenv() (string, error) {
	buf := &bytes.Buffer{}
	err := _run(buf, os.Stderr, "kenv", "-q", "smbios.system.uuid")
	if err != nil {
		return "", err
	}
	return _trim(buf.String()), nil
}
