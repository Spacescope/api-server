package core

import (
	"api-server/pkg/utils"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const USERID = "user_id"
const AUTHORIZATION = "authorization"

func AuthFromGateway() gin.HandlerFunc {
	return func(c *gin.Context) {
		app := utils.Gin{C: c}
		if len(c.Request.Header.Get(USERID)) == 0 {
			app.HTTPResponse(http.StatusForbidden, utils.ErrDataInfraAPIForbidden)
			c.Abort()
			return
		}

		if len(c.Request.Header.Get(AUTHORIZATION)) == 0 {
			app.HTTPResponse(http.StatusForbidden, utils.ErrDataInfraAPIForbidden)
			c.Abort()
			return
		}

		c.Next()
	}
}

func GinUidLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		fmt.Print(c.Request.Header.Get(USERID) + " ")
	}
}
