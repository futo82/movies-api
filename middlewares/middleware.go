package middlewares

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	verifier "github.com/okta/okta-jwt-verifier-golang"
)

// AuthMiddleware authenicates the request JWT token
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized."})
			return
		}
		tokenParts := strings.Split(authHeader, "Bearer ")
		if len(tokenParts) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized."})
			return
		}
		bearerToken := tokenParts[1]

		tv := map[string]string{}
		tv["aud"] = "api://default"
		tv["cid"] = os.Getenv("CLIENT_ID")
		jv := verifier.JwtVerifier{
			Issuer:           os.Getenv("ISSUER"),
			ClaimsToValidate: tv,
		}

		_, err := jv.New().VerifyAccessToken(bearerToken)

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Not Authorized."})
			return
		}

		c.Next()
	}
}
