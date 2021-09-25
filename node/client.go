package node

import (
	"Network/date"
	"fmt"
	"os"
	"strings"
)

func handleClient(node *Node) { //Клиент
	for {
		message := date.InputString()
		splited := strings.Split(message, " ") //Берём дынные и разбиваем

		switch splited[0] { //Команды клиента
		case "/exit":
			os.Exit(0)
		case "/connect":
			node.HandShake(splited[1], true)
		case "/network":
			node.PrintNetwork()
		case "/m":
			if len(splited) < 3 {
				fmt.Println("Не верное кол-во аргументов")
				continue
			}
			node.SendMessageTo(splited[1], []byte(splited[2]))
		default:
			node.SendMessageToAll([]byte(message))
		}
	}
}

func (node *Node) PrintNetwork() { //Ввывод всех подключений
	for addr := range node.Connections {
		fmt.Println(node.Connections[addr].Name + "|" + addr)
	}
}
