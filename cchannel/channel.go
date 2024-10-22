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

	if (c.buf.IsFull() || c.buf.maxSize == 0) && c.recvQ.IsEmpty() {
		c.sendQ.EQ(value)
		c.cond.Wait()
		return
	}
	c.buf.Enqueue(value)
	c.cond.Signal()
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
		c.recvQ.EQ(nil)
		c.cond.Wait()
		// unbuffered
		if !c.sendQ.IsEmpty() && c.buf.maxSize == 0 {
			if !c.recvQ.IsEmpty() {
				c.recvQ.DQ()
			}
			c.cond.Signal()
			return c.sendQ.DQ(), true
		}
	}

	if c.buf.IsFull() {
		value := c.buf.Dequeue()
		c.buf.Enqueue(c.sendQ.DQ())
		c.cond.Signal()
		return value, true
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
