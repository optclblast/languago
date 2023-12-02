package cache

import (
	"unsafe"
)

type Cache interface {
	Add(key string, value any) error
	Delete(key string) error
	Flush() error
	MemoryLimit(limit int)
}

type inmemory struct {
	// negative limit == "unlimited" storage size
	memoryLimit int
	storage     map[string]any
}

func (c *inmemory) Add(key string, value any) error {
	if unsafe.Sizeof(c.storage) >= uintptr(c.memoryLimit) &&
		c.memoryLimit > 0 {
		go func() {
			for k := range c.storage {
				delete(c.storage, k)
			}
		}()
	}

	c.storage[key] = value
	return nil
}

func (c *inmemory) Delete(key string) error {
	delete(c.storage, key)
	return nil
}

func (c *inmemory) Flush() error {
	clear(c.storage)
	return nil
}

func (c *inmemory) MemoryLimit(limit int) {
	c.memoryLimit = limit
}

// todo
type redis struct{}
