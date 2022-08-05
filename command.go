package main

type ID int

const (
	REG ID = iota
	JOIN
	LEAVE
	D_FILE
	S_FILE
	L_FILES
	USRS
	CHNS
	CHNS_FILE
)

type command struct {
	id       ID
	sender   string
	receiver string
	channel  string
	fileName string
	file     []byte
}
