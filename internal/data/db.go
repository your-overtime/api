package data

import (
	"fmt"
	"net/url"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Db struct
type Db struct {
	Conn *gorm.DB
}

// Init function return Db
func Init(user string, pw string, host string, name string) (*Db, error) {
	tz := os.Getenv("TZ")
	if len(tz) == 0 {
		tz = time.Now().Location().String()
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		url.QueryEscape(user),
		url.QueryEscape(pw),
		host,
		url.QueryEscape(name),
		url.QueryEscape(tz),
	)
	return InitWithDialector(mysql.Open(dsn))
}

func InitWithDialector(dialector gorm.Dialector) (*Db, error) {
	db := Db{}
	conn, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	db.Conn = conn

	// Migrate the schema
	conn.AutoMigrate(&ActivityDB{})
	conn.AutoMigrate(&UserDB{})
	conn.AutoMigrate(&TokenDB{})
	conn.AutoMigrate(&HolidayDB{})
	conn.AutoMigrate(&WorkDayDB{})
	conn.AutoMigrate(&WebhookDB{})

	(&db).MigrateActivityDuration()

	return &db, err
}
