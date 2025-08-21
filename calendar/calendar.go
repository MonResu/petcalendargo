package calendar

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/monresu/app/errorsApp"
	"github.com/monresu/app/events"
	"github.com/monresu/app/storage"
)

type Calendar struct {
	CalendarEvents map[string]*events.Event `json:"CalendarEvents"`
	Storage        storage.Store            `json:"-"`
	Notification   chan string              `json:"-"`
}

func NewCalendar(s storage.Store) *Calendar {
	return &Calendar{
		CalendarEvents: make(map[string]*events.Event),
		Storage:        s,
		Notification:   make(chan string),
	}
}

func (c *Calendar) Save() error {
	data, err := json.Marshal(c)
	if err != nil {
		return fmt.Errorf("Calendar Save: %w", err)
	}
	err = c.Storage.Save(data)
	if err != nil {
		return fmt.Errorf("Calendar Save: %w", err)
	}
	return nil
}

func (c *Calendar) Notify(msg string) {
	c.Notification <- msg
}

func (c Calendar) SetEventReminder(ID string, messageText string, at time.Time) error {
	isExists := c.CheckExists(ID)
	if isExists != nil {
		return fmt.Errorf("SetEventReminder: %w", isExists)
	}
	return c.CalendarEvents[ID].AddReminder(messageText, at, c.Notify)
}

func (c Calendar) CancelEventReminder(ID string) error {
	isExists := c.CheckExists(ID)
	if isExists != nil {
		return fmt.Errorf("CancelEventReminder: %w", isExists)
	}
	return c.CalendarEvents[ID].Reminder.Stop()
}

func (c Calendar) RemoveEventReminder(ID string) error {
	isExists := c.CheckExists(ID)
	if isExists != nil {
		return fmt.Errorf("RemoveEventReminder: %w", isExists)
	}
	c.CalendarEvents[ID].RemoveReminder()
	return nil
}

func (c *Calendar) Load() error {
	data, err := c.Storage.Load()
	if err != nil {
		return fmt.Errorf("Calendar Load: %w", err)
	}
	if len(data) == 0 {
		c.CalendarEvents = make(map[string]*events.Event)
		return fmt.Errorf("Calendar Load: %w", errorsApp.ErrLoadCalendarEmpty)
	}
	temp := struct {
		CalendarEvents map[string]*events.Event `json:"CalendarEvents"`
	}{}
	err = json.Unmarshal(data, &temp)
	if err != nil {
		return fmt.Errorf("Calendar Load: %w", err)
	}

	c.CalendarEvents = temp.CalendarEvents
	return nil
}

func (c *Calendar) AddEvent(title string, date string, priority events.Priority) (*events.Event, error) {
	e, err := events.NewEvent(title, date, priority)
	if err != nil {
		return &events.Event{}, fmt.Errorf("AddEvent: %w", err)
	}
	c.CalendarEvents[e.ID] = e

	return e, nil
}

func (c Calendar) GetEvents() map[string]*events.Event {
	return c.CalendarEvents
}

func (c *Calendar) DeleteEvent(ID string) error {
	isExists := c.CheckExists(ID)
	if isExists != nil {
		return fmt.Errorf("DeleteEvent: %w", isExists)
	}
	delete(c.CalendarEvents, ID)
	return nil

}

func (c *Calendar) EditEvent(ID string, title string, date string) error {
	e := c.CalendarEvents[ID]
	isExists := c.CheckExists(ID)
	if isExists != nil {
		return fmt.Errorf("EditEvent: %w", isExists)
	}

	err := e.UpdateEvent(title, date, events.PriorityHigh)
	return err
}

func (c Calendar) CheckExists(ID string) error {
	found := false
	for key := range c.CalendarEvents {
		if key == ID {
			found = true
		}
	}
	if !found {
		return errorsApp.ErrEventIDNotFound
	} else {
		return nil
	}
}
