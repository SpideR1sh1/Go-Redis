package main

import (
	"bufio"
	"io"
	"os"
	"sync"
	"time"
)

type Aof struct {
	file *os.File
	rd   *bufio.Reader
	mu   sync.Mutex
}

// NewAof initializes the AOF system and starts background syncing
func NewAof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: f,
		rd:   bufio.NewReader(f),
	}

	// Start a background goroutine for syncing the AOF file periodically
	go func() {
		for {
			aof.mu.Lock()
			aof.file.Sync() // Sync changes to disk
			aof.mu.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

// Close ensures the AOF file is properly closed
func (aof *Aof) Close() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Close()
}

// Write appends a marshaled value to the AOF file
func (aof *Aof) Write(value Value) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

// Read replays the AOF log by reading all saved commands and applying a handler function
func (aof *Aof) Read(fn func(value Value)) error {
	aof.mu.Lock()
	defer aof.mu.Unlock()

	// Reset the file pointer to the beginning
	_, err := aof.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	reader := NewResp(aof.file)

	for {
		value, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		fn(value)
	}

	return nil
}

// Flush manually triggers syncing of data to disk
func (aof *Aof) Flush() error {
	aof.mu.Lock()
	defer aof.mu.Unlock()
	return aof.file.Sync()
}
