// 20195845 김휘경
package main

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

const (
	LOGIN          = "1"
	CHAT           = "2"
	ROOM_MAX_USER  = 8 //1개의 채팅방에 최대 8명
	ROOM_MAX_COUNT = 1

	MAX_CLIENTS = 8 // 클라 최대 8명

	CMD_PREFIX = "/"
	CMD_LIST   = CMD_PREFIX + "list"
	CMD_DM     = CMD_PREFIX + "dm"
	CMD_EXIT   = CMD_PREFIX + "exit"
	CMD_VER    = CMD_PREFIX + "ver"
	CMD_RTT    = CMD_PREFIX + "rtt"
)

type Client struct {
	conn net.Conn
	read chan string
	quit chan int
	name string
	room *Room
}

type Room struct {
	num        int
	clientlist *list.List
}

var roomlist *list.List

func cancel(s chan os.Signal, listner net.Listener) {
	<-s
	//fmt.Println("\ngg~")
	//fmt.Println("\nsignal: ", s)
	clientNum--
	fmt.Printf("nickname left. There are %d users now\n", clientNum)
	os.Exit(0)
	listner.Close()
}

//var clientNum int // count client number
var clientNum int = 0 // command count
func main() {

	//ctrl+c
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	serverPort := "25845"

	ln, _ := net.Listen("tcp", ":"+serverPort)
	//fmt.Printf("The Server is ready to receive on port %s\n", serverPort)

	go cancel(c, ln)

	roomlist = list.New()
	for i := 0; i < ROOM_MAX_COUNT; i++ {
		room := &Room{i + 1, list.New()}
		roomlist.PushBack(*room)
	}

	defer ln.Close()

	for {
		//start := time.Now() // rtt ?
		// waiting connection
		conn, _ := ln.Accept()

		clientNum++
		go handleConnection(conn) // 고루틴

		Message := fmt.Sprintf("%d", clientNum)
		conn.Write([]byte(Message))

	}
}

func handleError(conn net.Conn, err error, errmsg string) {
	if conn != nil {
		conn.Close()
	}
	fmt.Println(err)
	fmt.Println(errmsg)
}

func handleConnection(conn net.Conn) {

	read := make(chan string)
	quit := make(chan int)
	client := &Client{conn, read, quit, "unknown", &Room{-1, list.New()}}

	go handleClient(conn, client)
	//fmt.Printf("remote Addr = %s\n", conn.RemoteAddr().String())
}

func handleClient(conn net.Conn, client *Client) {
	for {
		startTime := time.Now()
		select {
		case msg := <-client.read:
			if strings.HasPrefix(msg, "list") { // list
				//ListClients(client)
				var buffer bytes.Buffer
				for e := client.room.clientlist.Front(); e != nil; e = e.Next() {
					el := e.Value.(Client)
					//cn := el.name

					slice := strings.Split(client.conn.RemoteAddr().String(), ":")
					Message := fmt.Sprintf("<%s, %s, %s> ", el.name, slice[0], slice[1])
					buffer.WriteString(Message)
				}
				fmt.Fprintf(client.conn, "%s\n", buffer.String())

			} else if strings.HasPrefix(msg, "dm") { //DM
				sendToClientToClient(client, msg)
			} else if strings.HasPrefix(msg, "ver") {
				Message := fmt.Sprintf("Server version: 8.0\n")
				conn.Write([]byte(Message))
			} else if strings.HasPrefix(msg, "rtt") {
				elaspedTime := time.Since(startTime)
				// RTT
				Message := fmt.Sprintf("RTT = %.3f ms \n\n", float64(elaspedTime.Microseconds())/1000)
				conn.Write([]byte(Message))

			} else if strings.HasPrefix(msg, "exit") {
				clientNum--
				Message := fmt.Sprintf("%s left. There are %d users now\n", client.name, clientNum)
				conn.Write([]byte(Message))
				fmt.Printf("%s left. There are %d users now\n", client.name, clientNum)
				//fmt.Printf("gg~")
				conn.Close()
				break
			} else {
				sendToAllClients(client.name, msg) // 모든 클리에게
			}

		case <-client.quit:
			fmt.Printf("disconnect client, %s left. There are %d users now\n\n", client.name, clientNum)
			client.conn.Close()
			client.deleteFromList()
			return

		default:
			go recvFromClient(conn, client)
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func recvFromClient(conn net.Conn, client *Client) {
	buffer := make([]byte, 1024)
	client.conn.Read(buffer)
	recvmsg := string(buffer)

	strmsgs := strings.Split(recvmsg, "|")

	switch strmsgs[0] {
	case LOGIN:
		client.name = strings.TrimSpace(strmsgs[1])

		// 접속시
		fmt.Printf("%s joined from %s. There are users connected\n", client.name, conn.RemoteAddr().String())
		room := allocateEmptyRoom()
		if room.num < 1 {
			handleError(client.conn, nil, "max user limit!\n")
		}
		client.room = room

		if clientNum > 8 {
			clientNum--
			Message := fmt.Sprintf("chatting room full. cannot connect\n")
			conn.Write([]byte(Message))
			//handleError(client.conn, nil, "max user limit!\n")
		}

		if !client.dupUserCheck() { // 닉네임 이미 있음
			clientNum--
			Message := fmt.Sprintf("\nthat nickname is already used by another user. cannot connect.\n")
			conn.Write([]byte(Message))
			client.quit <- 0
			return
		}
		if strmsgs[0][0] == 'q' {
			clientNum--
			fmt.Printf("%s left. There are %d users now\n", client.name, clientNum)

		}
		//fmt.Printf("\nhello = %s, your room number is = %d\n", client.name, client.room.num)
		room.clientlist.PushBack(*client)

	case CHAT:
		fmt.Printf("\nrecv message = %s\n", strmsgs[1]) // 채팅 메시지 출력
		client.read <- strmsgs[1]
	}
}

func sendToClient(client *Client, sender string, msg string) {
	var buffer bytes.Buffer
	buffer.WriteString(sender) //name
	buffer.WriteString("> ")
	buffer.WriteString(msg)

	//fmt.Printf("client = %s ==> %s", client.name, buffer.String())

	fmt.Fprintf(client.conn, "%s", buffer.String())
}

func sendToClientDM(client *Client, sender string, msg string) {
	var buffer bytes.Buffer
	buffer.WriteString("from ")
	buffer.WriteString(sender) //name
	buffer.WriteString("> ")
	buffer.WriteString(msg)

	//fmt.Printf("client = %s ==> %s", client.name, buffer.String())

	fmt.Fprintf(client.conn, "%s", buffer.String())
}

func sendToClientToClient(client *Client, msg string) { // dm
	strmsgs := strings.Split(msg, " ")

	target := findClientByName(strmsgs[1])
	if target.conn == nil {
		fmt.Println("Can't find target User")
		return
	}
	sendToClientDM(target, client.name, strmsgs[2])
}

func sendToAllClients(sender string, msg string) {
	for re := roomlist.Front(); re != nil; re = re.Next() {
		r := re.Value.(Room)
		for e := r.clientlist.Front(); e != nil; e = e.Next() {
			c := e.Value.(Client)
			sendToClient(&c, sender, msg)
		}
	}
}

func (client *Client) deleteFromList() {
	for re := roomlist.Front(); re != nil; re = re.Next() {
		r := re.Value.(Room)
		for e := r.clientlist.Front(); e != nil; e = e.Next() {
			c := e.Value.(Client)
			if client.conn == c.conn {
				r.clientlist.Remove(e)
			}
		}
	}
}

func (client *Client) dupUserCheck() bool {
	for re := roomlist.Front(); re != nil; re = re.Next() {
		r := re.Value.(Room)
		for e := r.clientlist.Front(); e != nil; e = e.Next() {
			c := e.Value.(Client)
			if strings.Compare(client.name, c.name) == 0 {
				return false
			}
		}
	}

	return true
}

func allocateEmptyRoom() *Room {
	for e := roomlist.Front(); e != nil; e = e.Next() {
		r := e.Value.(Room)

		//fmt.Printf("clientlist len = %d", r.clientlist.Len())
		if r.clientlist.Len() < ROOM_MAX_USER {
			return &r
		}
	}

	// full room
	return &Room{-1, list.New()}
}

func findClientByName(name string) *Client { // dm
	for re := roomlist.Front(); re != nil; re = re.Next() {
		r := re.Value.(Room)
		for e := r.clientlist.Front(); e != nil; e = e.Next() {
			c := e.Value.(Client)
			cname := c.name
			//fmt.Printf("cname : %s name: %s\n ", cname, name)
			if strings.Compare(cname, name) == 0 {
				return &c
			}
		}
	}
	return &Client{nil, nil, nil, "unknown", nil}
}

// Sends to the client the list of chat rooms currently open.
func ListClients(client *Client) {
	for e := client.room.clientlist.Front(); e != nil; e = e.Next() {
		fmt.Printf("<%s, %s>\n", client.name, client.conn.RemoteAddr().String())
	}
	log.Println("client listed chat rooms\n")
}
