package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
)

var (
	DELIMITER = []byte(`\r\n`)
)

type client struct {
	conn        net.Conn
	outbound    chan<- command
	register    chan<- *client
	desregister chan<- *client
	username    string
}

// newClient - function that that create a new client inside of hub
// throught net.Conn to can have a register of clients.
//	 (net.Conn) conn - connection throught net.Conn
//	 (out chan<- command) out - receive a chan of kind command to edit and
//	 send the commands to the hub (server.)
//	 (reg chan<- *client) reg - receive a chann which it will send a pointer
//	 to the current client that will be register in the hub.
//	 (desrg chan<- *client) desrg - chann to send the pointer of the current
//	 client to unregister this.
//	 return: a pointer of kind *client (type client)
func newClient(conn net.Conn, out chan<- command, reg chan<- *client,
	desrg chan<- *client) *client {

	return &client{
		conn:        conn,
		outbound:    out,
		register:    reg,
		desregister: desrg,
	}
}

// read - function to receive command from client to send it to.
// handle function which send the order to the hub of the server.
//	return: an error in case of fail, nil in successs.
func (c *client) read() error {
	for {
		order, err := bufio.NewReader(c.conn).ReadBytes('\n')
		if err == io.EOF {
			c.desregister <- c
			c.conn.Write([]byte("connection was close!"))
			return nil
		}
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}
		if err := c.handle(order); err != nil {
			c.conn.Write([]byte(err.Error()))
		}
	}
	return nil
}

// handle - function that receive the command from read function
// and check the validity of it, if it's ok it send the correct
// order in the server hub.
//	([]byte) order - command to analyze to make a decision.
// 	return: error in case of fail, nil in success.
func (c *client) handle(order []byte) error {
	cmd := bytes.ToUpper(bytes.TrimSpace(bytes.Split(order, []byte(" "))[0]))
	args := bytes.TrimSpace(bytes.TrimPrefix(order, cmd))

	switch string(cmd) {
	case "REG":
		if err := c.reg(args); err != nil {
			c.err(err)
		}
		return nil
	case "JOIN":
		if err := c.join(args); err != nil {
			c.err(err)
		}
		return nil
	case "LEAVE":
		if err := c.leave(args); err != nil {
			c.err(err)
		}
		return nil
	case "S_FILE":
		if err := c.send(args); err != nil {
			c.err(err)
		}
		return nil
	case "D_FILE":
		if err := c.downFile(args); err != nil {
			c.err(err)
		}
		return nil
	case "L_FILES":
		c.listf()
		return nil
	case "USRS":
		c.users()
		return nil
	case "CHNS":
		c.chns()
		return nil
	case "CHNS_FILE":
		if err := c.listFilesch(args); err != nil {
			c.err(err)
		}
		return nil
	default:
		c.err(fmt.Errorf("unkown command %s", cmd))
		return nil
	}
}

// reg - function to register a new client in the hub.
//	([]byte) args - arguments that should have the correct username.
//	return: error in case of fail, nil in success.
func (c *client) reg(args []byte) error {
	u := bytes.TrimSpace(args)
	if u[0] != '@' {
		return fmt.Errorf("username must begin with '@'")
	}
	if len(u) == 0 {
		return fmt.Errorf("username cannot be blank")
	}
	c.username = string(args)
	c.register <- c
	fmt.Print(c.username + " was register!\n")
	return nil
}

// join - function to join a user to a channel register in
// the server hub.
//	([]byte) args - arguments with the channel name to join
//	return: error in case of fail, nil in success.
func (c *client) join(args []byte) error {
	channelID := bytes.TrimSpace(args)
	if channelID[0] != '#' {
		return fmt.Errorf("ERR channel ID must begin with #")
	}

	c.outbound <- command{
		receiver: string(channelID),
		sender:   c.username,
		id:       JOIN,
	}
	return nil
}

// leave - function to leave channel register in the server hub.
// ([]byte) args - arguments with the channel name  to leave.
//	return: error in case of fail, nil in success.
func (c *client) leave(args []byte) error {
	channelID := bytes.TrimSpace(args)
	if channelID[0] != '#' {
		return fmt.Errorf("ERR channelID must startr with '#'")
	}

	c.outbound <- command{
		receiver: string(channelID),
		sender:   c.username,
		id:       LEAVE,
	}
	return nil
}

// send - function to send a new file to another user or existent
// channel in the server hub.
//	([]byte) args - arguments with information about file name,
//	and channel name or user name to send the file.
//	return: error in case of fail, nil in success.
func (c *client) send(args []byte) error {
	fileName := bytes.Split(args, []byte(" ->"))[0]
	dest := bytes.TrimSpace(bytes.Split(args, []byte("->"))[1])
	if len(fileName) == 0 {
		return fmt.Errorf("file name cannot be blank")
	}
	if dest[0] == '#' || dest[0] == '@' {
		file, err := c.recvFile()
		if err != nil {
			return fmt.Errorf("file cannot be receive")
		}
		c.outbound <- command{
			sender:   c.username,
			receiver: string(dest),
			fileName: string(fileName),
			channel:  string(dest),
			file:     file,
			id:       S_FILE,
		}
		return nil
	}
	return fmt.Errorf("ERR channelID must start with '#'")
}

// downFile - function to download a existent file in the server.
//	([]byte) args - arguments with the file name to download.
//	return: error in case of fail, nil in success.
func (c *client) downFile(args []byte) error {
	fileName := bytes.TrimSpace(args)
	c.outbound <- command{
		sender:   c.username,
		fileName: string(fileName),
		id:       D_FILE,
	}
	return nil
}

// chns - function which allows you to get a list of all
// channels registered in the server hub.
func (c *client) chns() {
	c.outbound <- command{
		sender: c.username,
		id:     CHNS,
	}
}

// users - function  which allows you to get a list of all
// users registered in the server hub.
func (c *client) users() {
	c.outbound <- command{
		sender: c.username,
		id:     USRS,
	}
}

// listf - function which allows you to get a list of all
// existing files in the server hub.
func (c *client) listf() {
	c.outbound <- command{
		sender: c.username,
		id:     L_FILES,
	}
}

// listFilesch - function which allows you to get a list of all
// existing files linked to a specific channel.
//	([]byte) args - arguments that contain the name of channel which
//	we want get the list of files.
func (c *client) listFilesch(args []byte) error {
	channelID := bytes.TrimSpace(args)
	c.conn.Write(channelID)
	if channelID[0] != '#' {
		return fmt.Errorf("ERR channelID must start with '#'")
	}
	c.outbound <- command{
		channel: string(channelID),
		sender:  c.username,
		id:      CHNS_FILE,
	}

	return nil
}

// recvFIle - function that allow receive a file sended by a user.
//	return: []byte array with content of the file in success, error
//	in case of fail.
func (c *client) recvFile() ([]byte, error) {
	var file = make([]byte, 10000)
	c.conn.Write([]byte("Receiving file...\n"))
	for {
		_, err := c.conn.Read(file)
		//c.conn.SetReadDeadline(time.Now().Add(1 * time.Minute))
		if err != nil {
			fmt.Printf("conn.read() method execution error, error is:% v\n", err)
			c.conn.Write([]byte("File cannot be read, verify filename and channelname!\n"))
		} else {
			fmt.Print("File receive\n")
			return file, nil
		}
	}

}

// err - function to send and error to the user.
//	return: an error.
func (c *client) err(e error) {
	c.conn.Write([]byte("ERR: " + e.Error() + "\n"))
}
