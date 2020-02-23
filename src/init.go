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
func initArena(config *Config) *r2pipe.Pipe {

	logrus.Info("Initializing the arena")
	logrus.Debugf("Allocating %d bytes of memory...", config.Memsize)

	// allocate memory
	r2p, err := r2pipe.NewPipe(fmt.Sprintf("malloc://%d", config.Memsize))
	if err != nil {
		panic(err)
	}

	logrus.Info("Memoy successfully allocated")

	// define the architecture and the bitness
	_ = r2cmd(r2p, fmt.Sprintf("e asm.arch = %s", config.Arch))
	_ = r2cmd(r2p, fmt.Sprintf("e asm.bits = %d", config.Bits))

	// enable colors
	// _ = r2cmd(r2p, "e scr.color = 0")
	_ = r2cmd(r2p, "e scr.color = 3")
	_ = r2cmd(r2p, "e scr.color.args = true")
	_ = r2cmd(r2p, "e scr.color.bytes = true")
	_ = r2cmd(r2p, "e scr.color.grep = true")
	_ = r2cmd(r2p, "e scr.color.ops = true")
	_ = r2cmd(r2p, "e scr.bgfill = true")
	_ = r2cmd(r2p, "e scr.color.pipe = true")
	_ = r2cmd(r2p, "e scr.utf8 = true")

	// hex column width
	_ = r2cmd(r2p, "e hex.cols = 32")

	// initialize ESIL VM state
	logrus.Debug("Initializing the ESIL VM")
	_ = r2cmd(r2p, "aei")

	// initialize ESIL VM stack
	logrus.Debug("Initializing the ESIL Stack")
	_ = r2cmd(r2p, "aeim")

	// return the pipe
	return r2p
}

// genRandomOffsets returns random offsets for all bots
// This is used to get the offset the bots are initially placed in
func genRandomOffsets(config *Config) {

	logrus.Info("Generating random bot offsets")

	// define the amount of bots, an array to store the offsets in and a counter
	// to store the amount of tries it took to find a random positon for the bots
	var amountOfBots int = len(config.Bots)
	var offsets []int
	var roundCounter int = 0

	// seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())

	for {

		// reset the offsets array
		offsets = []int{}

		// define a random address
		// | ------------------------------------- | ----- |
		// | generate an address in this space
		address := rand.Intn(config.Memsize - config.MaxProgSize)
		logrus.Tracef("%d", address)

		// for all bots, try to generate another random address after the intially
		// generated address and test if it fits in memory
		for i := 0; i < amountOfBots; i++ {

			// append the address to the offsets array
			offsets = append(offsets, address)

			// define a min and max bound
			//
			// | ------|-|----------------------------------|-|
			// |       | |           in this space          | |
			// a       b c                                  d e
			//
			// a = 0x0
			// b = address
			// c = address + config.MaxProcSize (min)
			// d = config.Memsize - config.MaxProgSize (max)
			// e = config.Memsize
			min := address + config.MaxProgSize
			max := config.Memsize - config.MaxProgSize

			// if the new minimum bound is bigger or equal to the maximum bound,
			// discard this try and start with a fresh new initial address
			if min >= max {
				roundCounter++
				break
			}

			// generate a new address in the [min, max) range defined above
			address = rand.Intn(max-min) + min
			logrus.Tracef("%d", address)

			// If there isn't enough space inbetween the address and the biggest
			// possible address, as in, the biggest possible bot can't fit in that
			// space, discard and start with a new fresh initial address
			if (config.Memsize-config.MaxProgSize)-address < config.MaxProgSize {
				roundCounter++
				break
			}
		}

		// if the needed amount of offsets has been found, break out of the infinite loop
		if len(offsets) == amountOfBots {
			break
		}
	}

	logrus.Infof("Initial bot positions found after %d tries", roundCounter)

	// debug print all offsets
	var fields0 logrus.Fields = make(logrus.Fields)
	for i := 0; i < len(offsets); i++ {
		fields0[fmt.Sprintf("%d", i)] = offsets[i]
	}
	logrus.WithFields(fields0).Debug("Offsets")

	// shuffle the offsets
	rand.Shuffle(len(offsets), func(i, j int) {
		offsets[i], offsets[j] = offsets[j], offsets[i]
	})

	// debug print the shuffled offsets
	var fields1 logrus.Fields = make(logrus.Fields)
	for i := 0; i < len(offsets); i++ {
		fields1[fmt.Sprintf("%d", i)] = offsets[i]
	}
	logrus.WithFields(fields1).Debug("Shuffled offsets")

	config.RandomOffsets = offsets
}

// place the bot in the arena at the given address
func placeBot(r2p *r2pipe.Pipe, bot Bot, address int) {
	_ = r2cmd(r2p, fmt.Sprintf("wx %s @ %d", bot, address))
}
