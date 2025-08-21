package reminder

import (
	"fmt"
	"time"

	"github.com/monresu/app/errorsApp"
)

type Reminder struct {
	Message  string
	At       time.Time
	Sent     bool
	timer    *time.Timer
	notifyer func(string)
}

func NewReminder(message string, at time.Time, notify func(string)) (*Reminder, error) {
	now := time.Now().In(at.Location())
	if at.Before(now) {
		return &Reminder{}, errorsApp.ErrReminderInPast
	}

	return &Reminder{
		Message:  message,
		At:       at,
		Sent:     false,
		notifyer: notify,
	}, nil
}

func (r *Reminder) Start() {
	delay := time.Until(r.At)
	r.timer = time.AfterFunc(delay, r.Send)
}

func (r *Reminder) Send() {
	if r.Sent {
		return
	}
	r.notifyer(r.Message)
	r.Sent = true
}

func (r *Reminder) Stop() error {
	isStop := r.timer.Stop()
	if !isStop {
		return fmt.Errorf("Reminder Stop: %w", errorsApp.ErrReminderStop)
	}
	return nil
}
