package main

import "flag"

func parseConfig() Config {
	verbose = flag.Bool("v", false, "verbose output")
	flag.Parse()

	config := Config{
		Arch:        "x86",
		Bits:        32,
		Memsize:     1024,
		MaxProgSize: 64,
	}

	return config
}
