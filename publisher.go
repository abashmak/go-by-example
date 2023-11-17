package main

import (
	"fmt"
	"sync"
)

/*
Publisher routine to to publish decoded data
*/

func Publisher(inChan <-chan string, wg *sync.WaitGroup) {
	pktCount := 0

	defer wg.Done()

	// Needs to be a function to utilize up-to-date value of pktCount
	defer func() {
		fmt.Printf("Exiting publisher, after processing %d packets\n", pktCount)
	}()

	// Receive on inChan while it's open
	for data := range inChan {
		fmt.Printf("Publishing data: %s", data) // Simulate data publishing process
		pktCount += 1
	}
}
