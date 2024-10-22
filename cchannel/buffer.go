package cchannel

type Buffer struct {
	maxSize int
	buffer  []interface{}
}

func NewBuffer(size int) *Buffer {
	return &Buffer{
		maxSize: size,
		buffer:  make([]interface{}, 0, size),
	}
}

func (b *Buffer) IsEmpty() bool {
	return len(b.buffer) == 0
}

func (b *Buffer) IsFull() bool {
	return len(b.buffer) == b.maxSize
}

func (b *Buffer) Enqueue(value interface{}) {
	b.buffer = append(b.buffer, value)
}

func (b *Buffer) Dequeue() interface{} {
	if b.IsEmpty() {
		return nil
	}

	value := b.buffer[0]
	b.buffer = b.buffer[1:]
	return value
}
