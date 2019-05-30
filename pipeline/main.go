package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	fmt.Println("vim-go")
	test7()
}

func test1() {
	multiply := func(values []int, multiplier int) []int {
		multipledValues := make([]int, len(values))
		for i, v := range values {
			multipledValues[i] = v * multiplier
		}
		return multipledValues
	}

	add := func(values []int, additive int) []int {
		addedValue := make([]int, len(values))
		for i, v := range values {
			addedValue[i] = v + additive
		}
		return addedValue
	}

	ints := []int{1, 2, 3, 4}
	for _, v := range add(multiply(ints, 2), 1) {
		fmt.Println(v)
	}
}

func test2() {
	generator := func(done <-chan interface{}, intergers ...int) <-chan int {
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

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
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

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
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

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1, 2, 3, 4, 5, 6, 7)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

func test3() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
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

	take := func(done <-chan interface{}, inStream <-chan interface{}, num int) <-chan interface{} {
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

	done := make(chan interface{})
	defer close(done)

	for num := range take(done, repeat(done, 1), 10) {
		fmt.Printf("%v ", num)
	}
}

func test4() {
	repeatFn := func(
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

	take := func(done <-chan interface{}, inStream <-chan interface{}, num int) <-chan interface{} {
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

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Int() }
	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Printf("%v\n", num)
	}
}

func test5() {
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
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

	take := func(done <-chan interface{}, inStream <-chan interface{}, num int) <-chan interface{} {
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

	toString := func(done <-chan interface{}, in <-chan interface{}) <-chan string {
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

	done := make(chan interface{})
	defer close(done)

	var message string
	for token := range toString(done, take(done, repeat(done, "I", "am."), 5)) {
		message += token
	}
	fmt.Printf("message: %s...", message)
}

func test6() {
	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	stream := toInt(done, repeatFn(done, rand))
	fmt.Println("Primes:")
	for prime := range take(done, primeFinder(done, stream), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))
}

func test7() {
	rand := func() interface{} { return rand.Intn(50000000) }

	done := make(chan interface{})
	defer close(done)

	start := time.Now()

	randIntStream := toInt(done, repeatFn(done, rand))

	numFinder := runtime.NumCPU()
	fmt.Printf("Spinning up %d prime finders.\n", numFinder)
	finders := make([]<-chan interface{}, numFinder)

	fmt.Println("Primes:")
	for i := 0; i < numFinder; i++ {
		finders[i] = primeFinder(done, randIntStream)
	}

	for prime := range take(done, fanIn(done, finders...), 10) {
		fmt.Printf("\t%d\n", prime)
	}

	fmt.Printf("Search took: %v", time.Since(start))

}
