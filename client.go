package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type client struct {
	conn     net.Conn
	nick     string
	idUser   int
	room     *room
	commands chan<- command
}

func (c *client) readInput() {
	for {
		msg, err := bufio.NewReader(c.conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])

		switch cmd {
		case "/nick":
			c.commands <- command{
				id:     CMD_NICK,
				client: c,
				args:   args,
			}
		case "/join":
			c.commands <- command{
				id:     CMD_JOIN,
				client: c,
				args:   args,
			}
		case "/rooms":
			c.commands <- command{
				id:     CMD_ROOMS,
				client: c,
			}
		case "/msg":
			c.commands <- command{
				id:     CMD_MSG,
				client: c,
				args:   args,
			}
		case "/quit":
			c.commands <- command{
				id:     CMD_QUIT,
				client: c,
			}
		case "/register":
			c.commands <- command{
				id:     CMD_REGISTER,
				client: c,
				args:   args,
			}
		case "/login":
			c.commands <- command{
				id:     CMD_LOGIN,
				client: c,
				args:   args,
			}
		case "/play":
			c.commands <- command{
				id:     CMD_PLAY,
				client: c,
			}
		case "/move":
			c.commands <- command{
				id:     CMD_MOVE,
				client: c,
				args:   args,
			}
		case "/history":
			c.commands <- command{
				id:     CMD_HISTORY,
				client: c,
			}
		case "/rate":
			c.commands <- command{
				id:     CMD_RATE,
				client: c,
				args:   args,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

func (c *client) err(err error) {
	c.conn.Write([]byte("err: " + err.Error() + "\n"))
}

func (c *client) msg(msg string) {
	c.conn.Write([]byte("> " + msg + "\n"))
}
