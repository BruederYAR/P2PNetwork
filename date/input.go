package date

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
)

func InputString() string { //Чтение с консоли
	msg, _ := bufio.NewReader(os.Stdin).ReadString('\n') //Читаем буфер
	return strings.Replace(msg, "\n", "", -1)            //Убираем переходы на следующую строку и возвращаем сообщение
}



type Input struct {
	Modules    map[string]ModuleInfo
	ModulePath string
	OS string
	OSseparator string
	Args       map[string]string
}

func NewInput() *Input {
	var input = Input{
		Modules: make(map[string]ModuleInfo),
		OS: runtime.GOOS,
		Args:    make(map[string]string),
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
						moduleInfo := OpenModule(value + input.OSseparator + f.Name())
						moduleInfo.Path = input.ModulePath + input.OSseparator + f.Name()

						input.Modules[moduleInfo.Name] = moduleInfo  //Добавляем в карту модули по названию и данных о них в json
					}
				}
			}else{
				for _,f := range files{
					if !strings.Contains(f.Name(), "."){
						moduleInfo := OpenModule(value + input.OSseparator + f.Name())
						name := moduleInfo.Path
						moduleInfo.Path = value + input.OSseparator + f.Name()

						input.Modules[name] = moduleInfo  //Добавляем в карту модули по названию и данных о них в json
					}
				}
			}
			
			break
		}
	}

	return &input

}

