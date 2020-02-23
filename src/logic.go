package main

import (
	"log"

	"github.com/radare/r2pipe-go"
)

func dead(r2p *r2pipe.Pipe, botid int) bool {
	status := r2cmd(r2p, "?v 1+theend")

	if status != "" && status != "0x1" {
		log.Printf("[!] Bot %d has died", botid)
		return true
	}
	return false
}
