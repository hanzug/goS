package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/hanzug/goS/consts"
	e2 "github.com/hanzug/goS/consts/e"
	"github.com/hanzug/goS/pkg/ctl"
	"github.com/hanzug/goS/pkg/jwt"
)

// AuthMiddleware jwt鉴权中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e2.SUCCESS
		accessToken := c.GetHeader("access_token")
		refreshToken := c.GetHeader("refresh_token")

		if accessToken == "" {
			code = e2.InvalidParams
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e2.GetMsg(code),
				"data":   "accessToken 为空",
			})
			c.Abort()
			return
		}

		// todo：双token机制待完善，目前直接读取了两个token，可以先读取accessToken，如果过期则请求refreshToken
		newAccessToken, newRefreshToken, err := jwt.ParseRefreshToken(accessToken, refreshToken)
		if err != nil {
			code = e2.ErrorAuthCheckTokenFail
		}
		if code != e2.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e2.GetMsg(code),
				"data":   "jwt鉴权失败",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(newAccessToken)
		if err != nil {
			code = e2.ErrorAuthCheckTokenFail
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e2.GetMsg(code),
				"data":   err.Error(),
			})
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.ID, UserName: claims.Username}))
		SetToken(c, newAccessToken, newRefreshToken)
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}

func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(consts.AccessTokenHeader, accessToken)
	c.Header(consts.RefreshTokenHeader, refreshToken)
	//储存到浏览器的cookie中，每次自动把token加到头部
	c.SetCookie(consts.AccessTokenHeader, accessToken, consts.MaxAge, "/", "", secure, true)
	c.SetCookie(consts.RefreshTokenHeader, refreshToken, consts.MaxAge, "/", "", secure, true)
}

func IsHttps(c *gin.Context) bool {
	if c.GetHeader(consts.HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	} else {
		return false
	}
}
