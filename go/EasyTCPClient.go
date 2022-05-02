/* 20195845 hwikyungkim
 * TCP Client
 */
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

func cancel2(c chan os.Signal, conn net.Conn) {
	<-c
	conn.Write([]byte("5"))
	fmt.Println("\nBye bye~")
	conn.Close()
	os.Exit(0)
}

func main() {

	serverName := "nsl2.cau.ac.kr"
	serverPort := "25845"

	// connect to server
	conn, err := net.Dial("tcp", serverName+":"+serverPort)

	if err != nil {
		fmt.Println(err)
		return
	}
	localAddr := conn.LocalAddr().(*net.TCPAddr)
	fmt.Printf("Client is running on port %d\n", localAddr.Port)

	// ctrl +c
	cancleNoti := make(chan os.Signal)
	signal.Notify(cancleNoti, syscall.SIGINT, syscall.SIGTERM)
	go cancel2(cancleNoti, conn)

	defer conn.Close() // main 함수가 끝나기 직전에 TCP 연결을 닫음

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
			conn.Write([]byte("1" + inputString))
			endTime := time.Now().UnixMilli()
			rtt := endTime - startTime
			conn.Read(buffer)
			fmt.Printf("\nReply from server: %s", string(buffer))
			// RTT
			fmt.Printf("RTT : %d ms\n\n", rtt)
		} else if inputOption == "2\n" {
			startTime := time.Now().UnixMilli()
			// to server
			conn.Write([]byte("2"))
			// from server
			conn.Read(buffer)
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
			conn.Write([]byte("3"))
			conn.Read(buffer)
			endTime := time.Now().UnixMilli()
			rtt := endTime - startTime
			fmt.Printf("\nReply from server: client request number = %s\n", string(buffer))
			// RTT
			fmt.Printf("RTT : %d ms\n\n", rtt)
		} else if inputOption == "4\n" {
			startTime := time.Now().UnixMilli()
			conn.Write([]byte("4"))
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			endTime := time.Now().UnixMilli()
			rtt := endTime - startTime
			fmt.Printf("\nReply from server: %s", string(buffer))
			// RTT
			fmt.Printf("RTT : %d ms\n\n", rtt)

		} else if inputOption == "5\n" {
			conn.Write([]byte("5"))
			fmt.Println("Bye bye~")
			conn.Close()
			os.Exit(0)
		}

	}
}
