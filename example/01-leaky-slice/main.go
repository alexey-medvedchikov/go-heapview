package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/alexey-medvedchikov/go-heapview"
)

func main() {
	flag.Parse()

	size := 1000
	sizeArg := flag.Arg(0)

	if sizeArg != "" {
		var err error
		if size, err = strconv.Atoi(sizeArg); err != nil {
			log.Fatalf("could not parse %s as int", sizeArg)
		}
	}

	// A convenience function, feel free to copy-n-paste the code to costomize to your need.
	//   The only thing needed to inspect the heap is heap dump file, there's no other magic.
	heapview.SetupHandler(func() string { return "heapdump.dat" }, syscall.SIGUSR1)

	// As client and its contents are on the stack, we will report something like
	//   "an object of size 24 (slice size, see cache field) in function main.main
	//      holds len(Client.cache) * len(Message.contents) bytes of memory"
	var client Client

	for i := 0; i < size; i++ {
		// Memory allocation is done in MakeMessage, but it is not a cause of memory leak
		client.SaveMessage(MakeMessage())
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	// Use cache slice to make sure it is not GC-collected when we write heap dump
	fmt.Printf("cache length = %d", len(client.cache))
}

type Message struct {
	contents []byte
}

func MakeMessage() *Message {
	contents := make([]byte, 1024)
	for i := 0; i < len(contents); i++ {
		contents[i] = byte(i)
	}

	return &Message{
		contents: contents,
	}
}

type Client struct {
	cache []*Message
}

func (c *Client) SaveMessage(msg *Message) {
	// This is ultimately a line where leak happens. cache slice holds pointers to Messages while they
	//   hold pointers to contents slices.
	c.cache = append(c.cache, msg)
}
