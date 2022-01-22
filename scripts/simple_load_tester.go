package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"net/http"
	"sync"
	"time"
)

// This is a utility to run simple load tests on the app. It was made to fix the
// "concurrent map write" bug from FLAP-141
func main() {
	times := flag.Int("times", 100, "the number of times to make the request")
	flag.Parse()
	for r := range blastOff(*times, request) {
		if r.Err != nil {
			fmt.Println(r.Err)
		} else if bod, err := io.ReadAll(r.Response.Body); err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("STATUS: %v, BODY: %v", r.Response.Status, string(bod))
		}
	}
}

type Result struct {
	Response *http.Response
	Err      error
}

// Tweak this function to change the request made
func request() *http.Request {
	req, _ := http.NewRequest(
		http.MethodPost,
		"http://localhost:8080/simulate?tape=aaba",
		io.NopCloser(bytes.NewBufferString(`
		{
			"Start": "q0",
			"Type": "DFA",
			"Alphabet": "ab",
			"States": [
			{ "Id": "q0", "Ending": false },
			{ "Id": "q1", "Ending": true }
			],
			"Transitions": [
			{ "Start": "q0", "End": "q1", "Symbol": "a" },
			{ "Start": "q0", "End": "q0", "Symbol": "b" },
			{ "Start": "q1", "End": "q1", "Symbol": "b" },
			{ "Start": "q1", "End": "q0", "Symbol": "a" }
			]
		}
		`)))
	return req
}

// Executes an http request on `times` goroutines. Each go routine is equipped
// with a circuit breaker such that, if the request fails, the goroutine backs
// off for an increasingly longer period (starting from 50 milliseconds) before
// retrying the request. If the request fails 20 times in a row, the goroutine
// sends an error on the response channel and exits
func blastOff(times int, request func() *http.Request) <-chan Result {
	out := make(chan Result, times)

	var wait sync.WaitGroup
	for i := 0; i < times; i++ {
		wait.Add(1)
		go func() {
			circuitBreaker := 1
			for {
				resp, err := http.DefaultClient.Do(request())
				exit := func() {
					out <- Result{resp, err}
					wait.Done()
				}
				if err != nil {
					if circuitBreaker > 20 {
						exit()
						break
					}
					time.Sleep(time.Duration(math.Min(5, float64(circuitBreaker))) * 50 * time.Millisecond)
					circuitBreaker++
					continue
				}
				exit()
				break
			}
		}()
	}

	go func() {
		wait.Wait()
		close(out)
	}()

	return out
}
