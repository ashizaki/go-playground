package main

import "fmt"
import "math/rand"
import "time"

func main() {
	test2()
}

func test1() {
	complete := make(chan interface{})
	doWork := func(
		done <-chan time.Time,
		strings <-chan string,
	) <-chan interface{} {
		go func() {
			defer fmt.Println("doWork exited.")
			defer close(complete)
			for {
				select {
				case s := <-strings:
					fmt.Println(s)
				case <-done:
					return
				}
			}
		}()
		return complete
	}

	term := doWork(time.After(1*time.Second), nil)

	<-term
	fmt.Println("Done.")
}

func test2() {
	newRandStream := func(done <-chan interface{}) <-chan int {
		randStream := make(chan int)
		go func() {
			defer fmt.Println("closure exited.")
			defer close(randStream)
			for {
				select {
				case randStream <- rand.Int():
				case <-done:
					return
				}
			}
		}()
		return randStream
	}

	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	time.Sleep(1 * time.Second)
}
