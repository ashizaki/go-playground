package main

import "testing"

func BenchmarkGeneric(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range toString(done, take(done, repeat(done, "a"), b.N)) {
	}
}

func BenchmarkTyped(b *testing.B) {
	done := make(chan interface{})
	defer close(done)

	b.ResetTimer()
	for range takeString(done, repeatString(done, "a"), b.N) {
	}
}
