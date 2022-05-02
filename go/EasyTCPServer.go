/**
* 20195845 hwikyungkim
* TCP Server
 **/

package main

import (
	"bytes"
	"fmt"
	"math"
	"net"
	"os"
	"os/signal"
	"strconv"
	_ "strings"
	"syscall"
	"time"
)

func runtimeF(duration time.Duration) string { // time format

	HH := int64(math.Mod(duration.Hours(), 100))
	MM := int64(math.Mod(duration.Minutes(), 60))
	SS := int64(math.Mod(duration.Seconds(), 60))
	return fmt.Sprintf("%02d:%02d:%02d", HH, MM, SS)
}

func cancel(c chan os.Signal, conn net.Conn) {
	s := <-c
	fmt.Println("\nBye bye")
	fmt.Println("\nsignal: ", s)
	os.Exit(0)
	conn.Close()
}

func main() {

	serverPort := "25845"
	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("The Server is ready to receive on port %s\n", serverPort)
	var command int = 0 // command count

	//ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	for {
		conn, _ := listener.Accept()
		start := time.Now() // option4
		//defer listener.Close()
		go cancel(c, conn)

		for {
			buffer := make([]byte, 1024)
			count, _ := conn.Read(buffer)
			fmt.Printf("\nConnection request from %s\n", conn.RemoteAddr().String())
			command++

			n := string(buffer[0])
			if n == "1" {
				// option1: upper
				fmt.Printf("Command 1\n\n")
				conn.Write(bytes.ToUpper(buffer[1:count]))
			} else if n == "2" {
				// option2
				fmt.Printf("Command 2\n\n")
				address := conn.RemoteAddr().String()
				conn.Write([]byte(address))
			} else if n == "3" {
				// option3
				fmt.Printf("Command 3\n\n")
				com := strconv.Itoa(command)
				conn.Write([]byte(com))

			} else if n == "4" {
				// option4
				fmt.Printf("Command 4\n\n")
				runtime := runtimeF(time.Since(start))
				Message := fmt.Sprintf("running time = %s\n", runtime)
				conn.Write([]byte(Message))

			} else if n == "5" {
				// option5
				fmt.Printf("Command 5\n\n")
				conn.Close()
				fmt.Printf("Bye bye client : %s\n", conn.RemoteAddr().String())
				break

			}

		}
	}
}
