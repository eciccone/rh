package handler

import (
	"errors"
	"net/http"

	"github.com/eciccone/rh/api/repo/profile"
	"github.com/eciccone/rh/api/service"
	"github.com/gin-gonic/gin"
)

type profileHandler struct {
	profileService service.ProfileService
}

func NewProfileHandler(s service.ProfileService) profileHandler {
	return profileHandler{s}
}

func (h *profileHandler) GetProfile(c *gin.Context) error {
	profileId := c.GetString("sub")
	if profileId == "" {
		return errors.New("GetProfile failed to get subject, should have been set in middleware")
	}

	result, err := h.profileService.FetchProfile(profileId)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "profile found",
		"profile": result,
	})

	return nil
}

func (h *profileHandler) PostProfile(c *gin.Context) error {
	profileId := c.GetString("sub")
	if profileId == "" {
		return errors.New("PostProfile failed to get subject, should have been set in middleware")
	}

	var data profile.Profile
	if err := c.ShouldBindJSON(&data); err != nil {
		return ErrInvalidJSON
	}

	data.Id = profileId

	err := h.profileService.CreateProfile(data)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "profile created",
	})

	return nil
}
