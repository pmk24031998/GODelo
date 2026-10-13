package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	folders := []string{"Входящие", "Работа", "Личное"}
	in := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("=== папки ===")

		for i, name := range folders {
			fmt.Printf("%d. %s\n", i+1, name)
		}

		fmt.Println("0. Выход")
		fmt.Print("> ")
		choice, err := readLine(in)

		if err != nil {
			return
		}

		switch choice {
		case "0":
			fmt.Println("Пока!")
			return
		default:
			index, err := folderIndex(folders, choice)
			if err != nil {
				fmt.Println(err)
				continue
			}
			var exit bool
			folders, exit = openFolder(in, folders, folders[index])
			if exit {
				return
			}
		}
	}
}
func openFolder(in *bufio.Reader, folders []string, current string) ([]string, bool) {
	for {
		fmt.Printf("\n[%s]\n", current)
		fmt.Println("задач пока нет — занятие 02")
		fmt.Println("1. Создать папку")
		fmt.Println("2. Удалить папку")
		fmt.Println("3. Переименовать папку")
		fmt.Println("4. Назад")
		fmt.Println("0. Выход")
		fmt.Print("> ")

		choice, err := readLine(in)
		if err != nil {
			return folders, true
		}

		switch choice {
		case "1":
			fmt.Print("имя новой папки: ")
			name, err := readLine(in)
			if err != nil {
				return folders, true
			}
			folders, err = createFolder(folders, name)
			if err != nil {
				fmt.Println("ошибка:", err)
			}
		case "2":
			fmt.Print("какую папку удалить: ")
			name, err := readLine(in)
			if err != nil {
				return folders, true
			}
			folders, err = deleteFolder(folders, name)
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}
			if name == current {
				return folders, false
			}
		case "3":
			fmt.Print("какую папку переименовать: ")
			oldName, err := readLine(in)
			if err != nil {
				return folders, true
			}
			fmt.Print("новое имя: ")
			newName, err := readLine(in)
			if err != nil {
				return folders, true
			}
			folders, err = renameFolder(folders, oldName, newName)
			if err != nil {
				fmt.Println("ошибка:", err)
				continue
			}
			if oldName == current {
				current = newName
			}
		case "4":
			return folders, false
		case "0":
			fmt.Println("Пока!")
			return folders, true
		default:
			fmt.Println("нет такого пункта")
		}
	}
}

func createFolder(folders []string, name string) ([]string, error) {
	if name == "" {
		return folders, errors.New("пустое имя")
	}
	return append(folders, name), nil
}

func deleteFolder(folders []string, name string) ([]string, error) {
	if name == "" {
		return folders, errors.New("пустое имя")
	}

	updated := []string{}
	found := false
	for _, folder := range folders {
		if !found && folder == name {
			found = true
			continue
		}
		updated = append(updated, folder)
	}
	if !found {
		return folders, errors.New("такой папки нет")
	}
	return updated, nil
}

func renameFolder(folders []string, oldName string, newName string) ([]string, error) {
	if oldName == "" || newName == "" {
		return folders, errors.New("пустое имя")
	}
	for i, folder := range folders {
		if folder == oldName {
			folders[i] = newName
			return folders, nil
		}
	}
	return folders, errors.New("такой папки нет")
}

func folderIndex(folders []string, choice string) (int, error) {
	number, err := strconv.Atoi(choice)
	if err != nil || number < 1 || number > len(folders) {
		return 0, errors.New("такой папки нет")
	}
	return number - 1, nil
}

func readLine(in *bufio.Reader) (string, error) {
	text, err := in.ReadString('\n')
	text = strings.TrimSpace(text)
	if text == "" && err != nil {
		return "", err
	}
	return text, nil
}
