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
// @Success 201 {object} []pkg.Token
// @Router /token [get]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) GetTokens(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}

	ts, err := os.GetTokens()
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
// @Param token body pkg.InputToken true "Input token"
// @Success 201 {object} pkg.Token
// @Router /token [post]
// @Security BasicAuth
// @Security ApiKeyAuth
func (a *API) CreateToken(c *gin.Context) {
	os, err := a.getOvertimeServiceForUserFromRequest(c)
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
	t, err := os.CreateToken(it)
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
	os, err := a.getOvertimeServiceForUserFromRequest(c)
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	err = os.DeleteToken(uint(id))
	if err != nil {
		log.Debug(err)
		c.JSON(http.StatusInternalServerError, err)
	} else {
		c.JSON(http.StatusOK, "token deleted")
	}
}
