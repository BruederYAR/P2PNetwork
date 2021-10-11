package date

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type StringOptions struct { //ну почему в go нет стандартных значений для аргументов. боль
	P1 string
	P2 string
	P3 string
	P4 string
	P5 string
}

func InputString() string { //Чтение с консоли
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n') //Читаем буфер
	return strings.Replace(msg, "\n", "", -1)            //Убираем переходы на следующую строку и возвращаем сообщение
}

func RequestModule(dir string, message string) ([]byte, error) {
	cmd := returnCmd(dir)

	// Чтобы вводить что-то в стандартный поток ввода другой программы, нужно получить ее pipe.
	pipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	pipe.Write([]byte(message)) // Куда впоследствии можно что-то писать.
	pipe.Close()                // После ввода всех данных нужно обязательно его закрыть.

	output, err := cmd.Output() // Самый простой способ получить вывод другой программы, использовать:
	if err != nil {
		return nil, err
	}

	return output, nil
}

func returnCmd(dir string) *exec.Cmd{
	switch runtime.GOOS{
	case "windows":
		return exec.Command("cmd", "/C", dir)
	case "linux":
		return exec.Command("./" + dir)
	}
	return exec.Command("./" + dir)
}

func OpenModule(dir string) ModuleInfo {

	output, err := RequestModule(dir, "/cmd")

	if err != nil {
		fmt.Println(err)
		panic("Не удаётся открыть модуль " + dir)
	}

	var result ModuleInfo
	json.Unmarshal(output, &result)

	return result
}

type Input struct {
	Modules    map[string]ModuleInfo
	ModulePath string
	OS string
	OSseparator string
	Cmds       map[string]string
	Args       map[string]string
}

func NewInput() *Input {
	var input = Input{
		Modules: make(map[string]ModuleInfo),
		OS: runtime.GOOS,
		Args:    make(map[string]string),
		Cmds:    make(map[string]string),
	}

	if input.OS == "windows"{
		input.OSseparator = "\\"
	}else{
		input.OSseparator = "/"
	}

	otherArgs := os.Args[3:]
	for i := 0; i < len(otherArgs)-1; i++ { //Ищем ключи и значения для аргументов
		if otherArgs[i][0] == '-' && otherArgs[i+1][0] != '-' {
			input.Args[otherArgs[i]] = otherArgs[i+1]
		}
		if otherArgs[i][0] == '-' && otherArgs[i+1][0] == '-' {
			input.Args[otherArgs[i]] = ""
		}
	}

	//По аргументам выполняется функционал
	for keys, value := range input.Args {
		switch keys {
		case "-c": //Подключение разных модулей
			input.ModulePath = value
			files, err := ioutil.ReadDir(value)
			if err != nil {
				fmt.Println("Не удалось открыть деректорию с модулями")
				panic(err)
			}
			
			if len(files) == 0{
				fmt.Println("Модули в директории "+ input.ModulePath +" не обнаружены")
			} 

			if input.OS == "windows"{
				for _, f := range files { //Ищется exe
					if strings.Split(f.Name(), ".")[1] == "exe" {
						input.Modules[f.Name()] = OpenModule(value + input.OSseparator + f.Name()) //Добавляем в карту модули по названию и данных о них в json
	
						for i := 0; i < len(input.Modules[f.Name()].Cmds); i++ { //Добавляем команды молуоя в общий список команд
							input.Cmds[input.Modules[f.Name()].Cmds[i].Cmd] = f.Name()
						}
					}
				}
			}else{
				for _,f := range files{
					if !strings.Contains(f.Name(), "."){
						input.Modules[f.Name()] = OpenModule(value + input.OSseparator + f.Name()) //Добавляем в карту модули по названию и данных о них в json
	
						for i := 0; i < len(input.Modules[f.Name()].Cmds); i++ { //Добавляем команды молуоя в общий список команд
							input.Cmds[input.Modules[f.Name()].Cmds[i].Cmd] = f.Name()
						}
					}
				}
			}
			
			break
		}
	}

	return &input

}

func (input *Input) CommandExecute(com string) string {
	com = strings.TrimSpace(com)
	answer, err := RequestModule(input.ModulePath + input.OSseparator + input.Cmds[com], com)
	
	if err != nil {
		fmt.Println("Не удалось найти команду или модуль " + input.ModulePath + input.OSseparator + input.Cmds[com])
		return ""
	}

	return string(answer)
}
