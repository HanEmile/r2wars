package main

import (
	"fmt"
	"strings"

	"github.com/radare/r2pipe-go"
)

func stepIn(r2p *r2pipe.Pipe) string {
	_ = r2cmd(r2p, "aes")
	registers := r2cmd(r2p, "aerR")
	registersStripped := strings.Replace(registers, "\n", ";", -1)
	return registersStripped
}

func switchPlayer(currentPlayer int, config Config) int {
	return (currentPlayer + 1) % len(config.Bots)
}

func user(r2p *r2pipe.Pipe, id int, registers string, config Config) string {
	var res string

	// res += "\x1b[2J\x1b[0;0H"
	res += fmt.Sprintf("USER %d\n", id)
	res += fmt.Sprintf("%s\n", r2cmd(r2p, "aer"))
	res += "+++\n"
	res += fmt.Sprintf("%s\n", r2cmd(r2p, fmt.Sprintf("%s %d @ 0\n", "prx", config.Memsize)))
	res += "+++\n"
	res += fmt.Sprintf("%s\n", r2cmd("pxw 32 @r:SP"))
	res += fmt.Sprintf("%s\n", r2cmd("pD %d @ %s"%(size[uidx], orig[uidx]))

	return res
}
