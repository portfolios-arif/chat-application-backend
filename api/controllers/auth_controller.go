package controllers

import "github.com/gin-gonic/gin"

type AuthControllerInterface interface {
	CheckEmail(ctx *gin.Context)
	SendOtp(ctx *gin.Context)
	ValidateOtp(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	RefreshToken(ctx *gin.Context)
	Logout(ctx *gin.Context)
}
