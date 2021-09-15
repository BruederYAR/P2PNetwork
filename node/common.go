package node

import (
	"Network/date"
	"fmt"
	"net"
	"os"
	"strings"
)

type Node struct {
	Titles      map[int]string
	Types       map[int]string
	Connections map[string]string
	Input       *date.Input
	Address     Address
	Name        string //Имя узла
}

type Address struct {
	IPv4 string
	Port string
}

func Start(input *date.Input) {
	NewNode(os.Args[1], input).Run(handleServer, handleClient) //создаём ноду и закускаем узел
}

func NewNode(address string, input *date.Input) *Node { //Создание ноды
	splited := strings.Split(address, ":") //отделяем адрес от порта
	if len(splited) != 2 {
		return nil
	}
	return &Node{ //Возвращаем ноду, в которой карта подключений и текущий адрес
		Titles: map[int]string{
			0: "handshake",
			1: "date",
		},
		Types: map[int]string{
			0: "string",
			1: "json",
		},
		Input:       input,
		Connections: make(map[string]string),
		Address: Address{
			IPv4: splited[0],
			Port: ":" + splited[1],
		},
		Name: os.Args[2],
	}
}

func (node *Node) Run(handleServer func(*Node), handleClient func(*Node)) { //Выполняется запуск как сервер, так и клиента

	go handleServer(node)
	handleClient(node)
}





func (node *Node) ConnectTo(addresses []string, name string) { //Добавление в список подключений
	for _, addr := range addresses {
		if addr == "" {
			panic("Пустой адрес, всё плохо")
		}
		node.Connections[addr] = name

	}
}

func (node *Node) HandShake(address string, status bool) { //Рукопожатие при первом подключении
	var new_pack = date.Packege{
		From:  node.Address.IPv4 + node.Address.Port,
		To:    address,
		Name:  node.Name,
		Type:  node.Types[1],
		Date:  []byte{},
		Title: node.Titles[0],
	}

	if !status {
		fmt.Println("HandShake from", new_pack.From, "to", new_pack.To)
	}

	new_pack.Date, _ = date.HandShakeToJson(node.Connections, status) //Статус нужен для того, чтобы определять кто начал рукопожатие, иначе сеть будет постоянно их слать

	node.Send(&new_pack)

}

func (node *Node) SendMessageTo(To string, message string) {
	var new_pack = date.Packege{
		From:  node.Address.IPv4 + node.Address.Port,
		Name:  node.Name,
		Title: node.Titles[1],
		Type:  node.Types[0],
		Date:  []byte(message),
	}

	if strings.Contains(To, ":") {
		new_pack.To = To
	} else {
		for key, item := range node.Connections { //Поиск адреса по имени
			if To == item {
				new_pack.To = key
			}
		}
	}

	if new_pack.To == "" {
		fmt.Println("Не удалось найти получателя", To)
		return
	}

	node.Send(&new_pack)
}

func (node *Node) SendMessageToAll(message string) { //Отправка сообщений всем
	var new_pack = date.Packege{
		From:  node.Address.IPv4 + node.Address.Port,
		Name:  node.Name,
		Title: node.Titles[1],
		Type:  node.Types[0],
		Date:  []byte(message),
	}
	for addr := range node.Connections { //Переборам отправляем сообщение
		new_pack.To = addr
		node.Send(&new_pack)
	}
}

func (node *Node) Send(pack *date.Packege) { //Отправление данных конкретному пользователю
	conn, err := net.Dial("tcp", pack.To) //Подключаемся
	if err != nil {                       //Если подключение не прошло, забываем о узле
		delete(node.Connections, pack.To)
		return
	}
	defer conn.Close()

	byte_array, err := date.ToByteArray(*pack)
	conn.Write(byte_array) //Отправляем
}
