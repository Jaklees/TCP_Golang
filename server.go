package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms map[string]*room
	commands chan command
}

func newServer() *server {
	return &server{
		rooms: make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) run(){
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NICK:
			s.nick(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.lobbies(cmd.client, cmd.args)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client, cmd.args)
		}
	}
}
 func (s *server) newClient(conn net.Conn) {
 	log.Printf("New Client has connected: %s, Welcome!", conn.RemoteAddr().String())

 	c := &client{
			conn: conn,
			nick: "anon",
			commands: s.commands,
 	}
 	c.readInput()
}

func (s *server) nick(c *client, args []string) {
	c.nick = args[1]
	c.msg(fmt.Sprintf("Nick name set: %s ",c.nick))
}

func (s *server) join(c *client, args []string) {
	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name: roomName,
			members: make(map[net.Addr]*client),
		}
		s.rooms[roomName] = r
	}
	r.members[c.conn.RemoteAddr()] = c

	s.quitRoom(c)
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s joined",c.nick))
	c.msg(fmt.Sprintf("Welcome to %s", r.name))
}

func (s *server) lobbies(c *client, args []string){
	var rooms []string
	for name := range s.rooms {
		rooms = append(rooms, name)
	}
	c.msg(fmt.Sprintf("Rooms open: %s",strings.Join(rooms, ", ")))
}

func (s *server) msg(c *client, args []string) {
	if c.room == nil {
		c.err(errors.New("error: no room"))
		return
	}
	c.room.broadcast(c,c.nick+":  "+strings.Join(args[1:len(args)]," "))
}

func (s *server) quit(c *client, args []string) {
	log.Printf("Client has disconnected: %s", c.conn.RemoteAddr().String())

	s.quitRoom(c)

	c.msg("Goodbye! Hope to see you again")
	c.conn.Close()
}

func (s *server) quitRoom(c *client){
	if c.room != nil {
		delete(c.room.members,c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left, Goodbye!", c.nick))
	}
}