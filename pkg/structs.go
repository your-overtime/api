package pkg

import (
	"time"

	"gorm.io/gorm"
)

type InputEmployee struct {
	Name            string        `json:"name"`
	Surname         string        `json:"surname"`
	Login           string        `json:"login"`
	Password        string        `json:"password"`
	WeekWorkingTime time.Duration `json:"week_working_time"`
}

func (u *InputEmployee) toEmployee() Employee {
	return Employee{
		User: &User{
			Name:     u.Name,
			Surname:  u.Surname,
			Login:    u.Login,
			Password: u.Password,
		},
		WeekWorkingTime: u.WeekWorkingTime,
	}
}

type User struct {
	gorm.Model
	Name     string  `json:"name"`
	Surname  string  `json:"surname"`
	Login    string  `json:"login"`
	Password string  `json:"-"`
	Tokens   []Token `json:"tokens"`
}

type Token struct {
	gorm.Model
	CreationTime time.Time `json:"creation_time"`
	Name         string    `json:"name"`
	UserID       int       `json:"user_id"`
}

type Employee struct {
	*User           `gorm:"embedded"`
	WeekWorkingTime time.Duration `json:"week_working_time"`
}

type Activity struct {
	gorm.Model
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end"`
	Description string     `json:"description"`
	UserID      int        `json:"user_id"`
}

type InputActivity struct {
	Start       *time.Time `json:"start"`
	End         *time.Time `json:"end,omitempty"`
	Description string     `json:"description"`
}

type Hollyday struct {
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"`
}

type Overview struct {
	Date                     time.Time `json:"date"`
	WeekNumber               int       `json:"week_number"`
	OvertimeInMinutes        int       `json:"overtime_in_minutes"`
	WeekWorkingTimeInMinutes int       `json:"week_working_time_in_minutes"`
	ActiveActivity           Activity  `json:"active_activity"`
}
