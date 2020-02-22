package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

var (
	verbose *bool
)

func main() {

	log.Println("---")
	log.Println("[i] Parse the config")
	config := parseConfig()

	// build the bots
	log.Println("---")
	log.Println("[i] Build the bots")

	for i := 0; i < 2; i++ {
		bot := buildBot(config, "bots/warrior.asm")
		config.Bots = append(config.Bots, bot)
	}
	// bot2 := buildBot(config, "bots/warrior.asm")
	// config.Bots = append(config.Bots, bot2)
	// bot3 := buildBot(config, "bots/warrior.asm")
	// config.Bots = append(config.Bots, bot3)

	// initialize the arena
	log.Println("---")
	log.Println("[i] Initialize the Arena")
	r2p := initArena(config)

	randomOffsets := getRandomOffsets(config)

	// place each bot in the arena
	log.Println("---")
	log.Println("[i] Place the bots")
	for bot := 0; bot < len(config.Bots); bot++ {
		// Place the bot in the arena
		log.Printf("[i] Placing bot %d", bot)
		address := randomOffsets[bot]
		placeBot(r2p, config.Bots[bot], address)

		// store the initial address of the bot int the struct field
		config.Bots[bot].Addr = address

		// define the instruction point and the stack pointer
		log.Printf("[i] setting up the PC and SP for bot %d", bot)
		_ = r2cmd(r2p, fmt.Sprintf("aer PC=%d", address))
		_ = r2cmd(r2p, fmt.Sprintf("aer SP=SP+%d", address))

		// dump the registers of the user for being able to switch inbetween them
		initialRegisers := strings.Replace(r2cmd(r2p, "aerR"), "\n", ";", -1)
		config.Bots[bot].Regs = initialRegisers

		// print the instruction point and the stack pointer
		botStackPointer := r2cmd(r2p, "aerR~esp[2]")
		log.Printf("[i] bot %d esp = %s", bot, botStackPointer)
		botInstructionPointer := r2cmd(r2p, "aerR~eip[2]")
		log.Printf("[i] bot %d eip = %s", bot, botInstructionPointer)
	}

	// handle errors in esil
	_ = r2cmd(r2p, "e cmd.esil.todo=f theend=1")
	_ = r2cmd(r2p, "e cmd.esil.trap=f theend=1")
	_ = r2cmd(r2p, "e cmd.esil.intr=f theend=1")
	_ = r2cmd(r2p, "e cmd.esil.ioer=f theend=1")
	_ = r2cmd(r2p, "f theend=0")

	fmt.Println(r2cmd(r2p, fmt.Sprintf("b %d", config.Memsize)))

	// start the competition
	i := 0
	for true {

		// Step, then print the users registers
		registers := stepIn(r2p)
		config.Bots[i].Regs = registers
		fmt.Println(user(r2p, i, registers, config))

		// switch players
		i = switchPlayer(r2p, i, config)

		// sleepti
		time.Sleep(100 * time.Millisecond)
	}

	r2p.Close()
}
