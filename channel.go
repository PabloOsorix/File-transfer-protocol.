package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type channel struct {
	name         string
	participants map[*client]bool
	files        map[string]*file
}

// newChannel function create a new channel in the server
//	 (string) channelName - channel name to create
// 	 return : pointer to a channel struct (type channel).
func newChannel(channelName string) *channel {
	return &channel{
		name:         channelName,
		participants: make(map[*client]bool),
		files:        make(map[string]*file),
	}

}

// broadcast function allow the server share
// a given file with all of participants in the channel
//	(string) user - name of the user that sends the file.
// 	(string) fileName - name of the file to share in the channel.
func (c *channel) broadchast(user string, fileName string) {
	path := filepath.Join("files", fileName)
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	buf := make([]byte, 10000)

	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("sending file completed \n")
				return
			} else {
				fmt.Print(err.Error())
			}
		}
		for cl := range c.participants {
			if cl.username == user {
				continue
			}
			cl.conn.Write([]byte("file will be send from a channel!\n"))
			time.Sleep(10 * time.Second)
			cl.conn.Write([]byte(fileName + "\n"))
			time.Sleep(10 * time.Second)
			cl.conn.Write(buf[:n])
			cl.conn.Write([]byte("send by: " + user + "\n"))
		}
	}
}
