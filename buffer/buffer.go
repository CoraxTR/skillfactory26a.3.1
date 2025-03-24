package buffer

import (
	"fmt"
	"log"
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
	log.Printf("[INFO] New buffer created")
	return &buffer{
		storage:      make([]int, bufferSize),
		lifetime:     bufferLifeTime,
		writePointer: 0,
		readPointer:  0,
		isEmpty:      true,
	}
}

func (b *buffer) WriteToBuffer(data int) {
	log.Printf("[INFO] Writing %d to buffer", data)
	b.storage[b.writePointer] = data
	b.writePointer++
	b.isEmpty = false
	if b.writePointer == bufferSize {
		log.Printf("[WARNING] Buffer overflow")
		b.writePointer = 0
	}
}

func (b *buffer) ReadAllFromBuffer() ([]int, error) {
	if b.isEmpty == true {
		log.Printf("[WARNING] Trying to read empty buffer")
		return nil, fmt.Errorf("Buffer is empty")
	}
	data := b.storage[b.readPointer:b.writePointer]
	log.Printf("[INFO] Reading %d from buffer", data)
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
