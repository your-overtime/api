package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/your-overtime/api/pkg"
)

// GetTokens godoc
// @Tags token
// @Summary Retrieves tokens
// @Produce json
// @Success 200 {object} []pkg.Token
// @Router /token [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetTokens(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	ts, err := a.os.GetTokens(*e)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusCreated, ts)
	}
}

// CreateToken godoc
// @Tags token
// @Summary creates a token
// @Produce json
// @Consume json
// @Param bottles body pkg.InputToken true "input token"
// @Success 200 {object} pkg.Token
// @Router /token [post]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) CreateToken(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
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
}

// DeleteToken godoc
// @Tags token
// @Summary Delete a token
// @Produce json
// @Success 200 {object} pkg.Token
// @Param id path string true "Token id"
// @Router /token/{id} [delete]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) DeleteToken(c *gin.Context) {
	e, err := a.getEmployeeFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
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
}
