package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"git.goasum.de/jasper/overtime/pkg"
	"github.com/gin-gonic/gin"
)

// API struct
type API struct {
	os     pkg.OvertimeService
	es     pkg.EmployeeService
	router *gin.Engine
	host   string
}

func (a *API) getEmployeeFromRequest(c *gin.Context) (*pkg.Employee, error) {
	token := c.Request.FormValue("token")
	if len(token) > 0 {
		return a.es.FromToken(token)
	}
	authHeaderSlice := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authHeaderSlice) == 2 {
		switch strings.ToLower(authHeaderSlice[1]) {
		case "basic":
			payload := []byte{}
			_, err := base64.StdEncoding.Decode(payload, []byte(authHeaderSlice[2]))
			if err != nil {
				return nil, pkg.ErrUserNotFound
			}
			basicAuth := strings.Split(string(payload), ":")
			if len(basicAuth) != 2 {
				return nil, pkg.ErrUserNotFound
			}
			return a.es.Login(basicAuth[0], basicAuth[1])
		default:
			return a.es.FromToken(authHeaderSlice[1])
		}

	}

	return nil, pkg.ErrUserNotFound
}

func (a *API) createEndPoints() {
	api := a.router.Group("/api")

	v1 := api.Group("/v1")
	v1.GET("overview/current", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		overview, err := a.os.CalcCurrentOverview(*e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, overview)
		}
	})
	v1.GET("overview", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		overview, err := a.os.CalcOverviewForThisYear(*e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, overview)
		}
	})
	v1.POST("/activity/:desc", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		fmt.Println(e)
		desc := c.Param("desc")
		ac, err := a.os.StartActivity(desc, *e)
		if err == pkg.ErrActivityIsRunning {
			c.JSON(http.StatusConflict, err.Error())
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.DELETE("/activity", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		ac, err := a.os.StopRunningActivity(*e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.POST("/activity", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		var ia pkg.InputActivity
		err = c.Bind(&a)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		act := pkg.Activity{
			UserID:      e.ID,
			Start:       &ia.Start,
			End:         &ia.End,
			Description: ia.Description,
		}
		ac, err := a.os.AddActivity(act, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, ac)
		}
	})
	v1.GET("/activity/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		h, err := a.os.GetActivity(uint(id), *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.DELETE("/activity/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err = a.os.DelActivity(uint(id), *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	})
	v1.POST("/hollyday", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		var ih pkg.InputHollyday
		err = c.Bind(a)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		ho := pkg.Hollyday{
			UserID:      e.ID,
			Start:       ih.Start,
			End:         ih.End,
			Description: ih.Description,
		}
		h, err := a.os.AddHollyday(ho, *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.GET("/hollyday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		h, err := a.os.GetHollyday(uint(id), *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, h)
		}
	})
	v1.DELETE("/hollyday/:id", func(c *gin.Context) {
		e, err := a.getEmployeeFromRequest(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err)
		}
		id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
		err = a.os.DelHollyday(uint(id), *e)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
		} else {
			c.JSON(http.StatusOK, "")
		}
	})
}

// Init API server
func Init(os pkg.OvertimeService, es pkg.EmployeeService) *API {
	return &API{
		router: gin.Default(),
		os:     os,
		es:     es,
	}
}

// Start API server
func (a API) Start(host string) {
	a.createEndPoints()
	panic(a.router.Run(host))
}
