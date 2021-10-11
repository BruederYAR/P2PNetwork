package main

import (
	"Network/node"
	"Network/date"
	"os"
)

func init() {
	if len(os.Args) < 3 { //Если длина аргументов при включении меньше 3 - вызываем панику
		panic("len args < 3")
	}
}

func main() {
	input := date.NewInput()
	node.Start(input)
}
