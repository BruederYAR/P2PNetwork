package node

import (
	"Network/date"
	"encoding/json"
	"fmt"
	"net"
)

func handleServer(node *Node) { //Запуск сервера
	listen, err := net.Listen("tcp", "0.0.0.0"+node.Address.Port) //Слушаем определенный порт
	if err != nil {                                               //если есть ошибки вызываем панику
		panic("listen err")
	}
	defer listen.Close() //Ошибок нет - закрываем прослушку
	for {
		conn, err := listen.Accept() //Принимаем подключение
		if err != nil {              //При ошибке выходим из цикла и заного начинаем слушать
			break
		}
		go handleConnection(node, conn) //читаем данные
	}
}

func handleConnection(node *Node, conn net.Conn) { //Читаем данные
	defer conn.Close()
	var (
		buffer  = make([]byte, 512)
		message []byte
		pack    date.Packege
	)
	for {
		lenght, err := conn.Read(buffer) //Читаем всё в буфер
		if err != nil {
			break
		}

		message = append(message[:], buffer[:lenght]...) //Записываем только до длины, чтобы убрать мусор
	}

	pack, err := date.ToPackege(message) //Переводим в пакет
	if err != nil {                      //При ошибке метод закрываеться
		return
	}

	if pack.Name == node.Name { //Если вызвали сами себя, то выключаем метод
		return
	}

	node.ConnectTo([]string{pack.From}, pack.Name) //записываем того, кто отослал данные

	switch pack.Title {
	case node.Titles[1]: //date
		switch pack.Type {
		case node.Types[0]:
			fmt.Println(string(pack.Date)) //Выводим данные

			if string(pack.Date)[0] == '/' {
				message := node.Input.CommandExecute(string(pack.Date))
				node.SendMessageTo(pack.From, message)
			}
		}

	case node.Titles[0]: //Рукопожатие handshake
		var handShake date.HandShake
		json.Unmarshal(pack.Date, &handShake) //Забираем список узлов

		if handShake.Status { //Если начало рукопожатия
			node.HandShake(pack.From, false) //Отправляем узлы обратно
		}

		for _, local_node := range handShake.Nodes { //Добавляем узлы в локальные список
			if local_node.Name != node.Connections[local_node.Address] && local_node.Name != node.Name { //Если узел, который нам прислали был не известен, то выполняем рукопожатие с ним
				node.HandShake(local_node.Address, true)
			}
			node.Connections[local_node.Address] = local_node.Name
		}
	}

}