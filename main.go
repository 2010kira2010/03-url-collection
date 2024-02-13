package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type Item struct {
	Name string
	Date time.Time
	Tags string
	Link string
}

// NewItem создает новый экземпляр структуры Item с заполненными полями
func NewItem(name, tags, link string) *Item {
	return &Item{
		Name: name,
		Date: time.Now(),
		Tags: tags,
		Link: link,
	}
}

func main() {
	defer func() {
		// Завершаем работу с клавиатурой при выходе из функции
		_ = keyboard.Close()
	}()

	fmt.Println("Программа для добавления url в список")
	fmt.Println("Для выхода и приложения нажмите Esc")

	items, err := loadItemsFromFile("urls.txt")
	if err != nil {
		log.Fatal("Ошибка чтения файла: ", err)
	}

OuterLoop:
	for {
		// Подключаем отслеживание нажатия клавиш
		if err := keyboard.Open(); err != nil {
			log.Fatal(err)
		}

		char, key, err := keyboard.GetKey()
		if err != nil {
			log.Fatal(err)
		}

		switch char {
		case 'a':
			if err := keyboard.Close(); err != nil {
				log.Fatal(err)
			}

			// Добавление нового url в список хранения
			fmt.Println("Введите новую запись в формате <url описание теги>")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			args := strings.Fields(text)
			if len(args) < 3 {
				fmt.Println("Введите правильный аргументы в формате url описание теги")
				continue OuterLoop
			}

			item := NewItem(args[1], args[2], args[0])
			items = append(items, item)

			err := saveItemsToFile(items, "urls.txt")
			if err != nil {
				log.Fatal("Ошибка сохранения в файл: ", err)
			}

		case 'l':
			// Вывод списка добавленных url. Выведите количество добавленных url и список с данными url
			// Вывод в формате
			// Имя: <Описание>
			// URL: <url>
			// Теги: <Теги>
			// Дата: <дата>

			for _, item := range items {
				fmt.Printf("Имя: %s\n", item.Name)
				fmt.Printf("URL: %s\n", item.Link)
				fmt.Printf("Теги: %s\n", item.Tags)
				fmt.Printf("Дата: %s\n", item.Date.Format(time.RFC3339))
				fmt.Println()
			}
			fmt.Printf("Количество добавленных url: %d\n", len(items))

		case 'r':
			if err := keyboard.Close(); err != nil {
				log.Fatal(err)
			}
			// Удаление url из списка хранения
			fmt.Println("Введите имя ссылки, которое нужно удалить")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			name := strings.TrimSpace(text)

			removed := false
			for i, item := range items {
				if item.Name == name {
					items = append(items[:i], items[i+1:]...)
					removed = true
					break
				}
			}

			if removed {
				fmt.Printf("Ссылка %s удалена\n", name)
				err := saveItemsToFile(items, "urls.txt")
				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Printf("Ссылка %s не найдена\n", name)
			}
		default:
			// Если нажата Esc выходим из приложения
			if key == keyboard.KeyEsc {
				return
			}
		}
	}
}
