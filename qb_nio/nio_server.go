package qb_nio

import (
	"bufio"
	"crypto/rsa"
	"encoding/gob"
	"fmt"
	"net"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NioServer struct {

	//-- private --//
	uuid       string
	port       int
	listener   net.Listener
	clients    int
	clientsMap map[string]*client
	mux        sync.Mutex
	handler    NioMessageHandler
	stopChan   chan bool
	active     bool
	// RSA
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

type client struct {
	Id         string
	publicKey  *rsa.PublicKey // public key of the client
	sessionKey []byte         // session key for the client
}

type NioMessageHandler func(message *NioMessage) interface{}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioServer) GetUUID() string {
	if nil != instance {
		return instance.uuid
	}
	return ""
}

func (instance *NioServer) IsOpen() bool {
	if nil != instance {
		return instance.active
	}
	return false
}

func (instance *NioServer) Port() int {
	if nil != instance {
		return instance.port
	}
	return 0
}

func (instance *NioServer) Open() error {
	if nil != instance {
		if !instance.active {
			instance.active = true
			instance.stopChan = make(chan bool, 1)

			err := instance.initRSA()
			if nil != err {
				return err
			}

			listener, err := net.Listen("tcp", fmt.Sprintf(":%v", instance.port))
			if nil != err {
				return err
			}
			instance.listener = listener

			// main listener loop
			go instance.open()
		}
	}
	return nil
}

func (instance *NioServer) Close() error {
	if nil != instance {
		if instance.active {
			instance.active = false
			var err error
			if nil != instance.listener {
				err = instance.listener.Close()
			}
			instance.stopChan <- true
			return err
		}
	}
	return nil
}

// Wait is stopped
func (instance *NioServer) Join() {
	// locks and wait for exit response
	<-instance.stopChan
}

func (instance *NioServer) ClientsCount() int {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		return instance.clients
	}
	return 0
}

func (instance *NioServer) ClientsId() []string {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()
		keys := make([]string, 0, len(instance.clientsMap))
		for k := range instance.clientsMap {
			keys = append(keys, k)
		}
		return keys
	}
	return []string{}
}

func (instance *NioServer) OnMessage(callback NioMessageHandler) {
	if nil != instance {
		instance.handler = callback
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioServer) initRSA() error {
	if nil != instance && nil == instance.privateKey {
		// TODO: implement loading from file

		// auto-generates
		pri, pub, err := keysGenerate(KEY_LEN)
		if nil != err {
			return err
		}
		instance.privateKey = pri
		instance.publicKey = pub
	}
	return nil
}

func (instance *NioServer) open() {
	for {
		// accept connections
		conn, err := instance.listener.Accept()
		if err != nil {
			// error accepting connection
			continue
		}
		go instance.handleConnection(conn)
	}
}

func (instance *NioServer) incClients(conn net.Conn) *client {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		session := newSessionKey()

		c := new(client)
		c.Id = conn.RemoteAddr().String()
		c.sessionKey = session[:]
		instance.clients++
		instance.clientsMap[c.Id] = c
		return c
	}
	return nil
}

func (instance *NioServer) decClients(id string) {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		_, ok := instance.clientsMap[id]
		if ok {
			instance.clients--
			delete(instance.clientsMap, id)
		}
	}
}

func (instance *NioServer) handleConnection(conn net.Conn) {
	rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
	defer conn.Close()

	// new client connection
	c := instance.incClients(conn)

	// connection loop
	for {
		var message NioMessage
		dec := gob.NewDecoder(rw)
		err := dec.Decode(&message)
		if nil != err {
			if err.Error() == "EOF" {
				// client disconnected
			}
			// exit
			break
		}

		isHandshake := instance.isHandshake(&message)
		if isHandshake {
			// set public key for client
			c.publicKey = message.PublicKey
		}

		if !isHandshake && nil != instance.handler {
			if nil != c.publicKey && nil != c.sessionKey {
				// decode client message body
				if v, b := message.Body.([]byte); b {
					data, err := decrypt(v, c.sessionKey)
					if nil == err {
						message.Body = data
					} else {
						// encryption error
						fmt.Println("Http error decrypting data:", err)
					}
				}
			}
			customResponse := instance.handler(&message)
			if nil == customResponse {
				customResponse = true
			}

			err := sendResponse(customResponse, rw, c.sessionKey, c.publicKey, instance.publicKey, isHandshake)
			if err != nil {
				break
			}
		} else {
			// response OK (default)
			err := sendResponse(true, rw, c.sessionKey, c.publicKey, instance.publicKey, isHandshake)
			if err != nil {
				break
			}
		}
	}

	// client removed
	instance.decClients(c.Id)
}

func (instance *NioServer) isHandshake(message *NioMessage) bool {
	if v, b := message.Body.([]byte); b {
		return string(v) == string(HANDSHAKE.Body.([]byte))
	}
	return false
}

func sendResponse(body interface{}, rw *bufio.ReadWriter, sessionKey []byte, clientKey, serverKey *rsa.PublicKey, isHandshake bool) error {
	response := new(NioMessage)

	// public key is passed only with handshake
	if isHandshake {
		response.PublicKey = serverKey
		if nil != clientKey {
			response.SessionKey, _ = encryptKey(sessionKey, clientKey)
		}
	}

	s := serialize(body)

	// encode server message body
	if nil != clientKey && !isHandshake {
		data, err := encrypt(s, sessionKey)
		if nil == err {
			s = data
		} else {
			fmt.Println("Http error encrypting data", err)
		}
	}
	response.Body = s

	enc := gob.NewEncoder(rw)
	err := enc.Encode(response)
	if err != nil {
		return err
	}
	err = rw.Flush()
	return err
}
