package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func BenchmarkUnbufferdWrite(b *testing.B) {
	performWrite(b, tmpFileOrFatal())
}

func BenchmarkBufferdWrite(b *testing.B) {
	bufferedFile := bufio.NewWriter(tmpFileOrFatal())
	performWrite(b, bufferedFile)
}

func tmpFileOrFatal() *os.File {
	file, err := ioutil.TempFile("", "tmp")
	if err != nil {
		log.Fatalf("error :%v", err)
	}
	return file
}

func performWrite(b *testing.B, writer io.Writer) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for bt := range take(done, repeat(done, byte(0)), b.N) {
		writer.Write([]byte{bt.(byte)})
	}
}
