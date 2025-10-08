package server

import (
	"io"
	"log"
	"redis-clone/core"
	"syscall"
)

type FDcomm struct {
	Fd int
}

func (f *FDcomm) Read(b []byte) (int, error) {
	return syscall.Read(f.Fd, b)
}
func (f *FDcomm) Write(b []byte) (int, error) {
	return syscall.Write(f.Fd, b)
}
func readCommands(c io.ReadWriter) (core.Cmds, error) {
	//TODO: can only read a buffer of 4096 need to write a repeated read till EOF/delimiter logic

	buffer := make([]byte, 4096)
	n, err := c.Read(buffer)
	if err != nil {
		return nil, err
	}
	values, err := core.Decode(buffer[:n])
	log.Println("tarun--3", values)

	if err != nil {
		log.Println("tarun--4", err)

		return nil, err
	}

	var commands core.Cmds = make(core.Cmds, 0)
	log.Println("tarun--2", values)

	for _, v := range values {
		tokens, err := core.DecodeArrayString(v)
		if err != nil {
			return commands, err
		}
		log.Println("tarun--1", tokens)
		Cmd := &core.Cmd{
			Cmd:  tokens[0],
			Args: tokens[1:],
		}
		commands = append(commands, Cmd)

	}

	return commands, err
}

func writeCommand(c io.ReadWriter, val []byte) error {
	_, err := c.Write(val)

	return err
}
func writeErrorCommand(c io.ReadWriter, err error) error {
	_, errWrite := c.Write(core.EncodeError(err))
	return errWrite
}
