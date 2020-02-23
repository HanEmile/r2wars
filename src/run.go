package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/radare/r2pipe-go"
	"github.com/sirupsen/logrus"
)

// StepIn steps in and stores the state of the registers for the given bot
func stepIn(r2p *r2pipe.Pipe) {
	_ = r2cmd(r2p, "aes")

	// store the regisers
	registers := r2cmd(r2p, "aerR")
	registersStripped := strings.Replace(registers, "\n", ";", -1)
	return registersStripped
}

func switchPlayer(r2p *r2pipe.Pipe, currentPlayer int, config Config) int {

	// calculate the index of the nextPlayer
	nextPlayer := (currentPlayer + 1) % len(config.Bots)

	// restore the registers to the state of the next bot
	r2cmd(r2p, config.Bots[nextPlayer].Regs)

	return nextPlayer
}

func user(r2p *r2pipe.Pipe, id int, registers string, config Config) string {
	var res string

	res += "\x1b[2J\x1b[0;0H"
	res += fmt.Sprintf("USER %d\n", id)
	res += fmt.Sprintf("%s\n", r2cmd(r2p, "aer"))
	res += "+++\n"
	//res += fmt.Sprintf("%s\n", r2cmd(r2p, fmt.Sprintf("%s %d @ 0\n", "prc=f", config.Memsize)))
	res += fmt.Sprintf("%s\n", r2cmd(r2p, fmt.Sprintf("%s %d @ 0\n", "prx", config.Memsize)))
	res += "+++\n"
	res += fmt.Sprintf("%s\n", r2cmd(r2p, "pxw 32 @r:SP"))

	r2command := fmt.Sprintf("pd %d @ %d", len(config.Bots[id].Source)/2, config.Bots[id].Addr)
	r2cmdString := r2cmd(r2p, r2command)
	res += fmt.Sprintf("%s\n", r2cmdString)

	return res
}
