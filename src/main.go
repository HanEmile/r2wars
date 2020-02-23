package main

func main() {

	// initialize the game by parsing the config, defining the bots, building the
	// bots and generating random offsets where the bots should be placed in
	// memory
	config := parseConfig()
	defineBots(&config)
	buildBots(&config)
	genRandomOffsets(&config)

	// initialize the arena (allocate memory + initialize the ESIL VM & stack)
	r2p := initArena(&config)

	// place the bots in the arena
	placeBots(r2p, &config)

	// if an error occurs (interrupt, ioerror, trap, ...), the ESIL VM should set
	// a flag that can be used to determine if a player has died
	defineErrors(r2p)

	// run the actual game
	runGame(r2p, &config)

	r2p.Close()
}
