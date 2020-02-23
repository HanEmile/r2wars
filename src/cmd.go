package main

import (
	"github.com/radare/r2pipe-go"
	"github.com/sirupsen/logrus"
)

func r2cmd(r2p *r2pipe.Pipe, input string) string {

	logrus.Tracef("> %s", input)

	// send a command
	buf1, err := r2p.Cmd(input)
	if err != nil {
		panic(err)
	}

	// return the result of the command as a string
	return buf1
}
