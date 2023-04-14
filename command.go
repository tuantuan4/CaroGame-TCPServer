package main

type commandID int

const (
	CMD_NICK commandID = iota
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
	CMD_LOGIN
	CMD_REGISTER
	CMD_PLAY
	CMD_MOVE
	CMD_HISTORY
	CMD_RATE
)

type command struct {
	id     commandID
	client *client
	args   []string
}
