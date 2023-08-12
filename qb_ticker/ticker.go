package qb_ticker

import (
	"sync"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Ticker struct {
	timer     *time.Ticker
	timeout   time.Duration
	stopChan  chan bool
	paused    bool
	callback  TickerCallback
	locked    bool
	running   bool
	lockMux   sync.Mutex
	statusMux sync.Mutex
}

type TickerCallback func(*Ticker)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewTicker(timeout time.Duration, callback TickerCallback) *Ticker {
	instance := &Ticker{
		timer:    time.NewTicker(timeout),
		timeout:  timeout,
		stopChan: make(chan bool, 1),
		callback: callback,
	}
	instance.running = false

	return instance
}

// Tick Simple ticker loop
func Tick(timeout time.Duration, callback TickerCallback) *Ticker {
	et := NewTicker(timeout, callback)
	et.Start()

	return et
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *Ticker) IsRunning() bool {
	if nil != instance {
		instance.statusMux.Lock()
		defer instance.statusMux.Unlock()
		if nil != instance.timer {
			return instance.running
		}
	}
	return false
}

func (instance *Ticker) IsPaused() bool {
	if nil != instance {
		instance.statusMux.Lock()
		defer instance.statusMux.Unlock()
		if nil != instance.timer {
			return instance.paused
		}
	}
	return false
}

// Join Wait Ticker is stopped
func (instance *Ticker) Join() {
	// locks and wait for exit response
	<-instance.stopChan
}

// Start .... Start the timer
func (instance *Ticker) Start() {
	if nil != instance {
		instance.statusMux.Lock()
		defer instance.statusMux.Unlock()
		instance.start()
	}
}

// Stop ... stops the timer
func (instance *Ticker) Stop() {
	if nil != instance {
		instance.statusMux.Lock()
		defer instance.statusMux.Unlock()
		instance.stop()
	}
}

func (instance *Ticker) Pause() {
	if nil != instance {
		instance.statusMux.Lock()
		defer instance.statusMux.Unlock()
		if nil != instance.timer && !instance.paused {
			instance.paused = true
		}
	}
}

func (instance *Ticker) Resume() {
	if nil != instance {
		instance.statusMux.Lock()
		defer instance.statusMux.Unlock()
		if nil != instance.timer && instance.paused {
			instance.paused = false
		}
	}
}

func (instance *Ticker) Lock() {
	if nil != instance && !instance.locked {
		instance.locked = true
		instance.lockMux.Lock()
	}
}

func (instance *Ticker) Unlock() {
	if nil != instance && instance.locked {
		instance.lockMux.Unlock()
		instance.locked = false
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *Ticker) start() {
	instance.stop()
	instance.timer = time.NewTicker(instance.timeout)
	// infinite loop
	go instance.loop()
}

func (instance *Ticker) stop() {
	if nil != instance.timer {
		instance.timer.Stop()
		instance.timer = nil
	}
	if nil != instance.stopChan {
		instance.stopChan <- true
		instance.stopChan = make(chan bool, 1)
	}
	instance.running = false
}

func (instance *Ticker) loop() {
	if nil != instance && !instance.running {
		instance.running = true
		for {
			if nil != instance && nil != instance.timer {
				select {
				case <-instance.stopChan:
					return
				case <-instance.timer.C:
					// event
					if nil != instance.callback && !instance.paused {
						// thread safe call
						instance.Lock()
						instance.callback(instance)
						instance.Unlock()
					}
				}
			} else {
				return
			}
		}
	}

}
