# genetic_r2_bots

building r2 bots using genetic algorithms

## r2wars

The goal of r2wars is to place two bots in as shared memory space and let
them battle against each other until one of them cannot run anymore, because
of some kind of broken instruction.

A more informal README can be found [here](https://github.com/radareorg/radare2-extras/tree/master/r2wars).

## Usage

You'll probably first of all want to simply play with the two provided bots. In order to do so, run the game like this:

```go
go run ./... -t 1s -v ./bots/warrior0.asm ./bots/warrior1.asm
```

This runs the game with a round duration of 1 second and an info verbosity
level using the two provided bots. You can attach more bots if you'd like,
each bot increases the arena size by 512 bytes by default.

You can tweak most of the settings as displayed in the help:

```go
$ go run ./... -h
Usage of src:
  -arch string
    	bot architecture (mips|arm|x86) (default "x86")
  -bits int
    	bot bitness (8|16|32|64) (default 32)
  -maxProgSize int
    	the maximum bot size (default 64)
  -memPerBot int
    	the amount of memory each bot should add to the arena (default 512)
  -t duration
    	The duration of a round (default 250ns)
  -v	info
  -vv
    	debug
  -vvv
    	trace
```