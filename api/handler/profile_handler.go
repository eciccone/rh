package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/eciccone/rh/api/repo/profile"
	"github.com/eciccone/rh/api/rherr"
	"github.com/eciccone/rh/api/service"
	"github.com/gin-gonic/gin"
)

type profileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(s service.ProfileService) profileHandler {
	return profileHandler{s}
}

func (h *profileHandler) GetProfile(c *gin.Context) {
	profileId := c.GetString("sub")

	result, err := h.profileService.FetchProfile(profileId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "profile not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "profile found",
		"profile": result,
	})
}

func (h *profileHandler) PostProfile(c *gin.Context) {
	profileId := c.GetString("sub")

	var data profile.Profile
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid json data",
		})
		return
	}

	data.Id = profileId

	err := h.profileService.CreateProfile(data)
	if err != nil {
		if errors.Is(err, rherr.ErrProfileExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "profile already created",
			})
		} else if errors.Is(err, rherr.ErrUsernameTaken) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "username in use",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "profile created",
	})
}
