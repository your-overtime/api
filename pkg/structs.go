package pkg

import "time"

type User struct {
	Name    string
	Surname string
	Login   string
}

type Employee struct {
	*User
	WeekHours int
}

type Activity struct {
	Start       time.Time
	End         time.Time
	Description string
}

type Project struct {
	Name       string
	Customer   *string
	Activities []Activity
}
