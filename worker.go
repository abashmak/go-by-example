package main

import (
	"fmt"
	"sync"
)

func Worker(inChan <-chan []byte, outChan chan<- string, workerId int, wg *sync.WaitGroup) {
	pktCount := 0

	// Decrement semaphore before exiting the Worker scope
	defer wg.Done()

	// Needs to be a function to utilize updated values in workerId and pktCount
	defer func() {
		fmt.Printf("Exiting packet worker #%d, after processing %d packets\n", workerId, pktCount)
	}()

	// This loop will exit when inChan is closed
	for packet := range inChan {
		decoded := string(packet) // Simulate decoding/processing
		fmt.Printf("Packet worker #%d: decoded packet: %s", workerId, decoded)
		outChan <- decoded // Place data on output channel, delivered to Publisher thread
		pktCount += 1
	}
}
