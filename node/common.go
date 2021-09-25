package node

import (
	"Network/date"
	"crypto/rand"
	"crypto/rsa"
	"os"
	"strings"
)


type Node struct {
	Titles      map[int]string
	Types       map[int]string
	Connections map[string]*date.NodeInfo
	Input       *date.Input
	Address     Address
	Name        string         //Имя узла
	PrivateKey  rsa.PrivateKey //Приватный ключ для rsa
	PublicKey   rsa.PublicKey  //Публичный ключ для rsa
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
	PrivateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return &Node{
		Titles:      map[int]string{0: "handshake", 1: "date"},
		Types:       map[int]string{0: "string", 1: "json"},
		Connections: make(map[string]*date.NodeInfo),
		Input:       input,
		Address:     Address{IPv4: splited[0], Port: ":" + splited[1]},
		Name:        os.Args[2],
		PrivateKey:  *PrivateKey,
		PublicKey:   PrivateKey.PublicKey,
	}
}

func (node *Node) Run(handleServer func(*Node), handleClient func(*Node)) { //Выполняется запуск как сервер, так и клиента
	go handleServer(node)
	handleClient(node)
}

func (node *Node) ConnectTo(addresses []string, name string, publickey rsa.PublicKey) { //Добавление в список подключений
	for _, addr := range addresses {
		if addr == "" {
			panic("Пустой адрес, всё плохо")
		}
		node.Connections[addr] = &date.NodeInfo{ Name: name, PublicKey: publickey}

	}
}


