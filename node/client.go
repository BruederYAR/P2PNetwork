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
		case "/help":
			fmt.Println("^/exit", "Выход из программы")
			fmt.Println("^/connect", "Присоеденится к узлу. Ключи: 1)адрес(ip:port)")
			fmt.Println("^/network", "Вывести все присоедененные узлы и собственный адрес")
			fmt.Println("^/m", "Отправить сообщение. Ключи: 1)адрес(ip:port или имя) 2)сообщение")
		case "/connect":
			node.HandShake(splited[1], true)
		case "/network":
			node.PrintNetwork()
		case "/modules":
			for i,j := range node.Input.Modules{
				fmt.Println("Module", i, "\n", "Desk", j.Desk, "\n", "Path", j.Path)
				for _,info := range j.Cmds{
					fmt.Println("   ", info.Cmd, info.Desk)
				}
			}

		case "/m":
			if len(splited) < 3 {
				fmt.Println("Не верное кол-во аргументов")
				continue
			}
			node.SendMessageTo(splited[1], []byte(splited[2]))
		case "/mm":
			if len(splited) < 3{
				fmt.Println("Не верное кол-во аргументов")
				continue
			}
			node.ModuleRequestTo(splited[1], splited[2], splited[3])
		case "/mc":
			if len(splited) < 2{
				fmt.Println("Не верное кол-во аргументов")
				continue
			}
			node.CommandRequestTo(splited[1], []byte(splited[2]))
		default:
			node.SendMessageToAll([]byte(message))
		}
	}
}

func (node *Node) PrintNetwork() { //Ввывод всех подключений
	fmt.Println("local address " + node.Address.IP + node.Address.Port)
	for addr := range node.Connections {
		fmt.Println(node.Connections[addr].Name + "|" + addr)
	}
}

