package pkg

import (
	"time"

	"gorm.io/gorm"
)

type InputEmployee struct {
	Name                     string
	Surname                  string
	Login                    string
	Password                 string
	WeekWorkingTimeInMinutes uint
	WorkingDays              string
}

func (u *InputEmployee) ToEmployee() Employee {
	return Employee{
		User: &User{
			Name:     u.Name,
			Surname:  u.Surname,
			Login:    u.Login,
			Password: u.Password,
		},
		WeekWorkingTimeInMinutes: u.WeekWorkingTimeInMinutes,
		WorkingDays:              u.WorkingDays,
	}
}

type User struct {
	gorm.Model
	Name     string
	Surname  string
	Login    string `gorm:"unique"`
	Password string `json:"-"`
	Tokens   []Token
}

type InputToken struct {
	Name string
}

type Token struct {
	gorm.Model
	Name   string
	UserID uint
	Token  string
}

type Employee struct {
	*User                    `gorm:"embedded"`
	WeekWorkingTimeInMinutes uint
	WorkingDays              string
}

type Activity struct {
	gorm.Model
	Start       *time.Time
	End         *time.Time
	Description string
	UserID      uint
}

type InputActivity struct {
	Start       *time.Time
	End         *time.Time
	Description string
}

type Hollyday struct {
	gorm.Model
	Start       time.Time
	End         time.Time
	Description string
	UserID      uint
}

type InputHollyday struct {
	Start       time.Time
	End         time.Time
	Description string
}

type WorkDay struct {
	gorm.Model
	Day        time.Time `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Overtime   int64
	ActiveTime int64
	UserID     uint `gorm:"UNIQUE_INDEX:compositeindex;type:text;not null"`
}

type Overview struct {
	Date                         time.Time
	WeekNumber                   int
	OvertimeThisDayInMinutes     int64
	ActiveTimeThisDayInMinutes   int64
	OvertimeThisWeekInMinutes    int64
	ActiveTimeThisWeekInMinutes  int64
	OvertimeThisMonthInMinutes   int64
	ActiveTimeThisMonthInMinutes int64
	OvertimeThisYearInMinutes    int64
	ActiveTimeThisYearInMinutes  int64
	ActiveActivity               *Activity
}
