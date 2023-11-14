package main

import (
	"fmt"
	"sync"
)

/*
Publisher routine to to publish decoded data
*/

func Publisher(in <-chan string, wg *sync.WaitGroup) {
	pktCount := 0

	defer wg.Done()

	defer func() {
		fmt.Printf("Exiting publisher, after processing %d packets\n", pktCount)
	}()

	for data := range in {
		fmt.Printf("Publishing data: %s", data) // Simulate data publishing process
		pktCount += 1
	}
}
