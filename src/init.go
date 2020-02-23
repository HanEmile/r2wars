package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	r2pipe "github.com/radare/r2pipe-go"
	"github.com/sirupsen/logrus"
)

func buildBots(config *Config) {

	logrus.Info("Building all bots")

	// build all the bots
	for i := 0; i < config.AmountOfBots; i++ {
		buildBot(i, config)
	}
}

// buildBot builds the bot located at the given path.
func buildBot(i int, config *Config) {

	logrus.Debugf("Building bot %d", i)

	// open radare without input for building the bot
	r2p1, err := r2pipe.NewPipe("--")
	if err != nil {
		panic(err)
	}
	defer r2p1.Close()

	// Compile a warrior using rasm2
	//
	// This uses the given architecture, the given bitness and the given path in
	// rasm2 to compile the bot
	botPath := config.Bots[i].Path
	radareCommand := fmt.Sprintf("rasm2 -a %s -b %d -f %s", config.Arch, config.Bits, botPath)
	botSource := r2cmd(r2p1, radareCommand)

	config.Bots[i].Source = botSource
}

// init initializes the arena
func initArena(config Config) *r2pipe.Pipe {
	log.Println("[+] Initializing the arena")
	log.Printf("[ ] Allocating %d bytes of memory...", config.Memsize)

	// alocate memory
	r2p, err := r2pipe.NewPipe(fmt.Sprintf("malloc://%d", config.Memsize))
	if err != nil {
		panic(err)
	}
	log.Println("[+] Memoy successfully allocated \\o/")

	// define the architecture and the bitness
	_ = r2cmd(r2p, fmt.Sprintf("e asm.arch = %s", config.Arch))
	_ = r2cmd(r2p, fmt.Sprintf("e asm.bits = %d", config.Bits))

	// enable colors
	_ = r2cmd(r2p, "e scr.color = true")

	log.Println("[+] Initializing the ESIL VM")
	// initialize ESIL VM state
	_ = r2cmd(r2p, "aei")

	// initialize ESIL VM stack
	_ = r2cmd(r2p, "aeim")

	// return the pipe
	return r2p

}

// getRandomOffsets returns random offsets for all bots
// This is used to get the offset the bots are initially placed in
func getRandomOffsets(config Config) []int {

	var amountOfBots int = len(config.Bots)
	var offsets []int
	var roundCounter int = 0

	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	for {
		// define an integer array to store the random offsets in
		//var offsets []int = []int{}

		// define a random address
		address := rand.Intn(config.Memsize - config.MaxProgSize)

		// for all bots, try to generate another random address after the intially
		// generated address and test if it fits in memory
		for i := 0; i < amountOfBots; i++ {
			offsets = append(offsets, address)

			// generate a random value in range [maxProgSize, maxProgSize + 300)
			address += rand.Intn(config.MaxProgSize+300) + config.MaxProgSize

			// if there is not enough memory remaining after the last generated
			// address, start from be beginning
			if address+config.MaxProgSize > config.Memsize {
				roundCounter++
				continue
			}
		}

		// if enough addresses have been generated, break out of the for loop
		break
	}

	log.Printf("[+] Initial bot positions found after %d trues", roundCounter)

	return offsets
}

// place the bot in the arena at the given address
func placeBot(r2p *r2pipe.Pipe, bot Bot, address int) {
	_ = r2cmd(r2p, fmt.Sprintf("wx %s @ %d", bot, address))
}
