package main

import (
	"fmt"

	"github.com/MonResu/petcalendargo/calendar"
	"github.com/MonResu/petcalendargo/cmd"
	"github.com/MonResu/petcalendargo/logger"
	"github.com/MonResu/petcalendargo/storage"
)

func Add(a int, b int) int {
	return a + b
}

func main() {
	s := storage.NewJsonStorage("calendar.json")
	c := calendar.NewCalendar(s)
	err := c.Load()
	if err != nil {
		fmt.Println(err)
	}
	err = logger.LogInit("app.log")
	if err != nil {
		fmt.Println(err)
	}
	cli := cmd.NewCmd(c)
	cli.Run()
	
}
