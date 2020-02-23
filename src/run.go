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

// runGame actually runs the game (surprise!)
func runGame(r2p *r2pipe.Pipe, config *Config) {

	// start the competition
	var botid int = 0
	var round int = 0
	for true {

		// load the registers
		r2cmd(r2p, config.Bots[botid].Regs)

		// Step
		stepIn(r2p)

		// store the regisers
		registers := r2cmd(r2p, "aerR")
		registersStripped := strings.Replace(registers, "\n", ";", -1)
		config.Bots[botid].Regs = registersStripped

		logrus.Info(arena(r2p, config, botid, round))

		if dead(r2p, botid) == true {
			logrus.Warnf("DEAD (round %d)", round)
			os.Exit(1)
		}

		// switch players, if the new botid is 0, a new round has begun
		botid = switchPlayer(r2p, botid, config)
		if botid == 0 {
			round++
		}

		// sleep only a partial of the total round time, as a round is made up of
		// the movements of multiple bots
		time.Sleep(config.GameRoundDuration / time.Duration(config.AmountOfBots))
	}
}
