package data

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Db struct
type Db struct {
	Conn *gorm.DB
}

// Init function return Db
func Init(user string, pw string, host string, name string) (*Db, error) {
	db := Db{}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		pw,
		host,
		name,
	)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	db.Conn = conn

	// Migrate the schema
	conn.AutoMigrate(&pkg.Trace{})
	conn.AutoMigrate(&pkg.Hop{})
	conn.AutoMigrate(&pkg.Tracer{})

	return &db, err
}

func (d *Db) FetchTraces(start *time.Time, end *time.Time, query *string) []pkg.Trace {
	traces := []pkg.Trace{}
	if query != nil {
		d.Conn.Where("hostname like  ?", fmt.Sprintf("%%%s%%", *query))
	}
	if start != nil && end == nil {
		d.Conn.Where("execution_time > ?", *start)
	}

	if end != nil && start == nil {
		d.Conn.Where("execution_time > ?", *end)
	}

	if end != nil && start != nil {
		d.Conn.Where("execution_time BETWEEN ? AND", *start, *end)
	}

	d.Conn.Find(&traces)

	return traces
}

func (d *Db) FetchTracer(query *string, token *string) []pkg.Tracer {
	tracer := []pkg.Tracer{}
	if query != nil {
		d.Conn.Where("name like ?", fmt.Sprintf("%%%s%%", *query))
	}
	if token != nil {
		d.Conn.Where("access_token = ?", *token)
	}

	d.Conn.Find(&tracer)
	return tracer
}
