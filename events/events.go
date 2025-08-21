package events

import (
	"errors"
	"fmt"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"github.com/monresu/app/reminder"
)

func getNextID() string {
	return uuid.New().String()
}

type Event struct {
	ID       string             `json:"id"`
	Title    string             `json:"title"`
	StartAt  time.Time          `json:"start_at"`
	Priority Priority           `json:"priority"`
	Reminder *reminder.Reminder `json:"reminder"`
}

func NewEvent(title string, dateStr string, priority Priority) (*Event, error) {
	if t, err := time.ParseInLocation("02.01.2006 15:04", dateStr, time.Local); err == nil {
		return createValidatedEvent(title, t, priority)
	}

	t, err := dateparse.ParseIn(dateStr, time.Local)
	if err != nil {
		return nil, errors.New("неверный формат даты. Используйте ДД.ММ.ГГГГ ЧЧ:ММ (11.08.2025 13:35) или ГГГГ-ММ-ДД ЧЧ:ММ")
	}

	return createValidatedEvent(title, t, priority)
}

func createValidatedEvent(title string, t time.Time, priority Priority) (*Event, error) {
	if !isValidTitle(title) {
		return nil, errors.New("заголовок должен содержать 3-50 символов (буквы, цифры, пробелы)")
	}
	if err := priority.Validate(); err != nil {
		return nil, errors.New("приоритет должен быть low, medium или high")
	}

	return &Event{
		ID:       getNextID(),
		Title:    title,
		StartAt:  t.In(time.Local),
		Priority: priority,
	}, nil
}

func (e *Event) AddReminder(message string, at time.Time, notify func(string)) error {
	r, err := reminder.NewReminder(message, at, notify)
	if err != nil {
		return err
	}
	e.Reminder = r
	r.Start()
	return nil
}

func (e *Event) RemoveReminder() {
	e.Reminder = nil
}

func (e Event) Print() {
	fmt.Println(e.Title, e.StartAt.Format("02.01.2006 15:04"), e.ID)
}

func (e *Event) UpdateEvent(title string, date string, priority Priority) error {
	if !isValidTitle(title) {
		return errors.New("неверный формат заголовка")
	}
	if priority.Validate() != nil {
		return errors.New("неверный формат приоритета")
	}
	time, err := dateparse.ParseAny(date)
	if err != nil {
		return errors.New("неверный формат даты")
	}
	e.Title = title
	e.StartAt = time
	return nil
}
