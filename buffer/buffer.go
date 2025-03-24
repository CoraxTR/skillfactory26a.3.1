package buffer

import (
	"fmt"
	"time"
)

const bufferSize int = 1024
const bufferLifeTime time.Duration = 5 * time.Second

type buffer struct {
	storage      []int
	lifetime     time.Duration
	writePointer int
	readPointer  int
	isEmpty      bool
}

func CreateBuffer() *buffer {
	return &buffer{
		storage:      make([]int, bufferSize),
		lifetime:     bufferLifeTime,
		writePointer: 0,
		readPointer:  0,
		isEmpty:      true,
	}
}

func (b *buffer) WriteToBuffer(data int) {
	b.storage[b.writePointer] = data
	b.writePointer++
	b.isEmpty = false
	if b.writePointer == bufferSize {
		fmt.Println("BufferOverFlow")
		b.writePointer = 0
	}
}

func (b *buffer) ReadAllFromBuffer() ([]int, error) {
	if b.isEmpty == true {
		return nil, fmt.Errorf("Buffer is empty")
	}
	data := b.storage[b.readPointer:b.writePointer]
	b.ClearBuffer()
	return data, nil
}

func (b *buffer) ClearBuffer() {
	b.writePointer = 0
	b.readPointer = 0
	b.isEmpty = true
}

func (b *buffer) GetBufferLifetime() time.Duration {
	return b.lifetime
}
