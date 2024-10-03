package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/radareorg/r2pipe-go"
	"github.com/sirupsen/logrus"
)

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

	// AmountOfBots defines the amount of bots taking part in the tournament
	AmountOfBots int

	// RandomOffsets defines the offset in memory where the bots should be placed
	RandomOffsets []int

	// GameRoundTime defines the length of a gameround
	GameRoundDuration time.Duration
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

func parseConfig() Config {
	arch := flag.String("arch", "x86", "bot architecture (mips|arm|x86)")
	bits := flag.Int("bits", 32, "bot bitness (8|16|32|64)")
	maxProgSize := flag.Int("maxProgSize", 64, "the maximum bot size")
	memPerBot := flag.Int("memPerBot", 512, "the amount of memory each bot should add to the arena")
	gameRoundDuration := flag.Duration("t", 250*time.Millisecond, "The duration of a round")

	v := flag.Bool("v", false, "info")
	vv := flag.Bool("vv", false, "debug")
	vvv := flag.Bool("vvv", false, "trace")

	flag.Parse()

	if *v == true {
		logrus.SetLevel(logrus.InfoLevel)
	} else if *vv == true {
		logrus.SetLevel(logrus.DebugLevel)
	} else if *vvv == true {
		logrus.SetLevel(logrus.TraceLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	// parse all trailing command line arguments as path to bot sourcecode
	amountOfBots := flag.NArg()

	memsize := amountOfBots * *memPerBot

	logrus.WithFields(logrus.Fields{
		"mem per bot":  *memPerBot,
		"amountOfBots": amountOfBots,
		"memsize":      memsize,
	}).Infof("Loaded config")

	// define a config to return
	config := Config{
		Arch:              *arch,
		Bits:              *bits,
		Memsize:           memsize,
		MaxProgSize:       *maxProgSize,
		AmountOfBots:      amountOfBots,
		GameRoundDuration: *gameRoundDuration,
	}

	return config
}

// define bots defines the bots given via command line arguments
func defineBots(config *Config) {

	logrus.Info("Defining the bots")

	// define a list of bots by parsing the command line arguments one by one
	var bots []Bot
	for i := 0; i < config.AmountOfBots; i++ {
		bot := Bot{
			Path: flag.Arg(i),
		}
		bots = append(bots, bot)
	}

	config.Bots = bots
}

func r2cmd(r2p *r2pipe.Pipe, input string) string {

	logrus.Tracef("> %s", input)

	// send a command
	buf1, err := r2p.Cmd(input)
	if err != nil {
		panic(err)
	}

	// return the result of the command as a string
	return buf1
}

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

	logrus.Info("Memory successfully allocated")

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
	_ = r2cmd(r2p, fmt.Sprintf("wx %s @ %d", bot.Source, address))
}

func placeBots(r2p *r2pipe.Pipe, config *Config) {

	logrus.Info("Placing the bots in the arena")

	// place each bot in the arena
	for bot := 0; bot < len(config.Bots); bot++ {

		// get the address where the bot should be placed
		address := config.RandomOffsets[bot]

		// Place the bot in the arena
		logrus.Debugf("[i] Placing bot %d at %d", bot, address)
		placeBot(r2p, config.Bots[bot], address)

		logrus.Debugf("\n%s", r2cmd(r2p, fmt.Sprintf("pd 0x8 @ %d", address)))

		// store the initial address of the bot in the according struct field
		config.Bots[bot].Addr = address

		// define the instruction point and the stack pointer
		_ = r2cmd(r2p, fmt.Sprintf("aer PC=%d", config.Bots[bot].Addr))
		_ = r2cmd(r2p, fmt.Sprintf("aer SP=SP+%d", config.Bots[bot].Addr))

		// dump the registers of the bot for being able to switch inbetween them
		// This is done in order to be able to play one step of each bot at a time,
		// but sort of in parallel
		initialRegisers := strings.Replace(r2cmd(r2p, "aerR"), "\n", ";", -1)
		config.Bots[bot].Regs = initialRegisers
	}
}

func defineErrors(r2p *r2pipe.Pipe) {
	// handle errors in esil
	_ = r2cmd(r2p, "e cmd.esil.todo=f theend=1")
	_ = r2cmd(r2p, "e cmd.esil.trap=f theend=1")
	_ = r2cmd(r2p, "e cmd.esil.intr=f theend=1")
	_ = r2cmd(r2p, "e cmd.esil.ioer=f theend=1")
	_ = r2cmd(r2p, "f theend=0")
}

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

func dead(r2p *r2pipe.Pipe, botid int) bool {
	status := strings.TrimSpace(r2cmd(r2p, "?v theend"))
	// fixme: on Windows, we sometimes get output *from other calls to r2*

	if status == "0x1" {
		logrus.Warnf("[!] Bot %d has died", botid)
		return true
	}
	if status != "0x0" {
		logrus.Warnf("[!] Got invalid status '%s' for bot %d", status, botid)
	}
	return false
}

// The Windows terminal doesn't render ANSI escape codes by default,
// but at least, it supports opting in since Windows 10 1709.
// We could use windigo to call kernel32.SetConsoleMode, but this still
// wouldn't be enough as logrus seems to escape them. But if we force logrus to
// use color, it seems to call this function for us. :)
// see https://learn.microsoft.com/en-us/windows/console/setconsolemode
// see https://superuser.com/a/1300251/329759
func forceColor() {
	logrus.Infof("Running on Windows; forcing color")
	// make logrus use color
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
		FullTimestamp: false,
	})
}

func main() {
	if runtime.GOOS == "windows" {
		forceColor()
	}
	fmt.Println("\x1b[33mhi\x1b[0m")

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
