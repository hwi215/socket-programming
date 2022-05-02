/**
* 20195845 hwikyungkim
* UDP Server
 */
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

func cancel(c chan os.Signal) {
	<-c
	fmt.Println("\nBye bye")
	os.Exit(0)
}

func main() {

	serverPort := "25845"
	pconn, _ := net.ListenPacket("udp", ":"+serverPort)
	fmt.Printf("The Server is ready to receive on port %s\n", serverPort)
	var command int = 0 // command count

	for {
		start := time.Now() // option4

		//ctrl+c
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		go cancel(c)

		for {
			buffer := make([]byte, 1024)
			count, addr, _ := pconn.ReadFrom(buffer)
			fmt.Printf("\nConnection request from %s\n", addr.String())
			command++

			n := string(buffer[0])
			if n == "1" {
				// option1: upper
				fmt.Printf("Command 1\n\n")
				pconn.WriteTo(bytes.ToUpper(buffer[1:count]), addr)
			} else if n == "2" {
				// option2
				fmt.Printf("Command 2\n\n")
				address := addr.String()
				pconn.WriteTo([]byte(address), addr)
			} else if n == "3" {
				// option3
				fmt.Printf("Command 3\n\n")
				com := strconv.Itoa(command)
				pconn.WriteTo([]byte(com), addr)

			} else if n == "4" {
				// option4
				fmt.Printf("Command 4\n\n")
				runtime := runtimeF(time.Since(start))
				Message := fmt.Sprintf("running time = %s\n", runtime)
				pconn.WriteTo([]byte(Message), addr)

			} else if n == "5" {
				// option5
				fmt.Printf("Command 5\n\n")
				//pconn.Close()
				fmt.Printf("Bye bye client : %s\n", addr.String())
				break

			}

		}
	}
}
