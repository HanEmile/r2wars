package main

import (
	"fmt"

	"github.com/radare/r2pipe-go"
)

func main() {
	// allocate 1024 bytes of memory
	r2p, err := r2pipe.NewPipe("malloc://1024")
	if err != nil {
		panic(err)
	}
	defer r2p.Close()

	// get a hexdump of the first 100 bytes allocated
	hexdump := r2cmd(r2p, "px 100")
	fmt.Println(hexdump)

	// compile a warrior using rasm2
	bot := r2cmd(r2p, "rasm2 -a x86 -b 32 -f bots/warrior.asm")
	fmt.Println(bot)
}
