package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

const manual = `
Notices: @ in Username and # in channelname are mandatory.
  -> in command S_FILE is mandatory.

  COMMANDS:
    REG <@USERNAME>
    JOIN <#CHANNELNAME>
    LEAVE <#CHANNELNAME>
    D_FILE <nameOfFile>
    S_FILE <nameOfFile> -> <#CHANNELNAME>
    L_FILES
    USRS
    CHNS
    CHNS_FILE <#CHANNELNAME>
    QUIT
`

func main() {

	connection, err := net.Dial("tcp", "localhost:8031")
	//var respo = make([]byte, 10000)
	ch := make(chan []byte)
	eCh := make(chan error)
	fileName := make(chan string)

	if err != nil {
		fmt.Printf("There was an error making a connection\n")
		return
	}

	go read(ch, eCh, fileName, connection)
	go response(ch)
	fmt.Print(manual)

	for {
		//go readFile(ch, connection)
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter your command: \n\n")
		input, _ := reader.ReadString('\n')
		cmd, info, _ := strings.Cut(input, " ")

		switch cmd {
		case "QUIT":
			fmt.Print("You close the session!\n")
			connection.Write([]byte(io.EOF.Error()))
			return
		case "REG":
			msg := []byte(info)
			if msg[0] != '@' {
				fmt.Printf("Error: username must begin with @\n")
				continue
			}
			connection.Write([]byte(input))
		case "JOIN":
			if info == "" {
				fmt.Printf("Error: you have to entry a channel!\n")
				continue
			}
			msg := []byte(info)
			if msg[0] != '#' {
				fmt.Printf("Error: channel name must begin with '#'\n")
				continue
			}
			connection.Write([]byte(input))
		case "LEAVE":
			if info == "" {
				fmt.Printf("Error: you have to entry a channel!\n")
				continue
			}
			msg := []byte(info)
			if msg[0] != '#' {
				fmt.Printf("Error: channel name must begin with '#'\n")
			}
			connection.Write([]byte(input))
		case "S_FILE":
			fileName, channel, _ := strings.Cut(info, " -> ")
			ch := []byte(channel)
			if ch[0] == '#' || ch[0] == '@' {
				connection.Write([]byte(input))
				time.Sleep(5 * time.Second)
				sendFile(fileName, connection)
				continue
			}
			fmt.Printf("Error: channel name must begin with '#'\nand user name must begin with '@'")
			continue
		case "L_FILES":
			connection.Write([]byte(input))
		case "D_FILE":
			if info == "" {
				fmt.Printf("Error: Name of file is mandatory!\n")
				continue
			}
			connection.Write([]byte(input))
			fileName <- info
			continue
		case "USRS":
			connection.Write([]byte(input))
		case "CHNS":
			connection.Write([]byte(input))
		case "CHNS_FILE":
			if info == "" {
				fmt.Printf("Error: Name of channel is mandatory!\n")
			}
			msg := []byte(info)
			if msg[0] != '#' {
				fmt.Printf("Error: channel name must begin with '#'\n")
				continue
			}
			connection.Write([]byte(input))
		case "\n":
			continue
		default:
			fmt.Printf("Invalid command\n\n")
		}
	}
}

// read - function that obtain all responses or files from the server.
//	(chan []byte) ch - all information coming from the server will be
//	past throught this chan.
//	(chan error) eCh- in case of error it will be pass throutg this chan.
//	(chan string) fileName - used in case of download a file from the server.
//	(net.Conn) connection -  connection with the server.
func read(ch chan []byte, eCh chan error, fileName chan string,
	connection net.Conn) {
	for {
		data, err := bufio.NewReader(connection).ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				fmt.Print("The server is down\n")
				return
			}
			eCh <- err
			return
		}
		if data != nil {
			msgCh := "file will be send from a channel!\n"
			msgUsr := "file will be send from user \n"

			if msgCh == string(data) || msgUsr == string(data) {
				time.Sleep(9 * time.Second)
				fileName, err := bufio.NewReader(connection).ReadBytes('\n')
				if err != nil {
					fmt.Printf("Error: %v", err)
					continue
				}
				fileContent, err := bufio.NewReader(connection).ReadBytes('\n')
				if err != nil {
					fmt.Printf("Error: %v", err)
					continue
				}
				ch <- data
				time.Sleep(3 * time.Second)
				fmt.Print("name of file is: " + string(fileName) + "\n")
				createFile(fileContent, string(fileName))
				continue
			}
		}
		select {
		case fileName := <-fileName:
			createFile(data, fileName)
			continue
		default:
			ch <- data
		}
		//ch <- data
	}
}

// response - function that prints all information passed throught
// the ch channel used in read function.
//	(chan []byte) ch - channel with the information coming from server.
func response(ch chan []byte) {
	for {
		fmt.Print(string(<-ch))
	}
}
