package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

// parseConfig parses the config needed to start the game
func parseConfig() Config {

	// bot configs
	arch := flag.String("arch", "x86", "bot architecture (mips|arm|x86)")
	bits := flag.Int("bits", 32, "bot bitness (8|16|32|64)")
	maxProgSize := flag.Int("maxProgSize", 64, "the maximum bot size")
	memPerBot := flag.Int("memPerBot", 512, "the amount of memory each bot should add to the arena")
	gameRoundDuration := flag.Duration("t", 250, "The duration of a round")

	v := flag.Bool("v", false, "info")
	vv := flag.Bool("vv", false, "debug")
	vvv := flag.Bool("vvv", false, "trace")

	// parse the flags
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
