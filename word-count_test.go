package word_count

import (
	"fmt"
	"strings"
	"testing"
)

func TestProcessFile_FileNotExist(t *testing.T) {
	c := NewCounter()
	c.SetChunkSize(4)
	c.ProcessFile("not_existing")
	if c.GetCount() != 0 {
		t.Fail()
	}
	if len(c.GetFirstNWords(10)) != 0 {
		t.Fail()
	}
}

func TestProcessFileCount(t *testing.T) {
	c := NewCounter()
	c.SetChunkSize(8 * 1024)
	c.ProcessFile("pan50.txt")
	if c.GetCount() == 0 {
		t.Fail()
	}
	if len(c.GetFirstNWords(10)) != 0 {
		t.Fail()
	}
}

func TestProcessFileWords(t *testing.T) {
	c := NewCounter()
	c.SetChunkSize(8 * 1024)
	c.SetMode(FRQ)
	c.ProcessFile("pan50.txt")
	count := c.GetCount()
	if count != 0 {
		t.Fail()
	}
	words := c.GetFirstNWords(10)
	if len(words) == 0 {
		t.Fail()
	}
}

func TestProcessChunkCount(t *testing.T) {
	s := "1 2 3 4 5"
	c := NewCounter()
	chunk := []byte(s)
	c.ProcessChunk(&chunk)
	if c.GetCount() != 5 {
		t.Fail()
	}
}

func TestProcessChunkWords(t *testing.T) {
	s := "ala ma kota ale nie ma psa"
	c := NewCounter()
	c.SetMode(FRQ)
	chunk := []byte(s)
	c.ProcessChunk(&chunk)
	first := c.GetFirstNWords(1)[0]
	if first.key != "ma" {
		t.Fatalf("there is no value")
	}
	if first.value != 2 {
		t.Fatalf("value is not 2")
	}
}

var inputParams = []struct {
	file string
	size int
}{
	{file: "test.txt", size: 2 * 1024},
	{file: "test.txt", size: 4 * 1024},
	{file: "test.txt", size: 8 * 1024},
	{file: "pan50.txt", size: 2 * 1024},
	{file: "pan50.txt", size: 4 * 1024},
	{file: "pan50.txt", size: 8 * 1024},
	{file: "large.txt", size: 2 * 1024},
	{file: "large.txt", size: 4 * 1024},
	{file: "large.txt", size: 8 * 1024},
}

func BenchmarkCountWords(b *testing.B) {
	for _, param := range inputParams {
		b.Run(fmt.Sprintf("%s_%d", strings.Split(param.file, ".")[0], param.size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := NewCounter()
				c.SetMode(COUNT)
				c.SetChunkSize(param.size)
				c.ProcessFile(param.file)
			}
		})
	}
}

func BenchmarkWordsFrequency(b *testing.B) {
	for _, param := range inputParams {
		b.Run(fmt.Sprintf("%s_%d", strings.Split(param.file, ".")[0], param.size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := NewCounter()
				c.SetMode(FRQ)
				c.SetChunkSize(param.size)
				c.ProcessFile(param.file)
			}
		})
	}
}

func BenchmarkWordsFrequencyWithGettingNWords(b *testing.B) {
	for _, param := range inputParams {
		b.Run(fmt.Sprintf("%s_%d", strings.Split(param.file, ".")[0], param.size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c := NewCounter()
				c.SetMode(FRQ)
				c.SetChunkSize(param.size)
				c.ProcessFile(param.file)
				c.GetFirstNWords(10)
			}
		})
	}
}
