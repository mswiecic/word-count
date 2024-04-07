package word_count

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
)

type Mode int

const (
	COUNT = iota
	FRQ
)

type Counter interface {
	SetChunkSize(size int)
	SetMode(mode Mode)
	ProcessFile(fileName string)
	ProcessChunk(chunk *[]byte)
	GetCount() int
	GetFirstNWords(n int) []pair
}
type counter struct {
	words     sync.Map
	file      string
	counter   int
	chunkSize int
	mode      Mode
	sync.RWMutex
}

// NewCounter returns a Counter object with initialized chunk size to 2048 and in mode COUNT words
func NewCounter() Counter {
	return &counter{
		counter:   0,
		chunkSize: 2048,
		mode:      COUNT,
	}
}

func (c *counter) SetChunkSize(size int) {
	c.chunkSize = size
}

func (c *counter) SetMode(mode Mode) {
	c.mode = mode
}

type pair struct {
	key   string
	value uint32
}

func (c *counter) increment(number int) {
	c.Lock()
	defer c.Unlock()
	c.counter += number
}

func (c *counter) append(key string) {
	var one uint32 = 1
	val, ok := c.words.LoadOrStore(key, &one)
	if ok {
		atomic.AddUint32(val.(*uint32), 1)
	}
}

// GetCount returns number of counted word
func (c *counter) GetCount() int {
	return c.counter
}

// GetFirstNWords return provided first most frequent N words in Counter ordered
// from the most popular to the least popular
func (c *counter) GetFirstNWords(n int) []pair {
	kvs := make([]pair, 0)
	c.words.Range(func(key, value any) bool {
		kvs = append(kvs, pair{key: key.(string), value: *value.(*uint32)})
		return true
	})
	sort.Slice(kvs, func(i, j int) bool {
		return kvs[i].value > kvs[j].value
	})
	if len(kvs) > n {
		return kvs[:n]
	} else {
		return kvs
	}
}

// ProcessFile starts working with a file
func (c *counter) ProcessFile(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
	}(file)

	r := bufio.NewReader(file)

	chunkPool := sync.Pool{
		New: func() interface{} {
			c := make([]byte, c.chunkSize)
			return &c
		},
	}

	c.counter = 0

	var wg sync.WaitGroup
	for {
		// Read the next chunk
		chunk := chunkPool.Get().(*[]byte)
		n, err := r.Read(*chunk)
		*chunk = (*chunk)[:n]

		if err != nil && err != io.EOF {
			fmt.Println(err)
		}

		// Check for EOF
		if err == io.EOF {
			break
		}

		bytes, err := r.ReadBytes('\n') // read entire line
		if err != nil {
			fmt.Println(err)
		}
		if len(bytes) > 0 {
			*chunk = append(*chunk, bytes...)
		}

		// Print the chunk
		// fmt.Printf("%s", chunk)
		wg.Add(1)
		go func(chunk *[]byte) {
			defer func() {
				wg.Done()
				chunkPool.Put(chunk)
			}()
			c.ProcessChunk(chunk)
		}(chunk)
	}
	wg.Wait()
}

// ProcessChunk count word or add to a map depends on Mode
func (c *counter) ProcessChunk(chunk *[]byte) {
	s := string(*chunk)
	fields := strings.Fields(s)
	switch c.mode {
	case COUNT:
		c.increment(len(fields))
	case FRQ:
		for _, field := range fields {
			c.append(field)
		}
	}
}
