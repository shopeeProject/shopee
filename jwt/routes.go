package jwthandler

import (
	"github.com/gin-gonic/gin"
	"github.com/shopeeProject/shopee/util"
)

const (
	TokenRefresh = "/token-refresh"
)

type RefreshTokenInput struct {
	RefreshToken string `json:"refreshToken"`
}

func TokenRefreshHandler(r *util.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		refreshToken := RefreshTokenInput{}
		c.Bind(&refreshToken)
		if refreshToken.RefreshToken == "" {
			c.JSON(409, gin.H{"message": "Empty Refresh Token recieved"})
			return
		}
		RefreshResponse := Refresh(refreshToken.RefreshToken, r)
		if !RefreshResponse.Success {
			c.JSON(409, gin.H{"message": RefreshResponse.Message})
			return
		}
		c.JSON(200, gin.H{
			"message": RefreshResponse.Message,
			"data":    RefreshResponse.Data,
		})
	}
}

func RegisterRoutes(router *gin.Engine, r *util.Repository) *gin.RouterGroup {
	// todo validate if request is coming from valid seller/user
	v1 := router.Group("/")
	{
		v1.POST(TokenRefresh, TokenRefreshHandler(r))
	}
	return v1
}
