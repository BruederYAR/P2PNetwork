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

	//fmt.Println(Modules)

	//Отправку еденичных сообщений определенному узлу
	//title с запросом. В котором передатся команда. Отправляется результат как title data
	//Выдовать по запросу
}

type all struct{
	input date.Input
}

func main() {
	input := date.NewInput()
	node.Start(input)
}
