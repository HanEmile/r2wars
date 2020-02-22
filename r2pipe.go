package main

import (
	"fmt"

	"github.com/radare/r2pipe-go"
)

func main() {
	// open a file
	// $ r2 ...
	r2p, err := r2pipe.NewPipe("/nix/store/xhwhakb1zcf5wl2a8575gcrnmbbqihm2-busybox-1.30.1/bin/ls")
	if err != nil {
		panic(err)
	}
	defer r2p.Close()

	// send a command
	// [0x004087e0]> ...
	buf1, err := r2p.Cmd("?E Hello World")
	if err != nil {
		panic(err)
	}

	// print the result of the first command
	fmt.Println(buf1)
}
