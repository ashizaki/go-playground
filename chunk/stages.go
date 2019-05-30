package main

import "sync"

func generator(done <-chan interface{}, intergers ...int) <-chan int {
	intStream := make(chan int, len(intergers))
	go func() {
		defer close(intStream)
		for _, i := range intergers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}

func multiply(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
	multipledStream := make(chan int)
	go func() {
		defer close(multipledStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipledStream <- i * multiplier:
			}
		}

	}()
	return multipledStream
}

func add(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case addedStream <- i * additive:
			}
		}

	}()
	return addedStream
}

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case stream <- v:
				}
			}
		}
	}()
	return stream
}

func take(done <-chan interface{}, inStream <-chan interface{}, num int) <-chan interface{} {
	outStream := make(chan interface{})
	go func() {
		defer close(outStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case outStream <- <-inStream:
			}
		}

	}()
	return outStream
}

func repeatFn(
	done <-chan interface{},
	fn func() interface{},
) <-chan interface{} {
	stream := make(chan interface{})
	go func() {
		defer close(stream)
		for {
			select {
			case <-done:
				return
			case stream <- fn():
			}
		}
	}()
	return stream
}

func toString(done <-chan interface{}, in <-chan interface{}) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v.(string):
			}
		}
	}()
	return out
}

func toInt(done <-chan interface{}, in <-chan interface{}) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for v := range in {
			select {
			case <-done:
				return
			case out <- v.(int):
			}
		}
	}()
	return out
}

func repeatString(done <-chan interface{}, values ...string) <-chan string {
	stream := make(chan string)
	go func() {
		defer close(stream)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case stream <- v:
				}
			}
		}
	}()
	return stream
}

func takeString(done <-chan interface{}, inStream <-chan string, num int) <-chan string {
	outStream := make(chan string)
	go func() {
		defer close(outStream)
		for i := 0; i < num; i++ {
			select {
			case <-done:
				return
			case outStream <- <-inStream:
			}
		}

	}()
	return outStream
}

func primeFinder(done <-chan interface{}, intStream <-chan int) <-chan interface{} {
	primeStream := make(chan interface{})
	go func() {
		defer close(primeStream)
		for integer := range intStream {
			integer -= 1
			prime := true
			for divisor := integer - 1; divisor > 1; divisor-- {
				if integer%divisor == 0 {
					prime = false
					break
				}
			}

			if prime {
				select {
				case <-done:
					return
				case primeStream <- integer:
				}
			}
		}
	}()
	return primeStream
}

func fanIn(done <-chan interface{}, channels ...<-chan interface{}) <-chan interface{} {
	var wg sync.WaitGroup
	out := make(chan interface{})
	multiplex := func(c <-chan interface{}) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case out <- i:
			}
		}
	}

	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
