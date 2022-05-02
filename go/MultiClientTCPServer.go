/**
* 20195845 hwikyungkim
* MultiClient TCPServer
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

// client ID
var clientId int  // client clientId
var clientNum int // 누적 클라이언트 수

func runtimeF(duration time.Duration) string { // time format

	HH := int64(math.Mod(duration.Hours(), 100))
	MM := int64(math.Mod(duration.Minutes(), 60))
	SS := int64(math.Mod(duration.Seconds(), 60))
	return fmt.Sprintf("%02d:%02d:%02d", HH, MM, SS)
}

func cancel(s chan os.Signal, listner net.Listener) {
	<-s
	fmt.Println("\nBye bye~")
	//fmt.Println("\nsignal: ", s)
	os.Exit(0)
	listner.Close()
}

func ConnHandler(conn net.Conn, clientId int) {
	start := time.Now() // option4
	var command int = 0 // command count

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
			fmt.Printf("Bye bye client %d : %s\n", clientId, conn.RemoteAddr().String())
			clientNum--
			fmt.Printf("\nClient %d disconnected. Number of connected clients = %d\n\n", clientId, clientNum)
			break

		}
	}
}

func timmer() {
	// start
	fmt.Printf("Number of connected clients = %d\n", clientNum)
	for {
		time.Sleep(time.Second * 60) // 1분마다
		fmt.Printf("Number of connected clients = %d\n", clientNum)
	}
}

func main() {

	//ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	serverPort := "25845"

	// clientID 초기화
	clientId = 0
	clientNum = 0

	listener, _ := net.Listen("tcp", ":"+serverPort)
	fmt.Printf("The Server is ready to receive on port %s\n", serverPort)

	/* 1분 마다 클라이언트 수 출력하기기	 */
	go timmer()

	go cancel(c, listener)

	for {
		/* 연결 끊기면 client ID & 서버의 클라이언트 수 반환
		“Client 1 connected. Number of connected clients = 2”
		“Client 2 disconnected. Number of connected clients = 1”
		*/
		conn, _ := listener.Accept()
		clientId++ // 클라이언트 ID값 증가
		clientNum++
		//connectedCount++
		fmt.Printf("Client %d connected. Number of connected clients = %d\n", clientId, clientNum)

		go ConnHandler(conn, clientId) // go rutine

	}
}
