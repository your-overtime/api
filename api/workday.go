package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

// GetWorkDays godoc
// @Tags workday
// @Summary Retrieves workdays
// @Produce json
// @Success 200 {object} []pkg.WorkDay
// @Param start query string true "Start date"
// @Param end query string true "End date"
// @Router /workday [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetWorkDays(c *gin.Context) {
	e, err := a.getUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
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
}

// CreateWorkDay godoc
// @Tags workday
// @Summary creates a workdays
// @Produce json
// @Consume json
// @Param workday body pkg.InputWorkDay true "Input workday"
// @Success 200 {object} pkg.WorkDay
// @Router /workday [post]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) CreateWorkDay(c *gin.Context) {
	e, err := a.getUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
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
		InputWorkDay: pkg.InputWorkDay{
			UserID:     e.ID,
			Day:        iw.Day,
			Overtime:   iw.Overtime,
			ActiveTime: iw.ActiveTime,
		},
	}
	h, err := a.os.AddWorkDay(wo, *e)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, h)
	}
}
