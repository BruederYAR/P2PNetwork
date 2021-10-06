package node

import (
	"Network/date"
	"crypto/rand"
	"crypto/rsa"
	"net"
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
	IP string
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
	newnode := &Node{
		Titles:      map[int]string{0: "handshake", 1: "date"},
		Types:       map[int]string{0: "string", 1: "json"},
		Connections: make(map[string]*date.NodeInfo),
		Input:       input,
		Name:        os.Args[2],
		PrivateKey:  *PrivateKey,
		PublicKey:   PrivateKey.PublicKey,
	}
	ipv4, ipv6 := LocalIpAddress(*input)
	switch splited[0] {
	case "":
		newnode.Address = Address{IP: ipv4, Port: ":" + splited[1]}
	case "ipv4":
		newnode.Address = Address{IP: ipv4, Port: ":" + splited[1]}
	case "ipv6":
		newnode.Address = Address{IP: ipv6, Port: ":" + splited[1]}
	case "null":
		newnode.Address = Address{IP: "", Port: ":" + splited[1]}
	default:
		newnode.Address = Address{IP: splited[0], Port: ":" + splited[1]}
	}

	return newnode
}

func LocalIpAddress(input date.Input) (ipv4 string, ipv6 string){

	// Получаем все доступные сетевые интерфейсы
	interfaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
 
	for _, interf := range interfaces {
		// Список адресов для каждого сетевого интерфейса
		addrs, err := interf.Addrs()
		if err != nil {
			panic(err)
		}
		
		if input.OS == "windows"{
			if !strings.Contains((strings.Split(addrs[1].String(), "/"))[0], "192.168"){
				continue
			}else{
				return (strings.Split(addrs[1].String(), "/"))[0], (strings.Split(addrs[0].String(), "/"))[0]
			}
		}
		if input.OS == "linux"{
			if !strings.Contains((strings.Split(addrs[0].String(), "/"))[0], "192.168"){
				continue
			}else{
				return (strings.Split(addrs[0].String(), "/"))[0], (strings.Split(addrs[1].String(), "/"))[0]
			}
		}
	}
	return 
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


