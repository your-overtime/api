package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
	"gorm.io/gorm"
)

// CreateHoliday godoc
// @Summary Creates a holiday
// @Produce json
// @Success 200 {object} pkg.Holiday
// @Router /holiday [post]
func (a *API) CreateHoliday(c *gin.Context) {
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
}

// UpdateHoliday godoc
// @Summary Updates a holiday
// @Produce json
// @Success 200 {object} pkg.Holiday
// @Param id path string true "Holiday id"
// @Router /holiday/{id} [put]
func (a *API) UpdateHoliday(c *gin.Context) {
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
}

// GetHoliday godoc
// @Summary Get a holiday by id
// @Produce json
// @Success 200 {object} pkg.Holiday
// @Param id path string true "Holiday id"
// @Router /holiday/{id} [get]
func (a *API) GetHoliday(c *gin.Context) {
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
}

// GetHolidays godoc
// @Summary Get a activities by start and end
// @Produce json
// @Param start query string true "Start date"
// @Param end query string true "Start date"
// @Success 200 {object} pkg.Holiday
// @Router /holiday [get]
func (a *API) GetHolidays(c *gin.Context) {
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
}

// UpdateHoliday godoc
// @Summary Delete a holiday
// @Produce json
// @Success 200 {object} pkg.Holiday
// @Param id path string true "Holiday id"
// @Router /holiday/{id} [get]
func (a *API) DeleteHoliday(c *gin.Context) {
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
}
