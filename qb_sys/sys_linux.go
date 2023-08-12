//go:build linux
// +build linux

package qb_sys

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

const (
	// dbusPath is the default path for dbus machine id.
	dbusPath = "/var/lib/dbus/machine-id"
	// dbusPathEtc is the default path for dbus machine id located in /etc.
	// Some systems (like Fedora 20) only know this path.
	// Sometimes it's the other way round.
	dbusPathEtc = "/etc/machine-id"
)

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
	gio := &InfoObject{Kernel: osInfo[0], Core: osInfo[1], Platform: osInfo[2], OS: osInfo[3], GoOS: runtime.GOOS, CPUs: runtime.NumCPU()}
	gio.Hostname, _ = os.Hostname()
	return gio
}

func _getInfo() string {
	cmd := exec.Command("uname", "-srio")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("getInfo:", err)
	}
	return out.String()
}

// machineID returns the uuid specified at `/var/lib/dbus/machine-id` or `/etc/machine-id`.
// If there is an error reading the files an empty string is returned.
// See https://unix.stackexchange.com/questions/144812/generate-consistent-machine-unique-id
func machineID() (string, error) {
	id, err := _readFile(dbusPath)
	if err != nil {
		// try fallback path
		id, err = _readFile(dbusPathEtc)
	}
	if err != nil {
		return "", err
	}
	return _trim(string(id)), nil
}
