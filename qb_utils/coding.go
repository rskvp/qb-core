package qb_utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"
	"hash"
	"html"
	"io"
	"net/url"
	"os"
	"strings"
)

type CodingHelper struct {
}

var Coding *CodingHelper

func init() {
	Coding = new(CodingHelper)
}

//----------------------------------------------------------------------------------------------------------------------
//	URL
//----------------------------------------------------------------------------------------------------------------------

func (instance *CodingHelper) UrlEncode(text string) string {
	return url.QueryEscape(text)
}

func (instance *CodingHelper) UrlDecode(text string) string {
	value, err := url.QueryUnescape(text)
	if nil != err {
		return text
	}
	return value
}

//UrlEncodeAll
// Convert every char in escaped query string
// i.e. "Ciao" become "%43%69%61%6F"
func (instance *CodingHelper) UrlEncodeAll(text string) (response string) {
	for _, c := range text {
		s := string(c)
		u := url.QueryEscape(s)
		if u == s {
			response += fmt.Sprintf("%%%X", c)
		} else {
			response += u
		}
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	html
//----------------------------------------------------------------------------------------------------------------------

//HtmlEncodeAll
// Convert every char in UTF-8 escaped HTML string
// i.e. "Ciao" become "&#67;&#105;&#97;&#111;"
func (instance *CodingHelper) HtmlEncodeAll(text string) (response string) {
	for _, c := range text {
		s := string(c)
		u := html.EscapeString(s)
		if u == s {
			response += fmt.Sprintf("&#%v;", c)
		} else {
			response += u
		}
	}
	return
}

func (instance *CodingHelper) HtmlEncode(text string) (response string) {
	response = html.EscapeString(text)
	return
}

func (instance *CodingHelper) HtmlDecode(text string) (response string) {
	response = html.UnescapeString(text)
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	base64
//----------------------------------------------------------------------------------------------------------------------

func (instance *CodingHelper) EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func (instance *CodingHelper) DecodeBase64(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

func (instance *CodingHelper) WrapBase64(value string) string {
	if len(value) > 0 {
		if strings.Index(value, "base64-") == -1 {
			return "base64-" + instance.EncodeBase64([]byte(value))
		}
	}
	return value
}

func (instance *CodingHelper) UnwrapBase64(value string) string {
	if len(value) > 0 {
		clean := strings.Replace(value, "base64-", "", 1)
		if len(clean) > 3 {
			// regular base64 prefixed
			dec64, err := instance.DecodeBase64(clean)
			if nil == err {
				result := string(dec64)
				if clean != value {
					return result
				}
				if Compare.IsStringASCII(result) {
					return result
				}
			}
		}
	}
	return value
}

//----------------------------------------------------------------------------------------------------------------------
//	RSA
//----------------------------------------------------------------------------------------------------------------------

func (instance *CodingHelper) GenerateSessionKey() [32]byte {
	// crypto/rand.Reader is a good source of entropy for blinding the RSA
	// operation.
	rng := rand.Reader
	key := make([]byte, 32)
	if _, err := io.ReadFull(rng, key); err != nil {
		panic("RNG failure")
	}
	return sha256.Sum256(key)
}

// GenerateKeyPair generates a new an RSA keypair of the given bit size using the
// random source random (for example, crypto/rand.Reader).
func (instance *CodingHelper) GenerateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	return privkey, &privkey.PublicKey, nil
}

// PrivateKeyToBytes convert private key to bytes
func (instance *CodingHelper) PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)
	return privBytes
}

// PublicKeyToBytes convert public key to bytes
func (instance *CodingHelper) PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

// BytesToPrivateKey bytes to private key
func (instance *CodingHelper) BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		// is encrypted pem block
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	key, err := x509.ParsePKCS1PrivateKey(b)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// BytesToPublicKey bytes to public key
func (instance *CodingHelper) BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	enc := x509.IsEncryptedPEMBlock(block)
	b := block.Bytes
	var err error
	if enc {
		// is encrypted pem block"
		b, err = x509.DecryptPEMBlock(block, nil)
		if err != nil {
			return nil, err
		}
	}
	ifc, err := x509.ParsePKIXPublicKey(b)
	if err != nil {
		return nil, err
	}
	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not ok")
	}
	return key, nil
}

// EncryptWithPublicKey encrypts data with public key
func (instance *CodingHelper) EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}

// DecryptWithPrivateKey decrypts data with private key
func (instance *CodingHelper) DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	HASHING
//----------------------------------------------------------------------------------------------------------------------

func (instance *CodingHelper) MD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func (instance *CodingHelper) SHA1(text string) string {
	h := sha1.Sum([]byte(text))
	return fmt.Sprintf("%x", h)
}

func (instance *CodingHelper) SHA256FromFile(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", nil
	}
	defer f.Close()

	h := sha256.New()
	if _, e := io.Copy(h, f); e != nil {
		return "", nil
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func (instance *CodingHelper) SHA256FromText(text string) string {
	h := sha256.Sum256([]byte(text))
	return fmt.Sprintf("%x", h)
}

func (instance *CodingHelper) SHA256(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

func (instance *CodingHelper) SHA512(text string) string {
	h := sha512.Sum512([]byte(text))
	return fmt.Sprintf("%x", h)
}

func (instance *CodingHelper) Hash(h func() hash.Hash, secret []byte, message []byte) []byte {
	hasher := hmac.New(h, secret)
	hasher.Write(message)
	return hasher.Sum(nil)
}

func (instance *CodingHelper) HashSha256(secret []byte, message []byte) []byte {
	return instance.Hash(sha256.New, secret, message)
}

func (instance *CodingHelper) HashSha512(secret []byte, message []byte) []byte {
	return instance.Hash(sha512.New, secret, message)
}

// EncryptTextWithPrefix Encrypt using AES with a 32 byte key and adding a prefix to avoid multiple encryption
// Encrypted code is recognizable by the prefix.
// Useful for password encryption
func (instance *CodingHelper) EncryptTextWithPrefix(text string, key []byte) (string, error) {
	if !strings.HasPrefix(text, "enc-") {
		data, err := instance.EncryptBytesAES([]byte(text), Strings.FillLeftBytes(key, 32, '0'))
		if nil != err {
			return "", err
		}
		return "enc-" + instance.EncodeBase64(data), nil
	}
	return text, nil
}

func (instance *CodingHelper) DecryptTextWithPrefix(text string, key []byte) (string, error) {
	if strings.HasPrefix(text, "enc-") {
		text = text[4:]
		data, err := instance.DecodeBase64(text)
		if nil != err {
			return "", err
		}
		data, err = instance.DecryptBytesAES(data, Strings.FillLeftBytes(key, 32, '0'))
		if nil != err {
			return "", err
		}
		return string(data), nil
	}
	return text, nil
}

func (instance *CodingHelper) EncryptTextAES(text string, key []byte) ([]byte, error) {
	return instance.EncryptBytesAES([]byte(text), key)
}

func (instance *CodingHelper) EncryptFileAES(fileName string, key []byte, optOutFileName string) ([]byte, error) {
	data, err := IO.ReadBytesFromFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	encoded, err := instance.EncryptBytesAES(data, key)
	if err != nil {
		return []byte{}, err
	}

	// write file
	if len(optOutFileName) > 0 {
		_, err := IO.WriteBytesToFile(encoded, optOutFileName)
		if err != nil {
			return []byte{}, err
		}
	} else {
		_, err := IO.WriteBytesToFile(encoded, fileName)
		if err != nil {
			return []byte{}, err
		}
	}
	return encoded, nil
}

func (instance *CodingHelper) EncryptBytesAES(data []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key) // key must be 32 bytes
	// if there are any errors, handle them
	if err != nil {
		return []byte{}, err
	}

	// gcm or Galois/Counter Mode, is a mode of operation
	// for symmetric key cryptographic block ciphers
	// - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	gcm, err := cipher.NewGCM(c)
	// if any error generating new GCM
	// handle them
	if err != nil {
		return []byte{}, err
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, gcm.NonceSize())
	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return []byte{}, err
	}

	return gcm.Seal(nonce, nonce, data, nil), nil
}

func (instance *CodingHelper) DecryptTextAES(text string, key []byte) ([]byte, error) {
	return instance.DecryptBytesAES([]byte(text), key)
}

func (instance *CodingHelper) DecryptFileAES(fileName string, key []byte, optOutFileName string) ([]byte, error) {
	data, err := IO.ReadBytesFromFile(fileName)
	if err != nil {
		return []byte{}, err
	}
	encoded, err := instance.DecryptBytesAES(data, key)
	if err != nil {
		return []byte{}, err
	}

	// write file
	if len(optOutFileName) > 0 {
		_, err := IO.WriteBytesToFile(encoded, optOutFileName)
		if err != nil {
			return []byte{}, err
		}
	} else {
		_, err := IO.WriteBytesToFile(encoded, fileName)
		if err != nil {
			return []byte{}, err
		}
	}
	return encoded, nil
}

func (instance *CodingHelper) DecryptBytesAES(data []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return []byte{}, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return []byte{}, err
	}

	nonce, data := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return []byte{}, err
	}

	return plain, nil
}
