package main

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/your-overtime/api/v2/api"
	docs "github.com/your-overtime/api/v2/docs"
	"github.com/your-overtime/api/v2/internal/data"
	"github.com/your-overtime/api/v2/internal/service"

	log "github.com/sirupsen/logrus"
)

var version = "2.0.0"

// @title Your Overtime API
// @version 1.0
// @BasePath /api/v1
// @securityDefinitions.basic BasicAuth
//
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey AdminAuth
// @in query
// @name adminToken
func main() {
	docs.SwaggerInfo.Schemes = []string{"https", "http"}
	docs.SwaggerInfo.Version = version

	godotenv.Load()
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
	if strings.ToLower(os.Getenv("DEBUG")) == "true" {
		log.SetLevel(log.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	if len(os.Getenv("TZ")) > 0 {
		loc, err := time.LoadLocation(os.Getenv("TZ"))
		if err != nil {
			panic(err)
		}
		time.Local = loc
	}

	db, err := data.Init(
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
	if err != nil {
		panic(err)
	}

	ovs := service.Init(db)
	api := api.Init(ovs, os.Getenv("ADMIN_TOKEN"))
	api.Start(os.Getenv("HOST"))
}
