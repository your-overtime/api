package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/your-overtime/api/internal/service"
	"github.com/your-overtime/api/pkg"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// API struct
type API struct {
	os         *service.Service
	router     *gin.Engine
	host       string
	adminToken string
}

func (a *API) adminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.FormValue("adminToken")
		log.Debug(token, " ", a.adminToken, " ", token == a.adminToken)
		if token != a.adminToken {
			c.AbortWithError(http.StatusUnauthorized, errors.New("invalid token"))
		}
	}
}

func (a *API) getEmployeeFromRequest(c *gin.Context) (*pkg.Employee, error) {
	token := c.Request.FormValue("token")
	if len(token) > 0 {
		return a.os.FromToken(token)
	}
	authHeaderSlice := strings.Split(c.Request.Header.Get("Authorization"), " ")
	fmt.Println(authHeaderSlice)
	if len(authHeaderSlice) == 2 {
		switch strings.ToLower(authHeaderSlice[0]) {
		case "basic":
			payload, err := base64.StdEncoding.DecodeString(authHeaderSlice[1])
			if err != nil {
				log.Debug(err)
				return nil, pkg.ErrUserNotFound
			}
			basicAuth := strings.Split(string(payload), ":")
			if len(basicAuth) != 2 {
				log.Debug(pkg.ErrUserNotFound, " ", basicAuth, " ", authHeaderSlice)
				return nil, pkg.ErrUserNotFound
			}
			return a.os.Login(basicAuth[0], basicAuth[1])
		default:
			return a.os.FromToken(authHeaderSlice[1])
		}

	}

	return nil, pkg.ErrUserNotFound
}

func (a *API) createEndPoints() {
	api := a.router.Group("/api")

	v1 := api.Group("/v1")

	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Overview
	v1.GET("/overview", a.GetOverview)

	// activity
	v1.POST("/activity/:desc", a.StartActivity)
	v1.DELETE("/activity", a.StopActivity)
	v1.POST("/activity", a.CreateActivity)
	v1.PUT("/activity/:id", a.UpdateActivity)
	v1.GET("/activity/:id", a.GetActivity)
	v1.GET("/activity", a.GetActivities)
	v1.DELETE("/activity/:id", a.DeleteActivity)

	// holiday
	v1.POST("/holiday", a.CreateHoliday)
	v1.PUT("/holiday/:id", a.UpdateHoliday)
	v1.GET("/holiday/:id", a.GetHoliday)
	v1.GET("/holiday", a.GetHolidays)
	v1.DELETE("/holiday/:id", a.DeleteHoliday)

	// workday
	v1.GET("/workday", a.GetWorkDays)
	v1.POST("/workday", a.CreateWorkDay)

	// token
	v1.GET("/token", a.GetTokens)
	v1.POST("/token", a.CreateToken)
	v1.DELETE("/token/:id", a.DeleteToken)

	// account
	v1.GET("account", a.GetAccount)
	v1.PATCH("account", a.UpdateAccount)

	// employee
	authorizedV1 := v1.Group("/", a.adminAuth())
	{
		authorizedV1.POST("/employee", a.CreateEmployee)
	}
}

// Init API server
func Init(os *service.Service, adminToken string) *API {
	return &API{
		router:     gin.Default(),
		os:         os,
		adminToken: adminToken,
	}
}

func (a API) Start(host string) {
	a.createEndPoints()
	panic(a.router.Run(host))
}
