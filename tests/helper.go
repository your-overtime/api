package tests

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/your-overtime/api/internal/data"
	"github.com/your-overtime/api/pkg"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupDb(t *testing.T) data.Db {
	conn, err := gorm.Open(sqlite.Open("file::memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db := data.Db{Conn: conn}

	if err := conn.AutoMigrate(&pkg.Holiday{}); err != nil {
		t.Fatal(err)
	}
	if err := conn.AutoMigrate(&pkg.Activity{}); err != nil {
		t.Fatal(err)
	}
	if err := conn.AutoMigrate(&pkg.WorkDay{}); err != nil {
		t.Fatal(err)
	}
	return db
}

func Equals(t *testing.T, actual interface{}, expected interface{}) bool {
	aBytes, aErr := json.Marshal(actual)
	if aErr != nil {
		t.Fatal(aErr)
	}
	bBytes, bErr := json.Marshal(expected)
	if bErr != nil {
		t.Fatal(bErr)
	}

	return bytes.Equal(aBytes, bBytes)
}

func ParseDay(val string) time.Time {
	date, err := time.Parse("2006-01-02", val)

	if err != nil {
		panic(err)
	}
	return date
}

func ParseDayTime(val string) time.Time {
	daytime, err := time.Parse("2006-01-02 15:04", val)
	if err != nil {
		panic(err)
	}
	return daytime
}
