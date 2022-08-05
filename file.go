package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

type file struct {
	name string
	path string
}

// newFile - function to create a new register of file in
// the server hub.
//	(string) fileName - file name to create.
//	(string) filePath - path where the new file was create.
func newFile(fileName string, filePath string) *file {
	return &file{
		name: fileName,
		path: filePath,
	}
}

// createFile - Function to create a new file in the files directory
// of the server.
//	(string) fileName - name of the file to create.
// ([]byte) file - content of the file.
func createFile(fileName string, fileContent []byte) {
	//fmt.Printf("file was created!\n")
	path := filepath.Join("files", fileName)
	file, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	cleanContent := bytes.Trim(fileContent, "\x00")
	content := string(cleanContent)
	_, err2 := file.WriteString(content)
	if err2 != nil {
		log.Fatal(err2)
	}
	defer file.Close()
}

// chargeFiles - function to charge all existing files in the
// server to the server hub.
func chargeFiles(hub hub) {
	files, err := ioutil.ReadDir(hub.rootDir)
	if err != nil {
		fmt.Printf("Error: %v", err)
		return
	}
	for _, file := range files {
		path := "/" + hub.rootDir + "/" + file.Name()
		newFile := newFile(file.Name(), path)
		hub.files[file.Name()] = newFile
	}
}
