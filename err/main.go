package main

import "fmt"
import "net/http"

func main() {
	fmt.Println("vim-go")
	test2()
}

func test1() {
	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan *http.Response {
		response := make(chan *http.Response)
		go func() {
			defer close(response)
			for _, url := range urls {
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println(err)
					continue
				}
				select {
				case <-done:
					return
				case response <- resp:
				}
			}
		}()
		return response
	}

	done := make(chan interface{})
	defer close(done)

	urls := []string{"http://www.google.com", "http://badhost"}
	for responce := range checkStatus(done, urls...) {
		fmt.Printf("Response: %v\n", responce.Status)
	}
}

func test2() {
	type Result struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(
		done <-chan interface{},
		urls ...string,
	) <-chan Result {
		response := make(chan Result)
		go func() {
			defer close(response)
			for _, url := range urls {
				resp, err := http.Get(url)
				result := Result{Error: err, Response: resp}
				select {
				case <-done:
					return
				case response <- result:
				}
			}
		}()
		return response
	}

	done := make(chan interface{})
	defer close(done)

	errCount := 0
	urls := []string{"a", "http://www.google.com", "http://badhost", "b", "c", "d"}
	for result := range checkStatus(done, urls...) {
		if result.Error != nil {
			fmt.Printf("error : %v\n", result.Error)
			errCount++
			if errCount >= 3 {
				fmt.Println("too many error breaking!")
				break
			}
			continue
		}
		fmt.Printf("Response: %v\n", result.Response.Status)
	}

}
