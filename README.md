# go-by-example

A selective tour of Go by example, with emphasis on multi-threaded programming.

Topics covered:
* Variable declaration/assignment
* Range looping
* Panic
* Defer
* Thread signaling (work group, channel)
* Go-routines

## Architecture
This code implements the following simulated data pipeline:

![alt text](data-pipeline.png "Architecture")

## Compiling
`go build -o app`

## Running
1. Start the server: `./app`
2. Make sure [netcat](https://netcat.sourceforge.net/) is installed
3. Send sample TCP packets:
   ```sh
   for ((x=0; x<10; x++)); do echo -e "$x" | netcat localhost 8080 -c && sleep 1; done
   ```
