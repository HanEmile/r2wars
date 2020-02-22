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
}
