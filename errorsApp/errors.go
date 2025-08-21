package errorsApp

import "errors"

var ErrEventIDNotFound = errors.New("задача с таким ID не найдена")
var ErrLoadCalendarEmpty = errors.New("загружать нечего")
var ErrReminderInPast = errors.New("напоминание в прошлое")
var ErrReminderStop = errors.New("таймер уже истек или был ранее остановлен")
