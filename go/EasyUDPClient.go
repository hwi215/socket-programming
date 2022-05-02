/**
* 20195845 hwikyungkim
* UDP Client
**/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	_ "strconv"
	"strings"
	"syscall"
	"time"
)

func cancel(c chan os.Signal) {
	<-c
	fmt.Println("\nBye bye~")
	os.Exit(0)
}

func main() {

	serverName := "nsl2.cau.ac.kr"
	serverPort := "25845"

	// connect to server
	pconn, _ := net.ListenPacket("udp", ":")

	/*
		if err != nil {
			fmt.Println(err)
			return
		}

	*/
	localAddr := pconn.LocalAddr().(*net.UDPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	// ctrl +c
	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	go cancel(c)

	serverAddr, _ := net.ResolveUDPAddr("udp", serverName+":"+serverPort)

	defer pconn.Close() // main 함수가 끝나기 직전에 TCP 연결을 닫음

	for {
		// option menu
		fmt.Print("<Menu>\n1) convert text to UPPER-case\n2) get my IP address and port number\n3) get server request count\n4) get server running time\n5) exit\n")

		// input option
		kbReader := bufio.NewReader(os.Stdin)
		fmt.Print("Input option: ")
		inputOption, err := kbReader.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		buffer := make([]byte, 1024)

		if inputOption == "1\n" {
			fmt.Print("Input sentence: ")
			startTime := time.Now().UnixMilli()
			inputString, _ := kbReader.ReadString('\n')
			pconn.WriteTo([]byte("1"+inputString), serverAddr)
			endTime := time.Now().UnixMilli()
			rtt := endTime - startTime
			pconn.ReadFrom(buffer)
			fmt.Printf("\nReply from server: %s", string(buffer))
			// RTT
			fmt.Printf("RTT : %d ms\n\n", rtt)
		} else if inputOption == "2\n" {
			startTime := time.Now().UnixMilli()
			// to server
			pconn.WriteTo([]byte("2"), serverAddr)
			// from server
			pconn.ReadFrom(buffer)
			endTime := time.Now().UnixMilli()
			rtt := endTime - startTime
			s := string(buffer)
			s2 := strings.Split(s, ":")
			ipAddress := s2[0]
			portNumber := s2[1]
			fmt.Printf("\nReply from server: client IP = %s, port = %s\n", ipAddress, portNumber)
			// RTT
			fmt.Printf("RTT : %d ms\n\n", rtt)
		} else if inputOption == "3\n" {
			startTime := time.Now().UnixMilli()
			pconn.WriteTo([]byte("3"), serverAddr)
			pconn.ReadFrom(buffer)
			endTime := time.Now().UnixMilli()
			rtt := endTime - startTime
			fmt.Printf("\nReply from server: client request number = %s\n", string(buffer))
			// RTT
			fmt.Printf("RTT : %d ms\n\n", rtt)
		} else if inputOption == "4\n" {
			startTime := time.Now().UnixMilli()
			pconn.WriteTo([]byte("4"), serverAddr)
			buffer := make([]byte, 1024)
			pconn.ReadFrom(buffer)
			endTime := time.Now().UnixMilli()
			rtt := endTime - startTime
			fmt.Printf("\nReply from server: %s", string(buffer))
			// RTT
			fmt.Printf("RTT : %d ms\n\n", rtt)

		} else if inputOption == "5\n" {
			pconn.WriteTo([]byte("5"), serverAddr)
			fmt.Println("Bye bye~")
			pconn.Close()
			os.Exit(0)
		}

	}
}
