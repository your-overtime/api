package data

import (
	"time"

	"github.com/your-overtime/api/pkg"
	"gorm.io/gorm"
)

type ModelExtensions struct {
	CreatedAt time.Time
	UpdatedAt time.Time      `gorm:"index" json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.User        `gorm:"embedded"`
}

type TokenDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Token       `gorm:"embedded"`
}

type ActivityDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Activity    `gorm:"embedded"`
}

type HolidayDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Holiday     `gorm:"embedded"`
}

type WorkDayDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.WorkDay     `gorm:"embedded"`
}

type WebhookDB struct {
	ModelExtensions `gorm:"embedded"`
	pkg.Webhook     `gorm:"embedded"`
}
