package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
	"gorm.io/gorm"
)

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
// @Param desc path string true "Activity description"
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
// @Param id path string true "Activity id"
// @Router /activity/{id} [put]
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

// GetActivity godoc
// @Summary Get a activity by id
// @Produce json
// @Success 200 {object} pkg.Activity
// @Param id path string true "Activity id"
// @Router /activity/{id} [get]
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

// GetActivities godoc
// @Summary Get a activities by start and end
// @Produce json
// @Param start query string true "Start date"
// @Param end query string true "Start date"
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

// DeleteActivity godoc
// @Summary Delete a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Param id path string true "Activity id"
// @Router /activity/{id} [get]
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
