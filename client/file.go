package main

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

// sendFile - function to send a file to the server.
//	(string) fileName - filenamee to giving send.
//	(net.Conn) conn - living conection with the server.
func sendFile(fileName string, conn net.Conn) {

	var ph string
	fmt.Printf("Write directory of the file:\n")
	fmt.Scanf("%s", &ph)
	path := filepath.Join(ph, fileName)
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	data := make([]byte, 5000)
	_, err = file.Read(data)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	cleanData := bytes.Trim(data, "\x00")
	conn.Write(cleanData)
}

// createFile - function to create a new file coming from the server.
//	([]byte) fileContent - content coming from the server to create a
//	new file in the current directory.
//	(string) fileName - file name to create the new file.
//	return: error in case of fail, nil in success.
func createFile(fileContent []byte, fileName string) error {
	//path := filepath.Join("./", fileName)
	fileName = strings.TrimSpace(fileName)
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("%v", err)
		return err
	}
	cleanContent := bytes.Trim(fileContent, "\x00")
	content := string(cleanContent)
	_, err2 := file.WriteString(content)
	if err2 != nil {
		fmt.Printf("%v", err)
		return err2
	}
	defer file.Close()
	fmt.Printf("File %v was created!\n", fileName)
	return nil
}
