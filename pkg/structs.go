package pkg

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

type InputEmployee struct {
	Name                     string
	Surname                  string
	Login                    string
	Password                 string
	WeekWorkingTimeInMinutes uint
	NumWorkingDays           uint
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
		NumWorkingDays:           u.NumWorkingDays,
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
	NumWorkingDays           uint
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

type HolidayType string

const (
	HolidayTypeFree         HolidayType = "free"
	HolidayTypeSick         HolidayType = "sick"
	HolidayTypeLegalHoliday HolidayType = "legal_holiday"
)

var ErrInvalidType = errors.New("invalid type")

func StrToHolidayType(str string) (HolidayType, error) {
	switch strings.ToLower(str) {
	case "free":
		return HolidayTypeFree, nil
	case "sick":
		return HolidayTypeSick, nil
	case "legal_holiday":
		return HolidayTypeLegalHoliday, nil
	}

	return HolidayTypeFree, ErrInvalidType
}

type Holiday struct {
	gorm.Model
	Start       time.Time
	End         time.Time
	Description string
	Type        HolidayType
	UserID      uint
}

type InputHoliday struct {
	Start       time.Time
	End         time.Time
	Description string
	HolidayType HolidayType
}

type WorkDay struct {
	gorm.Model
	Day        time.Time `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
	Overtime   int64
	ActiveTime int64
	IsHoliday  bool
	UserID     uint `gorm:"UNIQUE_INDEX:compositeindex;index;not null"`
}

type InputWorkDay struct {
	Day        time.Time
	Overtime   int64
	ActiveTime int64
	UserID     uint
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
