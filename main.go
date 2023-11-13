package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
)

const (
	// Can be hard-coded, static config, or even dynamic config
	NumWorkers = 2
)

/*
Load-balancing TCP ingestion server
To send sample TCP packets on the command line:
for ((x=0; x<10; x++)); do echo -e "$x" | netcat localhost 8080 -c && sleep 1; done
*/
func main() {
	// Variable declaration without assignment, explicit type
	// Semaphore counter, thread safe, mutex, atomic integer
	var wg sync.WaitGroup

	// variable declaration with assignment, implicit type
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		// Throw runtime error and exit the program
		panic(fmt.Errorf("error establishing TCP listner: %v", err))
	} else {
		fmt.Printf("Successfully established TCP listner: %v\n", listener.Addr())
		// finally block
		defer listener.Close()
	}

	// buffered channel, like a queue, (can be unbuffered - mmediately blocking)
	ingestChannel := make(chan []byte, 100)
	defer close(ingestChannel)

	publishChannel := make(chan string, 100)
	defer close(publishChannel)

	fmt.Printf("Setting up %d ingestion workers\n", NumWorkers)
	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)                                        // Increment semaphore
		go Worker(ingestChannel, publishChannel, i, &wg) // Spin off new Worker thread
	}

	fmt.Println("Setting up 1 publisher")
	wg.Add(1)                         // Increment semaphore
	go Publisher(publishChannel, &wg) // Spin off new Publisher thread

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, os.Interrupt)
	go func() {
		for range sigChannel {
			fmt.Printf("\nReceived interrupt signal, closing channels and exiting...\n")
			close(ingestChannel)
			close(publishChannel)
			wg.Wait() // Block and wait for the semaphore count to reach 0
			os.Exit(0)
		}
	}()

	fmt.Println("Enterring packet listener loop...")
	for {
		// Needs to be new buffer for each iteration, because passed into worker by reference and could be clobbered
		buffer := make([]byte, 1024)

		conn, err := listener.Accept()
		if err != nil {
			panic(fmt.Errorf("unexpected TCP listner error: %v", err))
		}

		// Declared vars exist only within if/else scope
		if n, err := conn.Read(buffer); err != nil {
			panic(fmt.Errorf("unexpected TCP connection error: %v", err))
		} else {
			ingestChannel <- buffer[:n] // Place data on channel, delivered to single Worker thread in round-robin fashion
			conn.Close()
		}
	}
}
