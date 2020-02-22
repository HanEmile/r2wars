package main

// Config defines the meta config
type Config struct {

	// Arch defines the architecture the battle should run in
	Arch string

	// Bits defines the bitness
	Bits int

	// Memsize defines the arena size
	Memsize int

	// MaxProgSize defines the maximal bot size
	MaxProgSize int

	// Bots defines a list of bots to take part in the battle
	Bots []Bot
}

// Bot defines a bot
type Bot struct {

	// Path defines the path to the source of the bot
	Path string

	// Source defines the source of the bot after being compiled with rasm2
	Source string

	// Addr defines the initial address the bot is placed at
	Addr int

	// Regs defines the state of the registers of the bot
	// It is used to store the registers after each round and restore them in the
	// next round when the bot's turn has come
	Regs string
}
