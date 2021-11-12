package api

import (
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
	"github.com/your-overtime/api/internal/service"
	"github.com/your-overtime/api/pkg"

	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	docs "github.com/your-overtime/api/docs"
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
			c.AbortWithError(http.StatusUnauthorized, errors.New("Invalid token"))
		}
	}
}

func (a *API) getEmployeeFromRequest(c *gin.Context) (*pkg.Employee, error) {
	token := c.Request.FormValue("token")
	if len(token) > 0 {
		return a.os.FromToken(token)
	}
	authHeaderSlice := strings.Split(c.Request.Header.Get("Authorization"), " ")
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

	a.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.GET("overview", a.GetOverview)
	v1.POST("/activity/:desc", a.StartActivity)
	v1.DELETE("/activity", a.StopActivity)
	v1.POST("/activity", a.CreateActivity)
	v1.PUT("/activity/:id", a.UpdateActivity)
	v1.GET("/activity/:id", a.GetActivity)
	v1.GET("/activity", a.GetActivities)
	v1.DELETE("/activity/:id", a.DeleteActivity)

	v1.POST("/holiday", a.CreateHoliday)
	v1.PUT("/holiday/:id", a.UpdateHoliday)
	v1.GET("/holiday/:id", a.GetHoliday)
	v1.GET("/holiday", a.GetHolidays)
	v1.DELETE("/holiday/:id", a.DeleteHoliday)

	v1.GET("/workday", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		start, err := time.Parse(time.RFC3339Nano, c.Query("start"))
		if err != nil {
			log.Debug(start, err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		end, err := time.Parse(time.RFC3339Nano, c.Query("end"))
		if err != nil {
			log.Debug(end, err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		wds, err := a.os.GetWorkDays(start, end, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, wds)
		}
	})
	v1.POST("/workday", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		var iw pkg.InputWorkDay
		err = c.Bind(&iw)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		wo := pkg.WorkDay{
			UserID:     e.ID,
			Day:        iw.Day,
			Overtime:   iw.Overtime,
			ActiveTime: iw.ActiveTime,
		}
		h, err := a.os.AddWorkDay(wo, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.POST("/token", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		var it pkg.InputToken
		err = c.Bind(&it)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		t, err := a.os.CreateToken(it, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusCreated, t)
		}
	})
	v1.GET("/token", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}

		ts, err := a.os.GetTokens(*e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusCreated, ts)
		}
	})
	v1.DELETE("/token/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err = a.os.DeleteToken(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "token deleted")
		}
	})
	v1.GET("account", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		} else {
			c.JSON(http.StatusOK, e)
		}
	})
	v1.PATCH("account", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		var payload map[string]interface{}
		err = c.Bind(&payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
			return
		}
		e, err = a.os.UpdateAccount(payload, *e)
		if err != nil {
			log.Debug(err)
			if errors.Is(err, pkg.ErrDuplicateValue) {
				c.JSON(http.StatusBadRequest, err)
			} else {
				c.JSON(http.StatusInternalServerError, err)
			}
		} else {
			c.JSON(http.StatusOK, e)
		}
	})
	authorizedV1 := v1.Group("/", a.adminAuth())
	{
		authorizedV1.POST("/employee", func(c *gin.Context) {
			var ie pkg.InputEmployee
			err := c.Bind(&ie)
			if err != nil {
				log.Debug(err)
				c.JSON(http.StatusBadRequest, err)
				return
			}
			e, err := a.os.SaveEmployee(ie.ToEmployee(), "")
			if err != nil {
				log.Debug(err)
				c.JSON(http.StatusInternalServerError, err)
			} else {
				c.JSON(http.StatusCreated, e)
			}
		})
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
	docs.SwaggerInfo.Title = "Your Overtime Swagger API"
	docs.SwaggerInfo.BasePath = "/api/v1/"

	a.createEndPoints()
	panic(a.router.Run(host))
}
