package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Title  string
	Closed bool
}

func (t *Task) Close() error {
	if t.Closed {
		return errors.New("задачу уже закрыли")
	}
	t.Closed = true
	return nil
}

func (t Task) status() string {
	if t.Closed {
		return "закрыта"
	}
	return "открыта"
}

type Folder struct {
	Name  string
	Tasks []Task
}

func (f *Folder) AddTask(title string) error {
	if title == "" {
		return errors.New("пустой заголовок")
	}
	f.Tasks = append(f.Tasks, Task{Title: title})
	return nil
}

func (f *Folder) CloseTask(number int) error {
	if number < 1 || number > len(f.Tasks) {
		return errors.New("такой задачи нет")
	}
	return f.Tasks[number-1].Close()
}

func main() {
	folders := []Folder{
		{Name: "Входящие"},
		{Name: "Работа"},
		{Name: "Личное"},
	}
	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("=== папки ===")
		for i, folder := range folders {
			fmt.Printf("%d. %s\n", i+1, folder.Name)
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
			folders, exit = openFolder(in, folders, folders[index].Name)
			if exit {
				return
			}
		}
	}
}

func openFolder(in *bufio.Reader, folders []Folder, current string) ([]Folder, bool) {
	for {
		fmt.Printf("\n[%s]\n", current)
		printTasks(folders, current)
		fmt.Println("1. Создать задачу")
		fmt.Println("2. Закрыть задачу")
		fmt.Println("3. Создать папку")
		fmt.Println("4. Удалить папку")
		fmt.Println("5. Переименовать папку")
		fmt.Println("6. Назад")
		fmt.Println("0. Выход")
		fmt.Print("> ")

		choice, err := readLine(in)
		if err != nil {
			return folders, true
		}

		switch choice {
		case "1":
			fmt.Print("заголовок задачи: ")
			title, err := readLine(in)
			if err != nil {
				return folders, true
			}
			folders, err = addTask(folders, current, title)
			if err != nil {
				fmt.Println("ошибка:", err)
			}
		case "2":
			fmt.Print("номер задачи: ")
			raw, err := readLine(in)
			if err != nil {
				return folders, true
			}
			number, err := strconv.Atoi(raw)
			if err != nil {
				fmt.Println("ошибка: такой задачи нет")
				continue
			}
			folders, err = closeTask(folders, current, number)
			if err != nil {
				fmt.Println("ошибка:", err)
			}
		case "3":
			fmt.Print("имя новой папки: ")
			name, err := readLine(in)
			if err != nil {
				return folders, true
			}
			folders, err = createFolder(folders, name)
			if err != nil {
				fmt.Println("ошибка:", err)
			}
		case "4":
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
		case "5":
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
		case "6":
			return folders, false
		case "0":
			fmt.Println("Пока!")
			return folders, true
		default:
			fmt.Println("нет такого пункта")
		}
	}
}

func printTasks(folders []Folder, current string) {
	index, err := findFolder(folders, current)
	if err != nil {
		fmt.Println("ошибка:", err)
		return
	}
	if len(folders[index].Tasks) == 0 {
		fmt.Println("задач нет")
		return
	}
	for i, task := range folders[index].Tasks {
		fmt.Printf("%d. %s — %s\n", i+1, task.Title, task.status())
	}
}

func addTask(folders []Folder, folderName string, title string) ([]Folder, error) {
	index, err := findFolder(folders, folderName)
	if err != nil {
		return folders, err
	}
	err = folders[index].AddTask(title)
	return folders, err
}

func closeTask(folders []Folder, folderName string, number int) ([]Folder, error) {
	index, err := findFolder(folders, folderName)
	if err != nil {
		return folders, err
	}
	err = folders[index].CloseTask(number)
	return folders, err
}

func createFolder(folders []Folder, name string) ([]Folder, error) {
	if name == "" {
		return folders, errors.New("пустое имя")
	}
	return append(folders, Folder{Name: name}), nil
}

func deleteFolder(folders []Folder, name string) ([]Folder, error) {
	if name == "" {
		return folders, errors.New("пустое имя")
	}

	updated := []Folder{}
	found := false
	for _, folder := range folders {
		if !found && folder.Name == name {
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

func renameFolder(folders []Folder, oldName string, newName string) ([]Folder, error) {
	if oldName == "" || newName == "" {
		return folders, errors.New("пустое имя")
	}
	for i, folder := range folders {
		if folder.Name == oldName {
			folders[i].Name = newName
			return folders, nil
		}
	}
	return folders, errors.New("такой папки нет")
}

func findFolder(folders []Folder, name string) (int, error) {
	for i, folder := range folders {
		if folder.Name == name {
			return i, nil
		}
	}
	return 0, errors.New("такой папки нет")
}

func folderIndex(folders []Folder, choice string) (int, error) {
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
