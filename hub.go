package main

import (
	//"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type hub struct {
	registrations   chan *client
	deregistrations chan *client
	commands        chan command
	clients         map[string]*client
	channels        map[string]*channel
	files           map[string]*file
	rootDir         string
}

// newHub function to create a new center which handler all the
// clients, channels and files in the server.
func newHub() *hub {
	return &hub{
		registrations:   make(chan *client),
		deregistrations: make(chan *client),
		commands:        make(chan command),
		clients:         make(map[string]*client),
		channels:        make(map[string]*channel),
		files:           make(map[string]*file),
		rootDir:         "files",
	}
}

// HUB run - This function turn on a hub who is in charge to manage all the
// posible actions in the FTP server, if you want to join, leave,list,
// register, deregister to a channel or send, download, list a file that
// pettion is aproval do it for the hub, the hub is the brain of the server
func (h *hub) run() {
	for {
		select {
		case client := <-h.registrations:
			h.register(client)
		case client := <-h.deregistrations:
			h.deregister(client)
		case cmd := <-h.commands:
			switch cmd.id {
			case JOIN:
				h.joinChannel(cmd.sender, cmd.receiver)
			case LEAVE:
				h.leaveChannel(cmd.sender, cmd.receiver)
			// download file
			case D_FILE:
				h.downFile(cmd.sender, cmd.fileName)
			//	send file to channel or person.
			case S_FILE:
				h.sendFile(cmd.sender, cmd.receiver, cmd.fileName, cmd.file)
			// list all files
			case L_FILES:
				h.listFiles(cmd.sender)
			// list of users
			case USRS:
				h.listUsers(cmd.sender)
			// list of channels
			case CHNS:
				h.listChannels(cmd.sender)
			// list of files in a specific channel. (if user is in the channel)
			case CHNS_FILE:
				h.listChannelFiles(cmd.sender, cmd.channel)
			default:
				log.Fatal("Command not found")
			}
		}
	}
}

// The Hub use this function to register a new client in the server.
//	(type *client) c - pointer to a new client connection.
//
func (h *hub) register(c *client) {
	if _, exists := h.clients[c.username]; exists {
		c.username = ""
		c.conn.Write([]byte("ERR username taken\n"))
		return
	}
	h.clients[c.username] = c
	c.conn.Write([]byte("OK, successfully register\n"))
}

// The hub use this function to desregister a client
// from a the list of client and channels
// 	(*client) c - Pointer to a client with a living connection
func (h *hub) deregister(c *client) {
	if _, exists := h.clients[c.username]; exists {
		fmt.Printf("User %v was deregister.\n", c.username)
		delete(h.clients, c.username)
		for _, channel := range h.channels {
			delete(channel.participants, c)
		}
	}
}

// The hub use this function to join a user with a specific channel
// if the channel exits otherwiise it create a the channel.
//	(string) user - name of the user that want to join a channel.
// 	(string) chann - name of the channel to join or create.
func (h *hub) joinChannel(user string, chann string) {
	client, ok := h.clients[user]
	if ok {
		if channel, ok := h.channels[chann]; ok {
			channel.participants[client] = true
			client.conn.Write([]byte("User was registered in existent channel.\n"))
			return
		} else {
			ch := newChannel(chann)
			ch.participants[client] = true
			h.channels[chann] = ch
		}
		client.conn.Write([]byte("User was registered in new channel\n"))
	} else {
		fmt.Printf("User not found joinchannel(HUB)")
		//client.conn.Write([]byte("user not found"))
	}
}

// The hub use this function to leave of a specific channel if it exists.
// 	(string) user - user name  that wants exit to the channel.
// 	(string) chann - channel name to exit.
func (h *hub) leaveChannel(user string, chann string) {
	if client, ok := h.clients[user]; ok {
		if channel, ok := h.channels[chann]; ok {
			delete(channel.participants, client)
			resp := "You leave of " + channel.name + " channel"
			client.conn.Write([]byte(resp))
			return
		}
		client.conn.Write([]byte("Channel not found!\n"))
	}
}

// The hub use addFile when user wants to add a new file in a specific channel.
// 	(string) user - user name who sends the new file.
//	(string) chann - channel name which the file will be send.
//	(string) fileName - file name of the send file.
//	([]byte) file - content of the file to add a channel.
func (h *hub) sendFile(user string, dest string, fileName string, file []byte) {
	err := "You are not register in the channel! you need to be register in a channel to send files\n"
	path := filepath.Join(h.rootDir, fileName)
	if sender, ok := h.clients[user]; ok {
		if len(file) == 0 {
			sender.conn.Write([]byte("File not found!\n"))
			return
		}
		switch dest[0] {
		case '#':
			if channel, ok := h.channels[dest]; ok {
				if ok := channel.participants[sender]; !ok {
					sender.conn.Write([]byte(err))
					fmt.Printf("Canceled, %s isn't register in channel\n", sender.username)
					return
				}
				if _, ok := channel.files[fileName]; ok {
					sender.conn.Write([]byte("File with this name already exists\n"))
					return
				}
				new_file := newFile(fileName, path)
				channel.files[new_file.name] = new_file
				h.files[new_file.name] = new_file
				createFile(fileName, file)
				channel.broadchast(sender.username, new_file.name)
				sender.conn.Write([]byte("file was received and created"))
				fmt.Printf("File was received and shared in %v", channel.name)
			} else {
				sender.conn.Write([]byte("Channel not found\n"))
				return
			}
		case '@':
			user, ok := h.clients[dest]
			if !ok {
				sender.conn.Write([]byte("Err no such user!\n"))
				return
			}
			new_file := newFile(fileName, path)
			h.files[new_file.name] = new_file
			createFile(fileName, file)
			sendtoUser(fileName, file, sender.username, user)
		default:
			sender.conn.Write([]byte("ERR with S_FILE command!\n"))
		}
	}
}

// The hub use downFile to allow the client download a file from server.
//	(string) user - Name of the user that wants to download a file.
//	(string) fileName - name of a given file to download.
func (h *hub) downFile(user string, fileName string) {
	if client, ok := h.clients[user]; ok {
		path := filepath.Join(h.rootDir, fileName)
		file, err := os.Open(path)
		if err != nil {
			client.conn.Write([]byte("Error: " + err.Error()))
			return
		}
		defer file.Close()
		var buf = make([]byte, 10000)
		for {
			n, err := file.Read(buf)
			if err != nil {
				if err == io.EOF {
					fmt.Printf("sending file completed \n")
					//client.conn.Write([]byte("File will be send!\n"))
					client.conn.Write(buf[:n])
					//time.Sleep(5 * time.Second)
					//client.conn.Write([]byte("Server: File was send successful!"))
					return
				} else {
					fmt.Printf("%v", err)
					return
				}
			}
			client.conn.Write(buf[:n])
		}
	}
}

// The hub use listChannels for send a list of the current existent channels
// (string) user - name of the user/client.
func (h *hub) listChannels(user string) {
	if client, ok := h.clients[user]; ok {
		var channels []string

		if len(h.channels) == 0 {
			client.conn.Write([]byte("ERR no channels found"))
			return
		}

		for chann := range h.channels {
			channels = append(channels, chann)
		}
		resp := strings.Join(channels, "\n")
		client.conn.Write([]byte(resp + "\n"))
	}
}

// The hub use listFiles to return a list of all files in the server.
//	(string) user - name of the user who do the request.
func (h *hub) listFiles(user string) {
	if client, ok := h.clients[user]; ok {
		var files []string
		if len(h.files) == 0 {
			client.conn.Write([]byte("ERR no files found"))
			return
		}
		for fil := range h.files {
			files = append(files, "-"+fil)
		}
		resp := strings.Join(files, " ")
		client.conn.Write([]byte(resp + "\n"))
	}
}

// The hub uses listUsers to return the list of all users registered in the server.
//	(string) user - name of user who do the request.
func (h *hub) listUsers(user string) {
	if client, ok := h.clients[user]; ok {
		var users []string
		if len(h.clients) == 1 {
			client.conn.Write([]byte("ERR no users found"))
			return
		}
		for userName := range h.clients {
			if userName == client.username {
				continue
			}
			users = append(users, userName)
		}
		resp := strings.Join(users, ", ")
		client.conn.Write([]byte(resp + "\n"))
	}
}

// The hub use listChannelfiles to send a list of files register in the channel
//	(string) user - user name who do the request
// 	(string) channel - name of the chhannel that contain the files.
func (h *hub) listChannelFiles(user string, channel string) {
	if client, ok := h.clients[user]; ok {
		if chann, ok := h.channels[channel]; ok {
			var listFiles []string

			if len(chann.files) == 0 {
				client.conn.Write([]byte("Error: no files found in channel\n"))
				return
			}
			//listFiles = append(listFiles, chann.name)
			for fil := range chann.files {
				listFiles = append(listFiles, "\n"+"-"+fil)
			}
			resp := strings.Join(listFiles, ",")
			client.conn.Write([]byte(resp + "\n"))
		}
	}

}

//	sendtoUser - function use when user send a file to another file
//	 (string) fileName - file name which will be send.
//	 ([]byte) file - File content.
//	 (string) sender - User that send the file.
//	 (*client) user - User which the file will be send.
func sendtoUser(fileName string, file []byte, sender string, user *client) {
	user.conn.Write([]byte("file will be send from user!\n"))
	time.Sleep(10 * time.Second)
	user.conn.Write([]byte(fileName + "\n"))
	time.Sleep(10 * time.Second)
	user.conn.Write(file)
	user.conn.Write([]byte("send by: " + sender + "\n"))
}
