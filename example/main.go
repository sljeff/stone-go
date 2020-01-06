package main

import (
	"github.com/alecthomas/repr"
	"github.com/sljeff/stone-go"
)

func main() {
	sm, err := stone.ParseFile("config.stone")
	if err != nil {
		panic(err)
	}
	repr.Println(sm)
}
