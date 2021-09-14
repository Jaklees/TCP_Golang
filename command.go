package main

import "net"

type commandID int

const (
	CMD_NICK commandID = iota // value starts at 1 and increments
	CMD_JOIN
	CMD_ROOMS
	CMD_MSG
	CMD_QUIT
)

type command struct {
	id commandID
	client *client
	args []string
}

type room struct {
	name string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, m := range r.members {
		if sender.conn.RemoteAddr() != addr {
			m.msg(msg)
		}
	}
}