// 20195845 김휘경

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func cancel2(c chan os.Signal, conn net.Conn) {
	<-c
	fmt.Println("\ngg~")
	fmt.Fprintf(conn, "%s|p", LOGIN)
	fmt.Printf("nickname left. There are <N> users now")
	conn.Close()
	os.Exit(0)
}

const (
	LOGIN = "1"
	CHAT  = "2"
)

var clientNum int // count client number

func main() {

	serverName := "nsl2.cau.ac.kr"
	serverPort := "25845"

	conn, _ := net.Dial("tcp", serverName+":"+serverPort)

	localAddr := conn.LocalAddr().(*net.TCPAddr)
	//fmt.Printf("Client is running on port %d\n", localAddr.Port)

	// ctrl +c
	cancleNoti := make(chan os.Signal)
	signal.Notify(cancleNoti, syscall.SIGINT, syscall.SIGTERM)

	go cancel2(cancleNoti, conn)

	defer conn.Close() // main 함수가 끝나기 직전에 TCP 연결을 닫음
	//
	msgch := make(chan string)

	name := os.Args[1]

	//fmt.Printf("\n")

	// login
	fmt.Fprintf(conn, "%s|%s", LOGIN, name)

	buffer := make([]byte, 1024)
	conn.Read(buffer)
	//fmt.Printf("%s\n", string(buffer))

	// chat
	go handleRecvMsg(conn, msgch)
	fmt.Printf("[Welcome %s to CAU network class chat room at %s:%d.]\n", name, localAddr.IP, localAddr.Port)
	fmt.Printf("[There are %s users connected]\n\n", string(buffer))

	handleSendMsg(conn)

}

func handleError(conn net.Conn, errmsg string) {
	if conn != nil {
		conn.Close()
	}
	fmt.Println(errmsg)
}

func handleSendMsg(conn net.Conn) {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("\n") // msg 보내기
		text, _ := reader.ReadString('\n')

		text = strings.TrimLeft(text, "\\")

		if text == "rtt" {
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("%s\n", string(buffer))
		} else if text == "ver" {
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("%s\n", string(buffer))
		} else if text == "list" {
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("%s\n", string(buffer))

		} else if text == "exit" {
			buffer := make([]byte, 1024)
			conn.Read(buffer)
			fmt.Printf("%s\n", string(buffer))
			fmt.Printf("gg~")

		} else {
			fmt.Fprintf(conn, "%s|%s", CHAT, text)
		}

	}

}

func handleRecvMsg(conn net.Conn, msgch chan string) {
	for {
		select {
		case msg := <-msgch:
			fmt.Printf("%s\n", msg)
		default:
			go recvFromServer(conn, msgch)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func recvFromServer(conn net.Conn, msgch chan string) {
	msg, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		handleError(conn, "error")
		os.Exit(2)
		return
	}
	msgch <- msg
}
