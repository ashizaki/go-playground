package main

import (
	"fmt"
	"time"
)

func main() {
	test1()
}

func test1() {
	var or func(chs ...<-chan interface{}) <-chan interface{}
	or = func(chs ...<-chan interface{}) <-chan interface{} {
		switch len(chs) {
		case 0:
			return nil
		case 1:
			return chs[0]
		}

		orDone := make(chan interface{})
		go func() {
			defer close(orDone)
			switch len(chs) {
			case 2:
				select {
				case <-chs[0]:
				case <-chs[1]:
				}
			default:
				select {
				case <-chs[0]:
				case <-chs[1]:
				case <-chs[2]:
				case <-or(append(chs[3:], orDone)...):
				}
			}
		}()
		return orDone
	}

	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Second),
		sig(5*time.Second),
		sig(1*time.Second),
		sig(4*time.Second),
		sig(10*time.Second),
		sig(4*time.Second),
	)
	fmt.Printf("done after %v", time.Since(start))
}
