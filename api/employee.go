package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

// CreateEmployee godoc
// @Tags employee
// @Summary creates a employee
// @Produce json
// @Consume json
// @Param bottles body pkg.InputEmployee true "Input employee"
// @Success 200 {object} pkg.Employee
// @Router /employee [post]
// @Security AdminAuth
func (a *API) CreateEmployee(c *gin.Context) {
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
}
