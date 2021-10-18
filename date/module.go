package date

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

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

func returnCmd(dir string) *exec.Cmd {
	switch runtime.GOOS {
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

func (input *Input) CommandExecute(com string, module string) string {
	com = strings.TrimSpace(com); module = strings.TrimSpace(module)
	answer, err := RequestModule(input.Modules[module].Path, com)
	
	if err != nil {
		fmt.Println("Не удалось найти команду или модуль " + input.Modules[module].Path)
		return ""
	}


	return string(answer)
}
