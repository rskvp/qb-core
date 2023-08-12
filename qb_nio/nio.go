package qb_nio

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/rskvp/qb-core/qb_events"
	"github.com/rskvp/qb-core/qb_rnd"
	"github.com/rskvp/qb-core/qb_sys"
	"github.com/rskvp/qb-core/qb_utils"
)

type NIOHelper struct {
}

var NIO *NIOHelper

func init() {
	NIO = new(NIOHelper)
}

func (*NIOHelper) NewClient(host string, port int) *NioClient {
	instance := new(NioClient)
	instance.host = host
	instance.port = port
	instance.Timeout = 10 * time.Second
	instance.events = qb_events.Events.NewEmitter()
	instance.connected = false
	instance.closed = true
	instance.pingTimer = time.NewTicker(1 * time.Second)
	instance.EnablePing = false // ping disabled (avoid continuous connect/disconnect)

	sysid, err := qb_sys.Sys.ID()
	if nil != err {
		sysid = qb_rnd.Rnd.Uuid()
	}
	instance.uuid = fmt.Sprintf("[%v]:%v", sysid, port)

	return instance
}

func (*NIOHelper) NewServer(port int) *NioServer {
	instance := new(NioServer)
	instance.clients = 0
	instance.port = port
	instance.clientsMap = make(map[string]*client)
	instance.active = false

	sysid, err := qb_sys.Sys.ID()
	if nil != err {
		sysid = qb_rnd.Rnd.Uuid()
	}
	instance.uuid = fmt.Sprintf("[%v]:%v", sysid, port)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

const KEY_LEN = 1024 * 3

var (
	HANDSHAKE = &NioMessage{
		PublicKey: nil,
		Body:      []byte("ACK"),
	}
)

type NioMessage struct {
	PublicKey  *rsa.PublicKey // public key for response
	SessionKey []byte         // session key
	Body       interface{}    // message object
}

type NioSettings struct {
	Address string `json:"address"`
	host    string
	port    int
}

func (instance *NioSettings) Parse(text string) error {
	err := json.Unmarshal([]byte(text), &instance)
	instance.parseAddress(instance.Address)
	return err
}

func (instance *NioSettings) Host() string {
	if instance.port == 0 && len(instance.host) == 0 {
		instance.parseAddress(instance.Address)
	}
	return instance.host
}
func (instance *NioSettings) Port() int {
	if instance.port == 0 && len(instance.host) == 0 {
		instance.parseAddress(instance.Address)
	}
	return instance.port
}
func (instance *NioSettings) parseAddress(address string) {
	tokens := strings.Split(address, ":")
	switch len(tokens) {
	case 1:
		instance.port = qb_utils.Convert.ToInt(tokens[0])
	case 2:
		instance.host = tokens[0]
		instance.port = qb_utils.Convert.ToInt(tokens[1])
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func serialize(data interface{}) []byte {
	if nil != data {
		if v, b := data.([]byte); b {
			return v
		} else if v, b := data.(string); b {
			return []byte(v)
		} else if v, b := data.(error); b {
			data = map[string]interface{}{
				"error": v.Error(),
			}
		}
		return qb_utils.JSON.Bytes(data)
	}
	return []byte{}
}

func keysGenerate(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	return qb_utils.Coding.GenerateKeyPair(bits)
}

func newSessionKey() [32]byte {
	return qb_utils.Coding.GenerateSessionKey()
}

func encryptKey(data []byte, key *rsa.PublicKey) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := qb_utils.Coding.EncryptWithPublicKey(data, key)
		return response, err
	}
	return []byte{}, nil
}

func encrypt(data []byte, key []byte) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := qb_utils.Coding.EncryptBytesAES(data, key)
		return response, err
	}
	return []byte{}, nil
}

func decryptKey(data []byte, privateKey *rsa.PrivateKey) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := qb_utils.Coding.DecryptWithPrivateKey(data, privateKey)
		return response, err
	}
	return []byte{}, nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	if nil != data && len(data) > 0 {
		response, err := qb_utils.Coding.DecryptBytesAES(data, key)
		return response, err
	}
	return []byte{}, nil
}
