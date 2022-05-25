package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type authorizationHeader struct {
	IDToken string `header:"Authorization"`
}

func Validate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := &authorizationHeader{}

		if err := c.ShouldBindHeader(&header); err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": "internal server error",
			})
			return
		}

		bearerAndToken := strings.Split(header.IDToken, "Bearer ")
		if len(bearerAndToken) < 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "invalid access token",
			})
			return
		}

		token, err := jwt.Parse(
			[]byte(bearerAndToken[1]),
			jwt.WithKeySet(fetchTenantKeys()),
			jwt.WithValidate(true),
			jwt.WithAudience(os.Getenv("AUTH0_AUDIENCE")),
			jwt.WithAcceptableSkew(time.Minute))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "invalid access token",
			})
			return
		}

		c.Set("sub", token.Subject())
		c.Next()
	}

}

func fetchTenantKeys() jwk.Set {
	urlString := fmt.Sprintf("https://%s/.well-known/jwks.json", os.Getenv("AUTH0_DOMAIN"))
	set, err := jwk.Fetch(context.Background(), urlString)
	if err != nil {
		log.Fatalf("failed to parse tenant json web keys: %s\n", err)
	}
	return set
}
