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
}

// switchPlayer returns the id of the next Player
func switchPlayer(r2p *r2pipe.Pipe, currentPlayer int, config *Config) int {

	// calculate the index of the nextPlayer
	nextPlayer := (currentPlayer + 1) % config.AmountOfBots

	return nextPlayer
}

func arena(r2p *r2pipe.Pipe, config *Config, id, gen int) string {
	var res string = ""

	// clear the screen
	res += "\x1b[2J\x1b[0;0H"
	// res += fmt.Sprintf("%s\n", r2cmd(r2p, "?eg 0 0"))

	// print some general information such as the current user and the round the
	// game is in
	ip := fmt.Sprintf("%s\n", r2cmd(r2p, "aer~eip"))
	res += fmt.Sprintf("Round: %d \t\t User: %d \t\t ip: %s\n", gen, id, ip)

	// print the memory space
	res += fmt.Sprintf("%s\n", r2cmd(r2p, "pxa 0x400 @ 0"))
	// res += fmt.Sprintf("%s\n", r2cmd(r2p, fmt.Sprintf("pd 0x10 @ %d", config.Bots[id].Addr)))

	// res += fmt.Sprintf("%s\n", r2cmd(r2p, "prc 0x200 @ 0"))

	return res
}
