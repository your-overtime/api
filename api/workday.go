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
// @Param start query string true "Start date" format date-time
// @Param end query string true "End date" format date-time
// @Router /workday [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetWorkDays(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
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
	wds, err := os.GetWorkDays(start, end)
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
	os, err := a.getOvertimeServiceForUserFromRequest(c)
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
	user, _ := os.GetAccount()
	wo := pkg.WorkDay{
		InputWorkDay: pkg.InputWorkDay{
			UserID:     user.ID,
			Day:        iw.Day,
			Overtime:   iw.Overtime,
			ActiveTime: iw.ActiveTime,
		},
	}
	h, err := os.AddWorkDay(wo)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, h)
	}
}
