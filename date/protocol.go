package date

import (
	"bytes"
	"crypto/rsa"
	"encoding/gob"
	"encoding/json"
)

type Packege struct {
	To    string
	From  string
	Title string
	Name  string
	PublicKey rsa.PublicKey
	Type  string
	Date  []byte
}

type HandShake struct { //Информация о узлах при рукопожатии
	Nodes  []Node
	Status bool
}

type Node struct { //Узел адрес|Имя
	Address string
	PublicKey rsa.PublicKey
	Name    string
}

type NodeInfo struct {
	Name      string
	PublicKey rsa.PublicKey
}

func HandShakeToJson(nodes map[string]*NodeInfo, status bool) ([]byte, error) {
	var handShake = HandShake{} //Создание списка адресов для рукопожатия
	for addr := range nodes {
		handShake.Nodes = append(handShake.Nodes, Node{Address: addr, Name: nodes[addr].Name, PublicKey: nodes[addr].PublicKey})
	}
	handShake.Status = status

	return json.Marshal(handShake)
}

func ToByteArray(pack Packege) ([]byte, error) {
	var buffer bytes.Buffer
	enc := gob.NewEncoder(&buffer)

	err := enc.Encode(pack)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ToPackege(message []byte) (Packege, error) {
	var buffer bytes.Buffer
	var pack Packege

	buffer.Write(message)
	dec := gob.NewDecoder(&buffer)

	err := dec.Decode(&pack)
	if err != nil {
		return Packege{}, err
	}

	return pack, nil
}
