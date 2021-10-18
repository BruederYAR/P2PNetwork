package node

import (
	"Network/date"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
)

func (node *Node) recipientSearch(To string) (string, error){
	if strings.Contains(To, ":") { //Если адрес
		if node.Connections[To] == nil { //Если адреса нет в списке адресов, выполняем рукопожатия
			node.HandShake(To, true)
			return "", errors.New("[err] Неизвесный адрес, выполненно рукопожатие попробуйте ещё раз")
		}
		return To, nil
	} else { //Если имя
		for key, item := range node.Connections { //Поиск адреса по имени
			if item != nil{
				if To == item.Name {
					return key, nil
				}
			} 
		}
		if To == "" {
			return "", errors.New("Не удалось найти получателя " + To)
		}
	}
	return "", errors.New("Не удалось найти получателя " + To)
}


func (node *Node) HandShake(address string, status bool) { //Рукопожатие при первом подключении
	var new_pack = date.Packege{
		From:      node.Address.IP + node.Address.Port,
		To:        address,
		Name:      node.Name,
		PublicKey: node.PublicKey,
		Type:      node.Types[1],
		Date:      []byte{},
		Title:     node.Titles[0],
	}

	if !status {
		fmt.Println("HandShake from", new_pack.From, "to", new_pack.To)
	}

	new_pack.Date, _ = date.HandShakeToJson(node.Connections, status) //Статус нужен для того, чтобы определять кто начал рукопожатие, иначе сеть будет постоянно их слать

	node.Send(&new_pack)

}

func (node *Node) CommandRequestTo(To string, cmd []byte){
	var new_pack = date.Packege{
		From:      node.Address.IP + node.Address.Port,
		Name:      node.Name,
		PublicKey: node.PublicKey,
		Title:     node.Titles[3],
		Type:      node.Types[0],
	}

	To, err := node.recipientSearch(To)
	if err != nil{
		fmt.Println(err.Error())
		return
	}
	new_pack.To = To
	new_pack.Date = date.RSA_OAEP_Encrypt(cmd, node.Connections[new_pack.To].PublicKey)

	node.Send(&new_pack)
}

func (node *Node) ModuleRequestTo(To string, module string, cmd string) {
	var new_pack = date.Packege{
		From:      node.Address.IP + node.Address.Port,
		Name:      node.Name,
		PublicKey: node.PublicKey,
		Title:     node.Titles[2],
		Type:      node.Types[1],
	}

	To, err := node.recipientSearch(To)

	if err != nil{
		fmt.Println(err.Error())
		return
	}
	new_pack.To = To

	data := date.CmdRequest{Cmd: cmd, Module: module}
	message,err := json.Marshal(data)
	new_pack.Date = date.RSA_OAEP_Encrypt(message, node.Connections[new_pack.To].PublicKey)

	node.Send(&new_pack)
}

func (node *Node) SendMessageTo(To string, message []byte) {
	var new_pack = date.Packege{
		From:      node.Address.IP + node.Address.Port,
		Name:      node.Name,
		PublicKey: node.PublicKey,
		Title:     node.Titles[1],
		Type:      node.Types[0],
	}

	if strings.Contains(To, ":") { //Если адрес
		if node.Connections[To] == nil { //Если адреса нет в списке адресов, выполняем рукопожатия
			fmt.Println("[err] Неизвесный адрес, выполненно рукопожатие попробуйте ещё раз")
			node.HandShake(To, true)
			return
		}
		new_pack.To = To
	} else { //Если имя
		for key, item := range node.Connections { //Поиск адреса по имени
			if item != nil{
				if To == item.Name {
					new_pack.To = key
				}
			} 
		}
		if new_pack.To == "" {
			fmt.Println("Не удалось найти получателя", To)
			return
		}
	}

	new_pack.Date = date.RSA_OAEP_Encrypt(message, node.Connections[new_pack.To].PublicKey)

	node.Send(&new_pack)
}

func (node *Node) SendMessageToAll(message []byte) { //Отправка сообщений всем
	var new_pack = date.Packege{
		From:      node.Address.IP + node.Address.Port,
		Name:      node.Name,
		PublicKey: node.PublicKey,
		Title:     node.Titles[1],
		Type:      node.Types[0],
	}
	for addr := range node.Connections { //Переборам отправляем сообщение
		new_pack.To = addr
		new_pack.Date = date.RSA_OAEP_Encrypt(message, node.Connections[new_pack.To].PublicKey)
		node.Send(&new_pack)
	}
}

func (node *Node) Send(pack *date.Packege) { //Отправление данных конкретному пользователю
	conn, err := net.Dial("tcp", pack.To) //Подключаемся
	if err != nil {                       //Если подключение не прошло, забываем о узле
		delete(node.Connections, pack.To)
		fmt.Println("Ошибка подключения к", pack.To)
		return
	}
	defer conn.Close()

	byte_array, _ := date.ToByteArray(*pack)
	conn.Write(byte_array) //Отправляем
}