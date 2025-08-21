package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/google/shlex"
	"github.com/MonResu/petcalendargo/calendar"
	"github.com/MonResu/petcalendargo/events"
	"github.com/MonResu/petcalendargo/logger"
)

func (c *Cmd) executor(input string) {
	input = strings.TrimSpace(input)
	if input == "" {
		return
	}
	logger.PrintInfo(input)
	parts, err := shlex.Split(input)
	if err != nil {
		return
	}

	cmd := strings.ToLower(parts[0])

	switch cmd {
	case "add":
		if len(parts) < 4 {
			text := "Формат: add \"название события\" \"дата и время\" \"приоритет\""
			fmt.Println(text)
			logger.PrintError(text)
			return
		}

		title := parts[1]
		date := parts[2]
		priority := events.Priority(parts[3])

		e, err := c.calendar.AddEvent(title, date, priority)
		if err != nil {
			text := "Ошибка добавления: " + err.Error()
			fmt.Println(text)
			logger.PrintError(text)

		} else {
			text := "Событие: " + e.Title + " добавлено"
			fmt.Println(text)
			logger.PrintInfo(text)
		}
	case "list":
		events := c.calendar.GetEvents()
		for ID := range events {
			events[ID].Print()
		}

	case "remove":
		if len(parts) < 2 {
			text := "Формат: remove \"название события\""
			fmt.Println(text)
			logger.PrintError(text)
			return
		}
		id := parts[1]
		found := false
		for key := range c.calendar.CalendarEvents {
			if c.calendar.CalendarEvents[key].ID == id {
				c.calendar.DeleteEvent(c.calendar.CalendarEvents[key].ID)
				found = true
			}
		}
		if !found {
			text := "Задача не найдена"
			fmt.Println(text)
			logger.PrintError(text)
			return
		}

	case "update":
		if len(parts) < 5 {
			text := "Формат: update \"названия события\" \"новое название\" \"новая дата и время\" \"новый приоритет\""
			fmt.Println(text)
			logger.PrintError(text)
			return
		}
		id := parts[1]
		newTitle := parts[2]
		newDate := parts[3]
		newPriority := parts[4]
		found := false
		for key := range c.calendar.CalendarEvents {
			if c.calendar.CalendarEvents[key].ID == id {
				c.calendar.CalendarEvents[key].UpdateEvent(newTitle, newDate, events.Priority(newPriority))
			}
		}
		if !found {
			text := "Задача не найдена"
			fmt.Println(text)
			logger.PrintError(text)
			return
		}

	case "help":
		var rules = "Список всех команд:"
		for _, val := range commands {
			rules += "\n"
			s := val.name_command + ". " + val.description
			rules += s
		}
		fmt.Println(rules)
		logger.PrintInfo(rules)

	case "exit":
		if err := c.calendar.Save(); err != nil {
			text := fmt.Sprintf("Ошибка сохранения календаря: %v", err)
			fmt.Println(text)
			logger.PrintError(text)
		}
		logger.Close()
		close(c.calendar.Notification)
		os.Exit(0)

	case "add_reminder":
		if len(parts) < 4 {
			text := commands[3].description
			fmt.Println(text)
			logger.PrintError(text)
			return
		}
		id := parts[1]
		text_mess := parts[2]
		q := parts[3]
		time_mess, err := time.ParseInLocation("02.01.2006 15:04", q, time.Now().Location())
		if err != nil {
			text := "Ошибка в формате времени"

			fmt.Println(text)
			logger.PrintError(text)

			fmt.Println(err.Error())
			logger.PrintError(err.Error())

			return
		}
		err = c.calendar.SetEventReminder(id, text_mess, time_mess)
		if err != nil {
			fmt.Println(err.Error())
			logger.PrintError(err.Error())
			return
		} else {
			text := "Успешно"
			fmt.Println(text)
			logger.PrintError(text)

		}

	case "remove_reminder":
		id := parts[1]
		err := c.calendar.RemoveEventReminder(id)
		if err != nil {
			text := err.Error()
			fmt.Println(text)
			logger.PrintError(text)
			return
		}

	default:
		text := `Неизвестная команда:
		"Введите 'help' для списка команд`
		fmt.Println(text)
		logger.PrintError(text)
	}
}

type command struct {
	name_command string
	description  string
}

var commands = []command{
	{"add", "Добавить событие, формат: add \"название события\" \"дата и время\" \"приоритет\""},
	{"remove", "Удалить событие, формат: delete \"id события\""},
	{"edit", "Отредактировать событие, формат: change \"id события\" \"название события\" \"дата и время\" \"приоритет\""},
	{"add_reminder", "Добавить напоминание, формат: add \"id события\" \"текст напоминания\" \"дата и время напоминания\""},
	{"remove_reminder", "Удалить напоминание, формат: remove_reminder \"id напоминания\""},
	{"list", "Показать все события"},
	{"help", "Показать справку"},
	{"exit", "Выйти из программы"},
}

func (c *Cmd) completer(d prompt.Document) []prompt.Suggest {
	var suggestions = []prompt.Suggest{}
	for _, val := range commands {
		suggestions = append(suggestions, prompt.Suggest{Text: val.name_command, Description: val.description})
	}

	return prompt.FilterHasPrefix(suggestions, d.GetWordBeforeCursor(), true)
}

func (c *Cmd) Run() {
	prompt.OptionMaxSuggestion(16)
	p := prompt.New(
		c.executor,
		c.completer,
		prompt.OptionPrefix("> "),
	)
	go func() {
		for msg := range c.calendar.Notification {
			fmt.Println("🔔 Уведомление:", msg)
		}
	}()
	p.Run()
}

type Cmd struct {
	calendar *calendar.Calendar
}

func NewCmd(c *calendar.Calendar) *Cmd {
	cmd := &Cmd{
		calendar: c,
	}
	return cmd
}
