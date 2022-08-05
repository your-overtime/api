package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/v2/pkg"
)

// StopActivity godoc
// @Tags activity
// @Summary Stops a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Router /activity/stop [delete]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) StopActivity(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	ac, err := os.StopRunningActivity()
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
// @Tags activity
// @Summary Creates a activity
// @Produce json
// @Consume json
// @Param activity body pkg.InputActivity true "input activity"
// @Success 201 {object} pkg.Activity
// @Router /activity [post]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) CreateActivity(c *gin.Context) {
	var (
		ac  *pkg.Activity
		err error
	)

	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	var ia pkg.InputActivity
	err = c.BindJSON(&ia)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if ia.Start == nil {
		ac, err = os.StartActivity(ia.Description)
	} else {
		user, _ := os.GetAccount()
		act := pkg.Activity{
			UserID: user.ID,
			InputActivity: pkg.InputActivity{
				Start:       ia.Start,
				End:         ia.End,
				Description: ia.Description,
			},
		}
		ac, err = os.AddActivity(act)
	}
	if err == pkg.ErrEmptyDescriptionNotAllowed {
		c.JSON(http.StatusBadRequest, err.Error())
	} else if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusCreated, ac)
	}
}

// UpdateActivity godoc
// @Tags activity
// @Summary Updates a activity
// @Produce json
// @Consume json
// @Param activity body pkg.InputActivity true "input activity"
// @Success 200 {object} pkg.Activity
// @Param id path string true "Activity id"
// @Router /activity/{id} [put]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) UpdateActivity(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
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
		ID: uint(id),
		InputActivity: pkg.InputActivity{
			Start:       ia.Start,
			End:         ia.End,
			Description: ia.Description,
		},
	}
	ac, err := os.UpdateActivity(act)
	if err == pkg.ErrEmptyDescriptionNotAllowed {
		c.JSON(http.StatusBadRequest, err.Error())
	} else if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, ac)
	}
}

// GetActivity godoc
// @Tags activity
// @Summary Get a activity by id
// @Produce json
// @Success 200 {object} pkg.Activity
// @Param id path string true "Activity id"
// @Router /activity/{id} [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetActivity(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	h, err := os.GetActivity(uint(id))
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, h)
	}
}

// GetActivities godoc
// @Tags activity
// @Summary Get a activities by start and end
// @Produce json
// @Param start query string true "Start date" format date-time
// @Param end query string true "End date" format date-time
// @Success 200 {object} []pkg.Activity
// @Router /activity [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetActivities(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
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
	h, err := os.GetActivities(start, end)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, h)
	}
}

// DeleteActivity godoc
// @Tags activity
// @Summary Delete a activity
// @Produce json
// @Success 200 {object} pkg.Activity
// @Param id path string true "Activity id"
// @Router /activity/{id} [delete]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) DeleteActivity(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err = os.DelActivity(uint(id))
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, "activity deletet")
	}
}
