package cchannel

import (
	"sync"
)

const (
	Bidirectional = iota
	SendOnly
	ReceiveOnly
)

type CChannel struct {
	buf    Buffer
	cond   sync.Cond
	sendQ  WaitQueue
	recvQ  WaitQueue
	closed bool
	chType int
}

func NewChannel(size int, chType int) *CChannel {
	return &CChannel{
		buf:    *NewBuffer(size),
		cond:   *sync.NewCond(&sync.Mutex{}),
		sendQ:  *NewWQ(),
		recvQ:  *NewWQ(),
		chType: chType,
	}
}

func (c *CChannel) Send(value interface{}) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	if c.closed {
		panic("send on closed channel")
	}

	if c.chType == ReceiveOnly {
		panic("compilation error")
	}

	if c.buf.IsFull() || c == nil || c.buf.maxSize == 0 {
		c.sendQ.EQ(value)
	}
	c.buf.Enqueue(value)
}

func (c *CChannel) Receive() (interface{}, bool) {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	if c.chType == SendOnly {
		panic("compilation error")
	}

	if c.buf.IsEmpty() {
		if c.closed {
			return nil, false
		}
		// unbuffered
		if !c.sendQ.IsEmpty() {
			if !c.recvQ.IsEmpty() {
				c.recvQ.DQ()
			}
			return c.sendQ.DQ(), true
		}
		c.recvQ.EQ(nil)
	}

	if c.buf.IsFull() {
		if !c.recvQ.IsEmpty() {
			c.recvQ.DQ()
		}
		c.buf.Enqueue(c.sendQ.DQ())
		return c.buf.Dequeue(), true
	}

	return c.buf.Dequeue(), true
}

func (c *CChannel) Close() {
	if c == nil {
		panic("close of nil channel")
	}

	if c.closed {
		panic("close of closed channel")
	}

	if c.chType == ReceiveOnly {
		panic("compilation error")
	}

	c.cond.L.Lock()
	defer c.cond.L.Unlock()

	c.closed = true
	//c.cond.Broadcast()
}
