package date

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
)

type Packege struct {
	To    string
	From  string
	Title string
	Name  string
	Type  string
	Date  []byte
}

type HandShake struct { //Информация о узлах при рукопожатии
	Nodes  []Node
	Status bool
}

type Node struct { //Узел адрес|Имя
	Address string
	Name    string
}

func HandShakeToJson(nodes map[string]string, status bool) ([]byte, error) {
	var handShake = HandShake{} //Создание списка адресов для рукопожатия
	for addr := range nodes {
		handShake.Nodes = append(handShake.Nodes, Node{Address: addr, Name: nodes[addr]})
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
