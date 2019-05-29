package main

import "fmt"
import "sync"

func main() {
	test2()
}

func test1() {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begin\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutine...")
	close(begin)
	wg.Wait()
}

func test2() {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream <- i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for ret := range resultStream {
		fmt.Printf("Recieve: %d\n", ret)
	}
	fmt.Println("Done receiveing")
}
