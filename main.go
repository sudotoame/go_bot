package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	// Чтение токена из файла
	token, err := os.ReadFile("token.txt")
	if err != nil {
		log.Fatalf("Ошибка чтения token.txt: %v", err)
	}
	tokenStr := strings.TrimSpace(string(token))

	// Создание бота
	bot, err := tgbotapi.NewBotAPI(tokenStr)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Авторизован как %s", bot.Self.UserName)

	// Словарь для хранения TODO-листов пользователей
	todoLists := make(map[int64][]string)

	// Настройка получения обновлений
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Обработка обновлений
	for update := range updates {
		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		text := update.Message.Text
		log.Printf("[%d] %s", chatID, text)

		// Обработка команд
		if strings.HasPrefix(text, "/") {
			command := strings.TrimPrefix(text, "/")
			parts := strings.SplitN(command, " ", 2)
			cmd := parts[0]
			arg := ""
			if len(parts) > 1 {
				arg = parts[1]
			}

			switch cmd {
			case "add":
				if arg == "" {
					bot.Send(tgbotapi.NewMessage(chatID, "Используйте: /add <задача>"))
				} else {
					todoLists[chatID] = append(todoLists[chatID], arg)
					bot.Send(tgbotapi.NewMessage(chatID, "Задача добавлена!"))
				}
			case "list":
				tasks := todoLists[chatID]
				if len(tasks) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Ваш TODO-лист пуст."))
				} else {
					var msg string
					for i, task := range tasks {
						msg += fmt.Sprintf("%d. %s\n", i+1, task)
					}
					bot.Send(tgbotapi.NewMessage(chatID, msg))
				}
			case "done":
				tasks := todoLists[chatID]
				if len(tasks) == 0 {
					bot.Send(tgbotapi.NewMessage(chatID, "Список задач пуст."))
					break
				}

				if arg == "" {
					bot.Send(tgbotapi.NewMessage(chatID, "Используйте: /done <номер>"))
					break
				}

				index, err := strconv.Atoi(arg)
				if err != nil || index < 1 || index > len(tasks) {
					bot.Send(tgbotapi.NewMessage(chatID, "Неверный номер задачи."))
				} else {
					todoLists[chatID] = append(tasks[:index-1], tasks[index:]...)
					bot.Send(tgbotapi.NewMessage(chatID, "Задача выполнена!"))
				}
			case "help":
				helpText := "Доступные команды:\n" +
					"/add <задача> - добавить задачу\n" +
					"/list - показать список задач\n" +
					"/done <номер> - удалить задачу по номеру\n" +
					"/help - показать это сообщение"
				bot.Send(tgbotapi.NewMessage(chatID, helpText))
			default:
				bot.Send(tgbotapi.NewMessage(chatID, "Неизвестная команда. Используйте /help для справки."))
			}
		} else {
			bot.Send(tgbotapi.NewMessage(chatID, "Пожалуйста, используйте команды. Например, /help"))
		}
	}
}
