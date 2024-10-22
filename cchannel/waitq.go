package cchannel

type WaitQueue struct {
	queue []interface{}
}

func NewWQ() *WaitQueue {
	return &WaitQueue{
		queue: make([]interface{}, 0),
	}
}

func (wq *WaitQueue) IsEmpty() bool {
	return len(wq.queue) == 0
}

func (wq *WaitQueue) EQ(value interface{}) {
	wq.queue = append(wq.queue, value)
}

func (wq *WaitQueue) DQ() interface{} {
	if wq.IsEmpty() {
		return nil
	}
	value := wq.queue[0]
	wq.queue = wq.queue[1:]
	return value
}
