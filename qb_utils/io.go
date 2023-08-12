package qb_utils

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type IoHelper struct {
}

var IO *IoHelper

func init() {
	IO = new(IoHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *IoHelper) NewDownloadSession(actions interface{}) *DownloadSession {
	return newDownloadSession(actions)
}

func (instance *IoHelper) NewDownloader() *Downloader {
	return newDownloader()
}

func (instance *IoHelper) NewDownloaderAction(uid, source, sourceversion, target string) *DownloaderAction {
	return newAction(uid, source, sourceversion, target)
}

func (instance *IoHelper) FileSize(filename string) (int64, error) {
	info, err := os.Stat(filename)
	if nil != err {
		return 0, err
	}
	return info.Size(), nil
}

func (instance *IoHelper) Remove(filename string) error {
	return os.Remove(filename)
}

func (instance *IoHelper) RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func (instance *IoHelper) RemoveSilent(filename string) {
	_ = os.Remove(filename)
}

func (instance *IoHelper) RemoveAllSilent(path string) {
	_ = os.RemoveAll(path)
}

func (instance *IoHelper) MoveFile(from, to string) error {
	if b, _ := Paths.IsFile(to); !b {
		to = Paths.Concat(to, Paths.FileName(from, true))
	}
	_, err := instance.CopyFile(from, to)
	if nil != err {
		return err
	}
	return instance.Remove(from)
}

func (instance *IoHelper) CopyFile(src, dst string) (int64, error) {
	Paths.Mkdir(dst)

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func (instance *IoHelper) AppendTextToFile(text, file string) (bytes int, err error) {
	var f *os.File
	if b, _ := Paths.Exists(file); b {
		f, err = os.OpenFile(file,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.Create(file)
	}

	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		bytes, err = w.WriteString(text)
		_ = w.Flush()
	}
	return bytes, err
}

func (instance *IoHelper) WriteTextToFile(text, file string) (bytes int, err error) {
	var f *os.File
	f, err = os.Create(file)

	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		bytes, err = w.WriteString(text)
		_ = w.Flush()
	}
	return bytes, err
}

func (instance *IoHelper) AppendBytesToFile(data []byte, file string) (bytes int, err error) {
	var f *os.File
	if b, _ := Paths.Exists(file); b {
		f, err = os.OpenFile(file,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.Create(file)
	}

	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		bytes, err = w.Write(data)
		_ = w.Flush()
	}
	return bytes, err
}

func (instance *IoHelper) WriteBytesToFile(data []byte, file string) (bytes int, err error) {
	var f *os.File
	f, err = os.Create(file)
	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		bytes, err = w.Write(data)
		w.Flush()
	}
	return bytes, err
}

func (instance *IoHelper) ReadBytesFromFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := io.ReadAll(file)
	return b, err
}

func (instance *IoHelper) ReadTextFromFile(fileName string) (string, error) {
	b, err := instance.ReadBytesFromFile(fileName)
	return string(b), err
}

func (instance *IoHelper) Download(url string) ([]byte, error) {
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
				if resp.StatusCode < 300 {
					body, err := ioutil.ReadAll(resp.Body)
					if nil == err {
						return body, nil
					} else {
						return []byte{}, err
					}
				} else {
					return []byte{}, errors.New(fmt.Sprintf("%s: %s", resp.Status, url))
				}
			} else {
				return []byte{}, err
			}
		} else {
			// FILE SYSTEM
			return instance.ReadBytesFromFile(url)
		}
	}
	return []byte{}, Errors.Prefix(errors.New("missing_url"), "Missing Parameter 'url': ")
}

func (instance *IoHelper) ReadHashFromFile(fileName string) (string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func (instance *IoHelper) ReadHashFromBytes(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	hash := sha256.New()
	if _, err := io.Copy(hash, reader); err != nil {
		return "", err
	}
	return hex.EncodeToString(hash.Sum(nil)), nil
}

// ScanBytesFromFile read a file line by line
func (instance *IoHelper) ScanBytesFromFile(fileName string, callback func(data []byte) bool) error {
	if nil == callback {
		return errors.New("missing_callback")
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if callback(scanner.Bytes()) {
			// exit loop
			return nil
		}
	}

	return scanner.Err()
}

// ScanTextFromFile read a text file line by line
func (instance *IoHelper) ScanTextFromFile(fileName string, callback func(text string) bool) error {
	if nil == callback {
		return errors.New("missing_callback")
	}
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if callback(scanner.Text()) {
			// exit loop
			return nil
		}
	}

	return scanner.Err()
}

func (instance *IoHelper) ReadLinesFromFile(fileName string, count int) string {
	counter := 0
	var buf strings.Builder
	if count > 0 {
		_ = instance.ScanTextFromFile(fileName, func(text string) bool {
			counter++
			buf.WriteString(text + "\n")
			return counter == count
		})
	}
	return buf.String()
}

func (instance *IoHelper) ReadAllBytes(reader io.Reader) ([]byte, error) {
	return io.ReadAll(reader)
}

func (instance *IoHelper) ReadAllString(reader io.Reader) (string, error) {
	buf, err := io.ReadAll(reader)
	if nil != err {
		return "", err
	}
	return string(buf), nil
}

func (instance *IoHelper) Chmod(filename string, mode os.FileMode) (changed bool, err error) {
	var stats os.FileInfo
	stats, err = os.Stat(filename)
	if nil != err {
		return
	}
	err = os.Chmod(filename, mode)
	if nil != err {
		return
	}

	oldMode := stats.Mode()
	stats, err = os.Stat(filename)
	if nil != err {
		return
	}
	changed = oldMode != stats.Mode()

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
