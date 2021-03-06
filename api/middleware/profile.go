package middleware

import (
	"errors"
	"net/http"

	"github.com/eciccone/rh/api/service"
	"github.com/gin-gonic/gin"
)

func Profile(ps service.ProfileService) gin.HandlerFunc {
	return func(c *gin.Context) {
		profileID := c.GetString("sub")

		profile, err := ps.FetchProfile(profileID)
		if err != nil {
			if errors.Is(err, service.ErrNoProfile) {
				c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"msg": "profile not found",
				})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"msg": "internal server error",
				})
			}
			return
		}

		c.Set("username", profile.Username)
		c.Next()
	}
}
