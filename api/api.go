package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

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

// GetOverview godoc
// @Summary Retrieves overview of your overtime
// @Produce json
// @Success 200 {object} pkg.Overview
// @Router /overview [get]
func (a *API) GetOverview(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	overview, err := a.os.CalcOverview(*e, time.Now())
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, overview)
	}
}

// StartActivity godoc
// @Summary Starts a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity/:desc [post]
func (a *API) StartActivity(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	fmt.Println(e)
	desc := c.Param("desc")
	ac, err := a.os.StartActivity(desc, *e)
	if err == pkg.ErrActivityIsRunning {
		c.JSON(http.StatusConflict, err.Error())
	} else if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, ac)
	}
}

// StopActivity godoc
// @Summary Stops a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity [delete]
func (a *API) StopActivity(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	ac, err := a.os.StopRunningActivity(*e)
	if err != nil && err == pkg.ErrNoActivityIsRunning {
		c.JSON(http.StatusOK, err)
	} else if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, ac)
	}
}

// CreateActivity godoc
// @Summary Creates a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity [post]
func (a *API) CreateActivity(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	var ia pkg.InputActivity
	err = c.Bind(&ia)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	act := pkg.Activity{
		UserID:      e.ID,
		Start:       ia.Start,
		End:         ia.End,
		Description: ia.Description,
	}
	ac, err := a.os.AddActivity(act, *e)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, ac)
	}
}

// UpdateActivity godoc
// @Summary Updates a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity/:id [put]
func (a *API) UpdateActivity(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	var ia pkg.InputActivity
	err = c.Bind(&ia)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	act := pkg.Activity{
		Model:       gorm.Model{ID: uint(id)},
		UserID:      e.ID,
		Start:       ia.Start,
		End:         ia.End,
		Description: ia.Description,
	}
	ac, err := a.os.AddActivity(act, *e)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, ac)
	}
}

// UpdateActivity godoc
// @Summary Get a activity by id
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity/:id [get]
func (a *API) GetActivity(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h, err := a.os.GetActivity(uint(id), *e)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, h)
	}
}

// UpdateActivity godoc
// @Summary Get a activities by start and end
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity [get]
func (a *API) GetActivities(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	start, err := time.Parse(time.RFC3339, c.Query("start"))
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	end, err := time.Parse(time.RFC3339, c.Query("end"))
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	h, err := a.os.GetActivities(start, end, *e)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, h)
	}
}

// UpdateActivity godoc
// @Summary Delete a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity/:id [get]
func (a *API) DeleteActivity(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err = a.os.DelActivity(uint(id), *e)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, "")
	}
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

	v1.POST("/holiday", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		var ih pkg.InputHoliday
		err = c.Bind(&ih)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		ho := pkg.Holiday{
			UserID:      e.ID,
			Start:       ih.Start,
			End:         ih.End,
			Type:        ih.Type,
			Description: ih.Description,
		}
		h, err := a.os.AddHoliday(ho, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.PUT("/holiday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		var ih pkg.InputHoliday
		err = c.Bind(&ih)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		ho := pkg.Holiday{
			Model:       gorm.Model{ID: uint(id)},
			UserID:      e.ID,
			Start:       ih.Start,
			End:         ih.End,
			Type:        ih.Type,
			Description: ih.Description,
		}
		h, err := a.os.AddHoliday(ho, *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.GET("/holiday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		h, err := a.os.GetHoliday(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.GET("/holiday", func(c *gin.Context) {
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
		h := []pkg.Holiday{}
		typeStr := c.Query("type")
		if len(typeStr) > 0 {
			hType, err := pkg.StrToHolidayType(typeStr)
			if err != nil {
				log.Debug(end, err)
				c.JSON(http.StatusBadRequest, err)
				return
			}
			h, err = a.os.GetHolidaysByType(start, end, hType, *e)
		} else {
			h, err = a.os.GetHolidays(start, end, *e)
		}
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.DELETE("/holiday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusBadRequest, err)
			return
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err = a.os.DelHoliday(uint(id), *e)
		if err != nil {
			log.Debug(err)
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	})
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
