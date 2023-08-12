package qb_nio

import (
	"bufio"
	"crypto/rsa"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/rskvp/qb-core/qb_events"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NioClient struct {
	Timeout    time.Duration
	Secure     bool
	EnablePing bool

	//-- private --//
	uuid      string
	conn      net.Conn
	host      string
	port      int
	mux       sync.Mutex
	stopChan  chan bool
	events    *qb_events.Emitter
	connected bool
	closed    bool
	pingTimer *time.Ticker
	// RSA
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	serverKey  *rsa.PublicKey // server signature (got on handshake)
	sessionKey []byte
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioClient) GetUUID() string {
	if nil != instance {
		return instance.uuid
	}
	return ""
}

func (instance *NioClient) IsOpen() bool {
	if nil != instance {
		if !instance.closed {
			err := instance.test() // test for real feedback
			return nil == err
		}
	}
	return false
}

func (instance *NioClient) Open() error {
	if nil != instance {
		if instance.closed {
			instance.closed = false
			instance.stopChan = make(chan bool, 1)

			// start pinging remote server every 1 second
			go instance.startLoop()

			err := instance.initRSA()
			if nil != err {
				return err
			}
			return instance.handshake()
		}
	}
	return nil
}

func (instance *NioClient) Close() error {
	if nil != instance {
		if nil != instance.conn {
			instance.closed = true
			err := instance.conn.Close()

			// emit event
			instance.setConnected(false)

			return err
		}
		instance.stopChan <- true
	}
	return nil
}

// Join Wait is stopped
func (instance *NioClient) Join() {
	// locks and wait for exit response
	<-instance.stopChan
}

func (instance *NioClient) Send(data interface{}) (*NioMessage, error) {
	if nil != instance {

		// creates NIO message
		message := new(NioMessage)
		message.Body = data

		return instance.send(message, false)
	}
	return nil, nil
}

func (instance *NioClient) OnConnect(callback func(e *qb_events.Event)) {
	if nil != instance {
		instance.events.On("connect", callback)
	}
}

func (instance *NioClient) OffConnect() {
	if nil != instance {
		instance.events.Off("connect")
	}
}

func (instance *NioClient) OnDisconnect(callback func(e *qb_events.Event)) {
	if nil != instance {
		instance.events.On("disconnect", callback)
	}
}

func (instance *NioClient) OffDisconnect() {
	if nil != instance {
		instance.events.Off("disconnect")
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *NioClient) initRSA() error {
	if nil != instance && instance.Secure && nil == instance.privateKey {
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

func (instance *NioClient) setConnected(status bool) {
	if nil != instance {
		if status {
			if !instance.connected {
				instance.events.EmitAsync("connect")
			}
		} else {
			if instance.connected {
				instance.events.EmitAsync("disconnect")
			}
			// reset connection for next call to regenerate
			if nil != instance.conn {
				_ = instance.conn.Close()
				instance.conn = nil
			}
		}
		instance.connected = status
	}
}

func (instance *NioClient) test() error {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", instance.host, instance.port), instance.Timeout)
		if nil != conn {
			err = conn.Close()
		}
		return err
	}
	return nil
}

func (instance *NioClient) connect() (net.Conn, error) {
	if nil != instance {
		if nil == instance.conn {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", instance.host, instance.port), instance.Timeout)
			if nil == err {
				instance.conn = conn
				// trigger connect
				instance.setConnected(true)
			} else {
				// trigger disconnect
				instance.setConnected(false)
			}
			return conn, err
		}
		return instance.conn, nil
	}
	return nil, nil
}

func (instance *NioClient) handshake() error {
	if nil != instance {
		HANDSHAKE.PublicKey = instance.publicKey
		response, err := instance.send(HANDSHAKE, true)
		if nil != err {
			return err
		}
		instance.serverKey = response.PublicKey
	}
	return nil
}

func (instance *NioClient) ping() error {
	if nil != instance {
		err := instance.test()
		instance.setConnected(nil == err)
		if nil != err {
			return err
		}
	}
	return nil
}

func (instance *NioClient) startLoop() {

	for {
		if instance.closed {
			return
		}
		if nil != instance && nil != instance.pingTimer {
			select {
			case <-instance.stopChan:
				return
			case <-instance.pingTimer.C:
				// event
				if nil != instance && !instance.closed && instance.EnablePing {
					_ = instance.ping()
				}
			}
		}
	}

}

func (instance *NioClient) send(message *NioMessage, handshake bool) (*NioMessage, error) {
	if nil != instance {
		instance.mux.Lock()
		defer instance.mux.Unlock()

		conn, err := instance.connect()
		if nil != err {
			_ = instance.Close() // reset connection
			return nil, err
		}

		if handshake {
			message.PublicKey = instance.publicKey
		}

		// ENCRYPT BODY
		if !handshake && nil != instance.publicKey && len(instance.sessionKey) > 0 {
			v := serialize(message.Body)
			data, err := encrypt(v, instance.sessionKey)
			if nil == err {
				message.Body = data
			} else {
				return nil, errors.New("Client Encryption error")
			}
		} else {
			message.Body = serialize(message.Body)
		}

		rw := bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))
		enc := gob.NewEncoder(rw)
		err = enc.Encode(message)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Encode failed for message: %#v", message))
		}
		err = rw.Flush()
		if err != nil {
			return nil, errors.New("Flush failed.")
		}

		// read NIO response
		var response NioMessage
		dec := gob.NewDecoder(rw)
		err = dec.Decode(&response)
		if err != nil {
			// fmt.Println(errors.Wrap(err, "Client failed to read response"))
			return nil, errors.New("Client failed to read response")
		} else {
			// RESPONSE FROM SERVER
			if !handshake {
				if len(instance.sessionKey) > 0 {
					// DECRYPT BODY
					if v, b := response.Body.([]byte); b {
						data, err := decrypt(v, instance.sessionKey)
						if nil == err {
							response.Body = data
						}
					}
				}
			} else {
				// handshake
				if len(response.SessionKey) > 0 {
					data, err := decryptKey(response.SessionKey, instance.privateKey)
					if nil == err {
						instance.sessionKey = data
					}
				}
			}
			return &response, nil
		}
	}
	return nil, nil
}
